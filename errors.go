package main

import (
	"fmt"

	"github.com/NOVAPokemon/utils"
	"github.com/pkg/errors"
)

const (
	errorUserAlreadySubscribedFormat = "user %s is already subscribed"
	errorHandleUserFormat            = "error handling user %s"
	errorUserLeft                    = "user %s is not listening anymore"
)

var (
	errorNoNotificationId = errors.New("no notification id provided")
)

// Handler wrappers
func wrapAddNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, addNotificationName))
}

func wrapDeleteNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, deleteNotificationName))
}

func wrapGetListenersError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, getListenersName))
}

func wrapSubscribeNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, subscribeNotificationName))
}

func wrapUnsubscribeNotificationError(err error) error {
	return errors.Wrap(err, fmt.Sprintf(utils.ErrorInHandlerFormat, unsubscribeNotificationName))
}

// Other wrappers
func wrapHandleUserError(err error, username string) error {
	return errors.Wrap(err, fmt.Sprintf(errorHandleUserFormat, username))
}

func newUserAlreadySubscribedError(username string) error {
	return errors.New(fmt.Sprintf(errorUserAlreadySubscribedFormat, username))
}

func newUserAlreadyLeft(username string) error {
	return errors.New(fmt.Sprintf(errorUserLeft, username))
}
