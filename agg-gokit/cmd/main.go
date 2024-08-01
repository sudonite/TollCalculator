package main

import (
	"net"
	"net/http"
	"os"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	"github.com/joho/godotenv"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/sudonite/TollCalculator/agg-gokit/endpoint"
	"github.com/sudonite/TollCalculator/agg-gokit/service"
	"github.com/sudonite/TollCalculator/agg-gokit/transport"
)

func main() {

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var duration metrics.Histogram
	{
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "tollcalculator",
			Subsystem: "aggservice",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method", "success"})
	}

	var (
		service     = service.New(logger)
		endpoints   = endpoint.New(service, logger, duration)
		httpHandler = transport.NewHTTPHandler(endpoints, logger)
		listenAddr  = os.Getenv("AGG_HTTP_ENDPOINT")
	)

	httpListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}

	logger.Log("transport", "HTTP", "addr", listenAddr)
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}
