package client

import (
	"context"

	"github.com/sudonite/TollCalculator/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   types.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := types.NewAggregatorClient(conn)

	return &GRPCClient{
		Endpoint: endpoint,
		client:   c,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, aggReq *types.AggregateRequest) error {
	_, err := c.client.Aggregate(context.Background(), aggReq)
	return err
}
