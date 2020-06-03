package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	connectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "notifications_connected_clients",
		Help: "The total number of connected clients",
	})
)

func emitNrConnectedClients() {
	length := 0
	userChannels.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	connectedClients.Set(float64(length))
}

// metrics for prometheus
func recordMetrics() {
	go func() {
		for {
			emitNrConnectedClients()
			time.Sleep(2 * time.Second)
		}
	}()
}
