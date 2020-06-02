package kafka

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	ErrorInConsumerFormat = "error in consumer: %s"
	ErrorInProducerFormat = "error in producer: %s"
	ErrorInPublishMessage = "error publishing message"
)

func wrapConsumerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(ErrorInConsumerFormat, err))
}

func wrapProducerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(ErrorInProducerFormat, err))
}

func wrapPublishMessageError(err error) error {
	return errors.Wrap(err, ErrorInPublishMessage)
}
