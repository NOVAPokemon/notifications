package main

import (
	"fmt"
	"github.com/NOVAPokemon/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const host = "localhost"
const Port = 8010

var addr = fmt.Sprintf("%s:%d", host, Port)

func main() {
	r := utils.NewRouter(routes)

	log.Infof("Starting NOTIFICATIONS server in port %d...\n", Port)
	log.Fatal(http.ListenAndServe(addr, r))
}
