package kafka

import (
	"context"
	notificationMessages "github.com/NOVAPokemon/utils/websockets/notifications"
	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type NotificationsProducer struct {
	KafkaUrl string
	Username string
}

func (nc *NotificationsProducer) IssueOneNotification(notification notificationMessages.NotificationMessage) error {

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{nc.KafkaUrl},
		Topic:    nc.Username,
		Balancer: &kafka.LeastBytes{},
	})

	toSend := kafka.Message{
		Key:   []byte(nc.Username),
		Value: []byte(notification.SerializeToWSMessage().Serialize()),
	}

	if err := w.WriteMessages(context.Background(), toSend); err != nil {
		err = wrapPublishMessageError(err)
		log.Error(wrapProducerError(err))
	}
	if err := w.Close(); err != nil {
		log.Error(wrapProducerError(err))
		return err
	}
	return nil
}

func (nc *NotificationsProducer) PipeNotificationsFromChan(inChan chan notificationMessages.NotificationMessage, finishChan chan struct{}) error {

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
				Value: []byte(notification.SerializeToWSMessage().Serialize()),
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
