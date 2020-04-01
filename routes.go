package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const AddNotificationName = "ADD_NOTIFICATION"
const DeleteNotificationName = "DELETE_NOTIFICATION"
const SubscribeNotificationName = "SUBSCRIBE_NOTIFICATION"

const GET = "GET"
const DELETE = "DELETE"
const POST = "POST"

const IdPathVar = "id"

var SpecificNotificationRoute = fmt.Sprintf("/notification/{%s}", IdPathVar)

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
		Pattern:     SpecificNotificationRoute,
		HandlerFunc: DeleteNotificationHandler,
	},
	utils.Route{
		Name:        SubscribeNotificationName,
		Method:      GET,
		Pattern:     api.SubscribeNotificationPath,
		HandlerFunc: SubscribeToNotificationsHandler,
	},
}
