package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
	"github.com/sudonite/TollCalculator/aggregator/client"
	"github.com/sudonite/TollCalculator/types"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, aggClient client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("Kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("Kafka Consume error: %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON Serialization error: %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("Calculation error: %s", err)
		}

		req := &types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().UnixNano(),
			ObuID: int32(data.OBUID),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.Errorf("Aggregate error: %s", err)
			continue
		}
	}
}
