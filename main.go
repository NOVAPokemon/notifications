package main

import (
	"github.com/NOVAPokemon/utils"
)

const (
	host = utils.ServeHost
	port = utils.NotificationsPort
	serviceName = "MICROTRANSACTIONS"
)

func main() {
	utils.CheckLogFlag(serviceName)
	utils.StartServer(serviceName, host, port, routes)
}
