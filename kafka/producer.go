package kafka

import (
	"context"
	"time"

	notificationMessages "github.com/NOVAPokemon/utils/websockets/notifications"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// NotificationsProducer kafka notifications producer
type NotificationsProducer struct {
	KafkaUrl string
	Username string
}

// IssueOneNotification send a notification
func (nc *NotificationsProducer) IssueOneNotification(notification *notificationMessages.NotificationMessage) error {
	retry := false

	for retry {
		w := kafka.NewWriter(kafka.WriterConfig{
			BatchSize:    1,
			BatchTimeout: 10 * time.Millisecond,
			Brokers:      []string{nc.KafkaUrl},
			Topic:        nc.Username,
			Balancer:     &kafka.LeastBytes{},
		})

		toSend := kafka.Message{
			Key:   []byte(nc.Username),
			Value: notification.ConvertToWSMessage().Content.Serialize(),
		}

		if err := w.WriteMessages(context.Background(), toSend); err != nil {
			err = wrapPublishMessageError(err)
			log.Warn(wrapProducerError(err))
			retry = true
		} else {
			retry = false
		}

		if err := w.Close(); err != nil {
			log.Error(wrapProducerError(err))
			return err
		}
	}

	return nil
}

func (nc *NotificationsProducer) pipeNotificationsFromChan(inChan chan notificationMessages.NotificationMessage,
	finishChan chan struct{}) error {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{nc.KafkaUrl},
		Topic:    nc.Username,
		Balancer: &kafka.LeastBytes{},
	})

LOOP:
	for {
		select {
		case <-finishChan:
			break LOOP
		case notification, ok := <-inChan:
			if !ok {
				continue
			}

			toSend := kafka.Message{
				Key:   []byte(nc.Username),
				Value: notification.ConvertToWSMessage().Content.Serialize(),
			}

			if err := w.WriteMessages(context.Background(), toSend); err != nil {
				err = wrapPublishMessageError(err)
				log.Error(wrapProducerError(err))
			}
		}
	}
	if err := w.Close(); err != nil {
		log.Error(wrapProducerError(err))
		return err
	}
	return nil
}
