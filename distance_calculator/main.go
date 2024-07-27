package main

import (
	"log"

	"github.com/sudonite/TollCalculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:3000"
)

func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)

	/*
		grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
		if err != nil {
			log.Fatal(err)
		}
	*/

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}
