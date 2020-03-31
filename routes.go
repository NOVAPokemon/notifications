package main

import (
	"github.com/NOVAPokemon/notifications/exported"
	"github.com/NOVAPokemon/utils"
)

const AddNotificationName = "ADD_NOTIFICATION"
const DeleteNotificationName = "DELETE_NOTIFICATION"
const SubscribeNotificationName = "SUBSCRIBE_NOTIFICATION"

const GET = "GET"
const DELETE = "DELETE"

const POST = "POST"

var routes = utils.Routes{
	utils.Route{
		Name:        AddNotificationName,
		Method:      POST,
		Pattern:     exported.NotificationPath,
		HandlerFunc: AddNotificationHandler,
	},
	utils.Route{
		Name:        DeleteNotificationName,
		Method:      DELETE,
		Pattern:     exported.SpecificNotificationRoute,
		HandlerFunc: DeleteNotificationHandler,
	},
	utils.Route{
		Name:        SubscribeNotificationName,
		Method:      GET,
		Pattern:     exported.SubscribeNotificationPath,
		HandlerFunc: SubscribeToNotificationsHandler,
	},
}
