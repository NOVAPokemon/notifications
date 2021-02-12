package main

import (
	"os"

	"github.com/NOVAPokemon/notifications/metrics"
	"github.com/NOVAPokemon/utils"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
	"github.com/golang/geo/s2"
	log "github.com/sirupsen/logrus"
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

	location, exists := os.LookupEnv("LOCATION")
	if !exists {
		log.Fatal("no location in environment")
	}

	cellID := s2.CellIDFromToken(location)

	if !*flags.DelayedComms {
		commsManager = utils.CreateDefaultCommunicationManager()
	} else {
		commsManager = utils.CreateDefaultDelayedManager(false, &utils.OptionalConfigs{
			CellID: cellID,
		})
	}

	notificationdb.InitNotificationDBClient(*flags.ArchimedesEnabled)
	utils.StartServer(serviceName, host, port, routes, commsManager)
}
