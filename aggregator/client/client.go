package client

import (
	"context"

	"github.com/sudonite/TollCalculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
