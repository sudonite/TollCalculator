package main

import (
	"log"

	"github.com/sudonite/TollCalculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000/aggregate"
)

func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
