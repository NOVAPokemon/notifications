module github.com/NOVAPokemon/notifications

go 1.13

require (
	github.com/NOVAPokemon/authentication v0.0.7
	github.com/NOVAPokemon/utils v0.0.62
	github.com/gorilla/mux v1.7.4
	github.com/gorilla/websocket v1.4.2
	github.com/sirupsen/logrus v1.4.2
	go.mongodb.org/mongo-driver v1.3.1
)

replace github.com/NOVAPokemon/utils v0.0.62 => ../utils
