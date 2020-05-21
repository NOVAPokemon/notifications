package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorUserAlreadySubscribedFormat = "user %s is already subscribed"
	errorHandleUserFormat            = "error handling user %s"
)

var (
	errorNoNotificationId = errors.New("no notification id provided")
)

// Handler wrappers
func wrapAddNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, AddNotificationName))
}

func wrapDeleteNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, DeleteNotificationName))
}

func wrapGetListenersError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, GetListenersName))
}

func wrapSubscribeNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, SubscribeNotificationName))
}

func wrapUnsubscribeNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, UnsubscribeNotificationName))
}


// Other wrappers
func wrapHandleUserError(err error, username string) error {
	return errors.Wrap(err, fmt.Sprintf(errorHandleUserFormat, username))
}

func newUserAlreadySubscribedError(username string) error {
	return errors.New(fmt.Sprintf(errorUserAlreadySubscribedFormat, username))
}
