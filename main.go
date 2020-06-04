package main

import (
	"github.com/NOVAPokemon/notifications/metrics"
	"github.com/NOVAPokemon/utils"
)

const (
	host        = utils.ServeHost
	port        = utils.NotificationsPort
	serviceName = "MICROTRANSACTIONS"
)

func main() {
	metrics.RecordMetrics(&userChannels)
	utils.CheckLogFlag(serviceName)
	utils.StartServer(serviceName, host, port, routes)
}
