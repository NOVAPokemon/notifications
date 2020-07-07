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

// EmitSentNotificationKafka increment the number of notifications sent
func EmitSentNotificationKafka() {
	sentNotificationsLocal.Inc()
}

// EmitSentNotificationLocal increment the number of notifications sent to connected users
func EmitSentNotificationLocal() {
	receivedNotificationsKafka.Inc()
}

// EmitReceivedNotificationKafka increment the number of received notifications
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

// RecordMetrics start routine to record metrics for prometheus
func RecordMetrics(userChannels *sync.Map) {
	go func() {
		for {
			emitNrConnectedClients(userChannels)
			time.Sleep(2 * time.Second)
		}
	}()
}
