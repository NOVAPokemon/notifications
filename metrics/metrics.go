package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	connectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "notifications_connected_clients",
		Help: "The total number of connected clients",
	})

	sentNotificationsLocal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent_local",
		Help: "The total number of sent to local clients",
	})

	receivedNotificationsKafka = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_received_kafka",
		Help: "The total number of received notifications via kafka",
	})

	sentNotificationsKafka = promauto.NewCounter(prometheus.CounterOpts{
		Name: "notifications_sent_kafka",
		Help: "The total number of sent notifications via kafka",
	})
)

func EmitSentNotificationKafka() {
	sentNotificationsLocal.Inc()
}

func EmitSentNotificationLocal() {
	receivedNotificationsKafka.Inc()
}

func EmitReceivedNotificationKafka() {
	sentNotificationsKafka.Inc()
}

func emitNrConnectedClients(userChannels *sync.Map) {
	length := 0
	userChannels.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	connectedClients.Set(float64(length))
}

// metrics for prometheus
func RecordMetrics(userChannels *sync.Map) {
	go func() {
		for {
			emitNrConnectedClients(userChannels)
			time.Sleep(2 * time.Second)
		}
	}()
}
