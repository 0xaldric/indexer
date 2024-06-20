package parser

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/tonindexer/anton/abi/known"
	"github.com/tonindexer/anton/internal/app"
	"github.com/tonindexer/anton/internal/core"
)

type KafkaMessageSchema struct {
	OpName string `json:"op_name"`
	OpCode uint32 `json:"op_code"`
	Body  string `json:"body"`
}

func mapToString(m map[string]string) string {
	var str string
	for key, value := range m {
			str += fmt.Sprintf("%s: %s, ", key, value)
	}
	return str
}

func parseOperationAttempt(msg *core.Message, op *core.ContractOperation) error {
	msg.OperationName = op.OperationName
	if op.Outgoing {
		msg.SrcContract = op.ContractName
	} else {
		msg.DstContract = op.ContractName
	}

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
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()	
	topic := msg.OperationName

	payloadCell, err := cell.FromBOC(msg.Body)
	if err != nil {
		return errors.Wrap(err, "msg body from boc")
	}

	msgParsed, err := op.Schema.FromCell(payloadCell)
	if err != nil {
		return errors.Wrap(err, "msg body from boc")
	}

	msg.DataJSON, err = json.Marshal(msgParsed)
	if err != nil {
		return errors.Wrap(err, "json marshal parsed payload")
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(msg.DataJSON), &data); err != nil {
		return errors.Wrap(err, "json unmarshal parsed payload")
	}

	formattedData := make(map[string]string)
	for _, body := range op.Schema.Body {
		if value, ok := data[body.Name]; ok {
			// Marshal the struct directly
			jsonData, err := json.Marshal(map[string]interface{}{
				"name":  body.Name,
				"type":  body.Type,
				"value": value,
			})
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				continue
			}
			// Store the JSON string in formattedData map
			formattedData[body.Name] = string(jsonData)
		}
	}
	formmatedDataJSON, err := json.Marshal(formattedData)
	if err != nil {
		return errors.Wrap(err, "json marshal formatted data")
	}
	messageKafka := KafkaMessageSchema{
		OpName: op.OperationName,
		OpCode: op.OperationID,
		Body:   string(formmatedDataJSON),
	}
	messageKafkaJSON, err := json.Marshal(messageKafka)
	if err != nil {
		return errors.Wrap(err, "json marshal kafka message")
	}
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          messageKafkaJSON,
	}, nil)

	p.Flush(1000 * 1000)

	return nil
}

func (s *Service) parseDirectedMessage(ctx context.Context, acc *core.AccountState, msg *core.Message) error {
	if acc == nil {
		return errors.Wrap(app.ErrImpossibleParsing, "no account data")
	}
	if len(acc.Types) == 0 {
		return errors.Wrap(app.ErrImpossibleParsing, "no interfaces")
	}

	outgoing := acc.Address == msg.SrcAddress
	if outgoing && len(acc.Types) == 1 {
		msg.SrcContract = acc.Types[0]
	}
	if !outgoing && len(acc.Types) == 1 {
		msg.DstContract = acc.Types[0]
	}

	operations, err := s.ContractRepo.GetOperationsByID(ctx, msg.Type, acc.Types, outgoing, msg.OperationID)
	if err != nil {
		return errors.Wrap(err, "get contract operations")
	}

	switch len(operations) {
	case 0:
		return errors.Wrap(app.ErrImpossibleParsing, "unknown operation")
	case 1:
		return parseOperationAttempt(msg, operations[0])
	default:
		for _, op := range operations {
			switch op.ContractName {
			case known.NFTItem, known.NFTCollection, known.JettonMinter, known.JettonWallet:
				// firstly, skip standard contracts
			default:
				if err := parseOperationAttempt(msg, op); err == nil {
					return nil
				}
			}
		}
		var err error
		for _, op := range operations {
			if err = parseOperationAttempt(msg, op); err == nil {
				return nil
			}
		}
		return err
	}
}

func (s *Service) ParseMessagePayload(ctx context.Context, msg *core.Message) error {
	var err = app.ErrImpossibleParsing // save message parsing error to a database to look at it later

	// you can parse separately incoming messages to known contracts and outgoing message from them

	if len(msg.Body) == 0 {
		return errors.Wrap(app.ErrImpossibleParsing, "no message body")
	}

	errIn := s.parseDirectedMessage(ctx, msg.DstState, msg)
	if errIn != nil && !errors.Is(errIn, app.ErrImpossibleParsing) {
		err = errors.Wrap(errIn, "incoming")
	}
	if errIn == nil {
		return nil
	}

	errOut := s.parseDirectedMessage(ctx, msg.SrcState, msg)
	if errOut != nil && !errors.Is(errOut, app.ErrImpossibleParsing) {
		err = errors.Wrap(errOut, "outgoing")
	}
	if errOut == nil {
		return nil
	}

	return err
}
