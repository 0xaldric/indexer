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
		formmatedData := make(map[string]string)
		for i, body := range data {
			switch body := body.(type) {
			case string:
				formmatedData[i] = body
			}
		}
		jsonData, err := json.Marshal(formmatedData)
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
			Value:          messageKafkaJSON,
		}, nil)
		
		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal kafka message: %v\n", err))
		}
	}
}