package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = utils.ServeHost
const port = utils.NotificationsPort

var addr = fmt.Sprintf("%s:%d", host, port)

func main() {
	r := utils.NewRouter(routes)

	log.Infof("Starting NOTIFICATIONS server in port %d...\n", port)
	log.Fatal(http.ListenAndServe(addr, r))
}
