package indexer

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/allisson/go-env"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"

	"github.com/tonindexer/anton/internal/core"
)

type KafkaMessageSchema struct {
	OpName        string          `json:"op_name"`
	OpCode        uint32          `json:"op_code"`
	Hash          string          `json:"hash"`
	Value         uint64          `json:"value"`
	CreatedAt     int64           `json:"created_at"`
	CreatedLT     uint64          `json:"created_lt"`
	SrcAddress    string          `json:"src_address"`
	DstAddress    string          `json:"dst_address"`
	SrcBlockSeqNo uint32          `json:"src_block_seq_no"`
	DstBlockSeqNo uint32          `json:"dst_block_seq_no"`
	Body          json.RawMessage `json:"body"`
	BodyHex       string          `json:"body_hex"`
	SrcTxHash     string          `json:"src_tx_hash"`
	DstTxHash     string          `json:"dst_tx_hash"`
}

func (s *Service) createTopics() {
	NumPartitions := env.GetInt("NUM_PARTITIONS", 2)
	operations, err := s.contractRepo.GetOperations(context.Background())
	if err != nil {
		log.Error().Msg(fmt.Sprintf("get operations: %v\n", err))
	}
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": s.KafkaURI})
	if err != nil {
		log.Error().Msg(fmt.Sprintf("create admin client: %v\n", err))
	}
	defer adminClient.Close()
	for _, i := range operations {
		topicName := string(i.OperationName)
		adminClient.CreateTopics(context.Background(), []kafka.TopicSpecification{{
			Topic:         topicName,
			NumPartitions: NumPartitions,
			Config: map[string]string{
				"retention.ms":   "600000",
				"cleanup.policy": "delete",
				"segment.ms":     "600000",
				"segment.bytes":  "104857600",
			},
		}})
	}
}

func (s *Service) produceMessageLoop(msgChannel <-chan *core.Message) {
	s.createTopics()
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error().Msg(fmt.Sprintf("Delivery failed: %v\n", ev.TopicPartition))
				}
			}
		}
	}()
	for s.running() {
		msg := <-msgChannel
		if msg.OperationName == "" {
			log.Error().Msg("operation name is empty")
			continue
		}
		var data map[string]json.RawMessage
		if err := json.Unmarshal([]byte(msg.DataJSON), &data); err != nil {
			log.Error().Msg(fmt.Sprintf("json unmarshal parsed payload: %v\n", err))
		}
		formatData := make(map[string]json.RawMessage)
		for i, body := range data {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("json marshal body: %v\n", err))
			}
			formatData[i] = json.RawMessage(jsonBody)
		}
		jsonData, err := json.Marshal(formatData)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal format data: %v\n", err))
		}
		// get message data
		messageKafka := KafkaMessageSchema{
			OpName:        msg.OperationName,
			OpCode:        msg.OperationID,
			Hash:          base64.StdEncoding.EncodeToString(msg.Hash),
			Value:         uint64(0),
			CreatedAt:     msg.CreatedAt.Unix(),
			CreatedLT:     msg.CreatedLT,
			SrcBlockSeqNo: msg.SrcBlockSeqNo,
			DstBlockSeqNo: msg.DstBlockSeqNo,
			Body:          jsonData,
			BodyHex:       hex.EncodeToString(msg.Body),
			SrcTxHash:     hex.EncodeToString(msg.SrcTxHash),
			DstTxHash:     hex.EncodeToString(msg.DstTxHash),
		}
		if msg.SrcAddress.MustToTonutils() != nil {
			messageKafka.SrcAddress = msg.SrcAddress.MustToTonutils().String()
		}
		if msg.DstAddress.MustToTonutils() != nil {
			messageKafka.DstAddress = msg.DstAddress.MustToTonutils().String()
		}
		messageKafkaJSON, err := json.Marshal(messageKafka)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal kafka message: %v\n", err))
		}
		topicName := msg.OperationName
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
			Value:          json.RawMessage(messageKafkaJSON),
		}, nil)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal kafka message: %v\n", err))
		}
	}
}
