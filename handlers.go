package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/cookies"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

const serviceName = "Notifications"

var userChannels = UserNotificationChannels{
	channels: map[string]chan []byte{},
}

func AddNotificationHandler(w http.ResponseWriter, r *http.Request) {
	var request AddNotificationRequest
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		utils.HandleJSONDecodeError(&w, serviceName, err)
	}

	id := primitive.NewObjectID()
	notification := utils.Notification{
		Id:       id,
		Username: request.Username,
		Type:     request.Type,
		Content:  request.Content,
	}

	err = notificationdb.AddNotification(notification)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	username := request.Username

	channel := userChannels.Get(username)

	if channel == nil {
		log.Errorf("user %s is not listening for notifications", username)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(notification)

	if err != nil {
		utils.HandleJSONEncodeError(&w, serviceName, err)
		return
	}

	log.Info("got notification %+v", notification)
	channel <- jsonBytes
}

// Possibly useless endpoint since users dont need explicitly to delete
// the notifications read.
func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex := vars[IdPathVar]
	id, err := primitive.ObjectIDFromHex(idHex)

	if err != nil {
		log.Error("bad id in delete request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = notificationdb.RemoveNotification(id)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
	}
}

func SubscribeToNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		handleError(&w, "Connection error", err)
		return
	}

	claims, err := cookies.ExtractAndVerifyAuthToken(&w, r, serviceName)

	if err != nil {
		return
	}

	username := claims.Username

	if userChannels.Has(username) {
		return
	}

	channel := make(chan []byte)
	userChannels.Add(username, channel)

	go notifyUser(username, conn, channel)
}

func notifyUser(username string, conn *websocket.Conn, channel chan []byte) {
	defer closeUserListener(username, conn, channel)

	for {
		select {
		case jsonNotification := <-channel:
			err := conn.WriteMessage(websocket.TextMessage, jsonNotification)
			if err != nil {
				return
			}
		}
	}
}

func closeUserListener(username string, conn *websocket.Conn, channel chan []byte) {
	conn.Close()
	close(channel)
	userChannels.Remove(username)
}

func handleError(w *http.ResponseWriter, errorString string, err error) {
	log.Error(err)
	http.Error(*w, errorString, http.StatusInternalServerError)
	return
}
