package main

import (
	"fmt"
	"strings"

	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
)

const (
	addNotificationName         = "ADD_NOTIFICATION"
	deleteNotificationName      = "DELETE_NOTIFICATION"
	subscribeNotificationName   = "SUBSCRIBE_NOTIFICATION"
	unsubscribeNotificationName = "UNSUBSCRIBE_NOTIFICATION"
	getListenersName            = "GET_LISTENERS"
)

const (
	get    = "GET"
	delete = "DELETE"
	post   = "POST"
)

var routes = utils.Routes{
	api.GenStatusRoute(strings.ToLower(fmt.Sprintf("%s", serviceName))),
	utils.Route{
		Name:        addNotificationName,
		Method:      post,
		Pattern:     api.NotificationPath,
		HandlerFunc: addNotificationHandler,
	},
	utils.Route{
		Name:        deleteNotificationName,
		Method:      delete,
		Pattern:     api.SpecificNotificationRoute,
		HandlerFunc: deleteNotificationHandler,
	},
	utils.Route{
		Name:        subscribeNotificationName,
		Method:      get,
		Pattern:     api.SubscribeNotificationPath,
		HandlerFunc: subscribeToNotificationsHandler,
	},
	utils.Route{
		Name:        getListenersName,
		Method:      get,
		Pattern:     api.GetListenersPath,
		HandlerFunc: getOtherListenersHandler,
	},
	utils.Route{
		Name:        unsubscribeNotificationName,
		Method:      get,
		Pattern:     api.UnsubscribeNotificationPath,
		HandlerFunc: unsubscribeToNotificationsHandler,
	},
}
