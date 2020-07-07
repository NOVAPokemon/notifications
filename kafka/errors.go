package kafka

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	errorInConsumerFormat = "error in consumer: %s"
	errorInProducerFormat = "error in producer: %s"
	errorInPublishMessage = "error publishing message"
)

func wrapConsumerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(errorInConsumerFormat, err))
}

func wrapProducerError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(errorInProducerFormat, err))
}

func wrapPublishMessageError(err error) error {
	return errors.Wrap(err, errorInPublishMessage)
}
