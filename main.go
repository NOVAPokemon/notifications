package main

import (
	"github.com/NOVAPokemon/notifications/metrics"
	"github.com/NOVAPokemon/utils"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
)

const (
	host        = utils.ServeHost
	port        = utils.NotificationsPort
	serviceName = "MICROTRANSACTIONS"
)

func main() {
	metrics.RecordMetrics(&userChannelsMap)

	flags := utils.ParseFlags(serverName)

	if !*flags.LogToStdout {
		utils.SetLogFile(serverName)
	}

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		locationTag := utils.GetLocationTag(utils.DefaultLocationTagsFilename, serverName)
		commsManager = utils.CreateDefaultDelayedManager(locationTag, false)
	}

	notificationdb.InitNotificationDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
