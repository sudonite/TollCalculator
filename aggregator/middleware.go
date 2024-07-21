package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sudonite/TollCalculator/types"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Aggregated distance")
	}(time.Now())

	err = l.next.AggregateDistance(distance)
	return
}

func (l *LogMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)

		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalAmount
		}

		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"amount":   amount,
			"distance": distance,
		}).Info("Calculate Invoice")
	}(time.Now())

	inv, err = l.next.CalculateInvoice(obuID)
	return
}
