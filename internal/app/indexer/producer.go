package indexer

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"

	"github.com/tonindexer/anton/internal/core"
)

type KafkaMessageSchema struct {
	OpName string `json:"op_name"`
	OpCode uint32 `json:"op_code"`
	Body  json.RawMessage `json:"body"`
}

func (s *Service) produceMessageLoop(msgChannel <-chan *core.Message) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"});
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
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg.DataJSON), &data); err != nil {
			log.Error().Msg(fmt.Sprintf("json unmarshal parsed payload: %v\n", err))
		}
		formmatedData := make(map[string]json.RawMessage)
		data["srcAddress"] = msg.SrcAddress.MustToTonutils().String()
		data["dstAddress"] = msg.DstAddress.MustToTonutils().String()
		data["timestamp"] = msg.CreatedAt.Unix()
		for i, body := range data {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				log.Error().Msg(fmt.Sprintf("json marshal body: %v\n", err))
			}
			formmatedData[i] = json.RawMessage(jsonBody)
		}
		jsonData, err := json.Marshal(formmatedData)
		fmt.Println(formmatedData)
		fmt.Println(jsonData)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal formmated data: %v\n", err))
		}
		messageKafka := KafkaMessageSchema{
			OpName: msg.OperationName,
			OpCode: msg.OperationID,
			Body:   json.RawMessage(jsonData),
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