package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
)

const AddNotificationName = "ADD_NOTIFICATION"
const DeleteNotificationName = "DELETE_NOTIFICATION"
const SubscribeNotificationName = "SUBSCRIBE_NOTIFICATION"

const GET = "GET"
const DELETE = "DELETE"

const IdPathVar = "id"

const POST = "POST"

const NotificationPath = "/notification"
var SpecificNotificationPath = fmt.Sprintf("/notification/{%s}", IdPathVar)
const SubscribeNotificationPath = "/subscribe"

var routes = utils.Routes{
	utils.Route{
		Name:        AddNotificationName,
		Method:      POST,
		Pattern:     NotificationPath,
		HandlerFunc: AddNotificationHandler,
	},
	utils.Route{
		Name:        DeleteNotificationName,
		Method:      DELETE,
		Pattern:     SpecificNotificationPath,
		HandlerFunc: DeleteNotificationHandler,
	},
	utils.Route{
		Name:        SubscribeNotificationName,
		Method:      GET,
		Pattern:     SubscribeNotificationPath,
		HandlerFunc: SubscribeToNotificationsHandler,
	},
}
