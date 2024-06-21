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

func mapJSONToString(m map[string]string) string {
	var str string
	for _, value := range m {
		str += value + ","
	}
	// Remove the last comma
	if len(str) > 0 {
		str = str[:len(str)-1]
	}
	return str
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
		fmt.Println(len(msgChannel))
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(msg.DataJSON), &data); err != nil {
			log.Error().Msg(fmt.Sprintf("json unmarshal parsed payload: %v\n", err))
		}
		formmatedData := make(map[string]string)
		for i, body := range data {
			jsonData, err := json.Marshal(map[string]interface{}{
				"name":  i,
				"value": body,
			})
			if err != nil {
				log.Error().Msg(fmt.Sprintf("json marshal parsed payload: %v\n", err))
				continue
			}
			formmatedData[i] = string(jsonData)
		}
		messageKafka := KafkaMessageSchema{
			OpName: msg.OperationName,
			OpCode: msg.OperationID,
			Body:   json.RawMessage("[" + mapJSONToString(formmatedData) + "]"),
		}
		messageKafkaJSON, err := json.Marshal(messageKafka)
	
		topic := msg.OperationName
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          messageKafkaJSON,
		}, nil)
		
		if err != nil {
			log.Error().Msg(fmt.Sprintf("json marshal kafka message: %v\n", err))
		}
	}
}