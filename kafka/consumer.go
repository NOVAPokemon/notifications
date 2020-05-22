package kafka

import (
	"context"
	"encoding/json"
	ws "github.com/NOVAPokemon/utils/websockets"
	"github.com/NOVAPokemon/utils/websockets/notifications"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type NotificationsConsumer struct {
	KafkaUrl             string
	NotificationsChannel chan ws.Serializable
	FinishChan           chan struct{}
	Username             string
}

func (nc *NotificationsConsumer) PipeMessagesFromTopic() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{nc.KafkaUrl},
		Topic:     nc.Username,
		Partition: 0,
	})

LOOP:
	for {
		select {
		case <-nc.FinishChan:
			break LOOP
		default:
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			m, err := r.ReadMessage(ctx)
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
			msgParsed, err := ws.ParseMessage(&msgStr)
			if err != nil {
				log.Error(wrapConsumerError(err))
				continue
			}
			msg, err := notifications.DeserializeNotificationMessage(msgParsed)
			if err != nil {
				log.Error(wrapConsumerError(err))
				continue
			}
			nc.NotificationsChannel <- msg
		}
	}

	if err := r.Close(); err != nil {
		log.Error(wrapProducerError(err))
	}
	log.Warn("Kafka routine exiting...")
}
