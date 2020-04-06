package main

import (
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const AddNotificationName = "ADD_NOTIFICATION"
const DeleteNotificationName = "DELETE_NOTIFICATION"
const SubscribeNotificationName = "SUBSCRIBE_NOTIFICATION"
const UnsubscribeNotificationName = "UNSUBSCRIBE_NOTIFICATION"
const GetListenersName = "GET_LISTENERS"

const GET = "GET"
const DELETE = "DELETE"
const POST = "POST"

var routes = utils.Routes{
	utils.Route{
		Name:        AddNotificationName,
		Method:      POST,
		Pattern:     api.NotificationPath,
		HandlerFunc: AddNotificationHandler,
	},
	utils.Route{
		Name:        DeleteNotificationName,
		Method:      DELETE,
		Pattern:     api.SpecificNotificationRoute,
		HandlerFunc: DeleteNotificationHandler,
	},
	utils.Route{
		Name:        SubscribeNotificationName,
		Method:      GET,
		Pattern:     api.SubscribeNotificationPath,
		HandlerFunc: SubscribeToNotificationsHandler,
	},
	utils.Route{
		Name:        GetListenersName,
		Method:      GET,
		Pattern:     api.GetListenersPath,
		HandlerFunc: GetOtherListenersHandler,
	},
	utils.Route{
		Name:        UnsubscribeNotificationName,
		Method:      GET,
		Pattern:     api.UnsubscribeNotificationPath,
		HandlerFunc: UnsubscribeToNotificationsHandler,
	},
}
