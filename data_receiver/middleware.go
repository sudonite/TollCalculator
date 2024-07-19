package main

import "github.com/sudonite/TollCalculator/types"

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	return l.next.ProduceData(data)
}
