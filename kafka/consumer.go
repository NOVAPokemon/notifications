package kafka

import (
	"context"
	"encoding/json"
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
	NotificationsChannel chan ws.Serializable
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

			deserialized := &ws.Message{}
			err = json.Unmarshal(m.Value, deserialized)
			if err != nil {
				log.Error(wrapConsumerError(err))
				continue
			}
			msgStr := string(m.Value)
			msgParsed, err := ws.ParseMessage(msgStr)
			if err != nil {
				log.Error(wrapConsumerError(err))
				continue
			}
			msg, err := notifications.DeserializeNotificationMessage(msgParsed)
			if err != nil {
				log.Error(wrapConsumerError(err))
				continue
			}
			log.Infof("read notification %s from kafka: %d ms", msg.GetId(), after-before)

			metrics.EmitReceivedNotificationKafka()
			nc.NotificationsChannel <- msg
		}
	}

	if err := r.Close(); err != nil {
		log.Error(wrapProducerError(err))
	}
	log.Warn("Kafka routine exiting...")
}

// Close stop reading messages from topic
func (nc *NotificationsConsumer) Close() {
	close(nc.FinishChan)
}
