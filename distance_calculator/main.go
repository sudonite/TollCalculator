package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sudonite/TollCalculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://localhost:4000"
)

func main() {
	var (
		aggregatorHTTPEndpoint = "http://localhost" + os.Getenv("AGG_HTTP_ENDPOINT")
		aggregatorGRPCEndpoint = "http://localhost" + os.Getenv("AGG_GRPC_ENDPOINT")
	)

	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorHTTPEndpoint)

	grpcClient, err := client.NewGRPCClient(aggregatorGRPCEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	_ = grpcClient

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer.Start()
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
