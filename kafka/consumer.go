package kafka

import (
	"context"
	"time"

	"github.com/NOVAPokemon/notifications/metrics"
	ws "github.com/NOVAPokemon/utils/websockets"
	"github.com/NOVAPokemon/utils/websockets/notifications"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// NotificationsConsumer kafka notification consumer
type NotificationsConsumer struct {
	KafkaUrl             string
	NotificationsChannel chan *ws.WebsocketMsg
	FinishChan           chan struct{}
	Username             string
}

// PipeMessagesFromTopic send messages received from topic to channel
func (nc *NotificationsConsumer) PipeMessagesFromTopic() {
	r := kafka.NewReader(kafka.ReaderConfig{
		MaxWait:   1 * time.Second,
		Brokers:   []string{nc.KafkaUrl},
		Topic:     nc.Username,
		Partition: 0,
		MinBytes:  1e3,
		MaxBytes:  10e6,
	})

	defer close(nc.NotificationsChannel)

LOOP:
	for {
		select {
		case <-nc.FinishChan:
			break LOOP
		default:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			before := ws.MakeTimestamp()

			m, err := r.ReadMessage(ctx)
			cancel()

			after := ws.MakeTimestamp()

			if err != nil {
				continue
			}

			wsMsgContent := ws.ParseContent(m.Value)
			notificationMsg := wsMsgContent.Data.(notifications.NotificationMessage)
			log.Infof("read notification %s from kafka: %d ms", notificationMsg.Notification.Id, after-before)

			metrics.EmitReceivedNotificationKafka()
			nc.NotificationsChannel <- notificationMsg.ConvertToWSMessage()
		}
	}

	if err := r.Close(); err != nil {
		log.Error(wrapConsumerError(err))
	}
	log.Warn("Kafka routine exiting...")
}

// Close stop reading messages from topic
func (nc *NotificationsConsumer) Close() {
	close(nc.FinishChan)
}
