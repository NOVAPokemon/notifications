package main

import (
	"github.com/NOVAPokemon/notifications/metrics"
	"github.com/NOVAPokemon/utils"
)

const (
	host        = utils.ServeHost
	port        = utils.NotificationsPort
	serviceName = "NOTIFICATIONS"
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

	utils.StartServer(serviceName, host, port, routes, commsManager)
}
