package service

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/sudonite/TollCalculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
	log  log.Logger
}

func newLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
			log:  logger,
		}
	}
}

func (mw loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Aggregate", "took", time.Since(start), "obuID", dist.OBUID, "distance", dist.Value, "unix", dist.Unix, "err", err)
	}(time.Now())

	err = mw.next.Aggregate(ctx, dist)
	return err
}

func (mw loggingMiddleware) Calculate(ctx context.Context, obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Calculate", "took", time.Since(start), "obuID", obuID, "err", err)
	}(time.Now())

	inv, err = mw.next.Calculate(ctx, obuID)
	return inv, err
}

type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{
			next: next,
		}
	}
}

func (mw instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw instrumentationMiddleware) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, obuID)
}
