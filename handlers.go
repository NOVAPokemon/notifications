package main

import (
	"encoding/json"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
	"github.com/NOVAPokemon/utils/tokens"
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
	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		return
	}

	var request AddNotificationRequest
	err = json.NewDecoder(r.Body).Decode(&request)
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

	channel, ok := userChannels.Get(username)
	if !ok {
		log.Errorf("user %s is not listening for notifications", username)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(notification)
	if err != nil {
		utils.HandleJSONEncodeError(&w, serviceName, err)
		return
	}

	log.Infof("got notification from %s to %s", claims.Username, username)
	channel <- jsonBytes
}

// Possibly useless endpoint since users dont need explicitly to delete
// the notifications read.
func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex, ok := vars[api.IdPathVar]
	if !ok {
		log.Error("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

func GetOtherListenersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username, ok := vars[api.UsernamePathVar]
	if !ok {
		log.Error("no id provided")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		log.Error(err)
		return
	}

	usernames := userChannels.GetOthers(username)

	log.Info("returning others: ", len(usernames))

	jsonBytes, err := json.Marshal(usernames)
	if err != nil {
		utils.HandleJSONEncodeError(&w, serviceName, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		handleError(&w, "Error writing json to body", err)
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

	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)
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
