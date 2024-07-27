package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sudonite/TollCalculator/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(client client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: client,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	inv, err := h.client.GetInvoice(context.Background(), 605974)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())

		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "The listen address of the HTTP server")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:3000", "The listen address of the aggregator service")
	flag.Parse()

	logrus.Infof("Gateway HTTP server running on port %s", *listenAddr)

	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr)
		invHandler = NewInvoiceHandler(client)
	)

	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
