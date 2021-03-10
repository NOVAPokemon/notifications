package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	http "github.com/bruno-anjos/archimedesHTTPClient"

	"github.com/NOVAPokemon/notifications/metrics"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
	"github.com/NOVAPokemon/utils/tokens"
	ws "github.com/NOVAPokemon/utils/websockets"
	"github.com/NOVAPokemon/utils/websockets/notifications"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	keyType      = string
	userChannels = struct {
		notificationChannel chan *ws.WebsocketMsg
		finishChannel       chan struct{}
	}
	valueType = userChannels
)

var (
	userChannelsMap = sync.Map{}
	serverName      string
	commsManager    ws.CommunicationManager
)

func init() {
	if aux, exists := os.LookupEnv(utils.HostnameEnvVar); exists {
		serverName = aux
	} else {
		log.Fatal("could not load server name")
	}
	log.Info("Server name: ", serverName)
}

func addNotificationHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddNotificationError(err), http.StatusBadRequest)
		return
	}

	var notificationMsg notifications.NotificationMessage
	err = json.NewDecoder(r.Body).Decode(&notificationMsg)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddNotificationError(err), http.StatusBadRequest)
		return
	}

	username := notificationMsg.Notification.Username

	err = notificationdb.AddNotification(notificationMsg.Notification)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddNotificationError(err), http.StatusBadRequest)
		return
	}

	value, ok := userChannelsMap.Load(username)
	if !ok {
		notFoundErr := errors.New(fmt.Sprintf("user %s should be listening in the only server", username))
		utils.LogWarnAndSendHTTPError(&w, notFoundErr, http.StatusNotFound)
		return
	}

	metrics.EmitSentNotificationLocal()
	channels := value.(valueType)
	log.Infof("got notification %s from %s to %s (%d)", notificationMsg.Notification.Id, claims.Username,
		username, ws.MakeTimestamp())

	select {
	case <-channels.finishChannel:
		utils.LogAndSendHTTPError(&w, newUserAlreadyLeft(username), http.StatusNotFound)
		return
	case channels.notificationChannel <- notificationMsg.ConvertToWSMessage():
	}
}

func deleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idHex, ok := vars[api.IdPathVar]
	if !ok {
		err := wrapDeleteNotificationError(errorNoNotificationId)
		utils.LogAndSendHTTPError(&w, err, http.StatusBadRequest)
		return
	}

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapDeleteNotificationError(err), http.StatusBadRequest)
		return
	}

	err = notificationdb.RemoveNotification(id)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapDeleteNotificationError(err), http.StatusNotFound)
	}
}

func getOtherListenersHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetListenersError(err), http.StatusBadRequest)
		return
	}

	username := authToken.Username

	var usernames []string
	userChannelsMap.Range(func(key, value interface{}) bool {
		currUsername := key.(keyType)
		if currUsername != username {
			usernames = append(usernames, currUsername)
		}
		return true
	})

	log.Info("returning others: ", len(usernames))

	jsonBytes, err := json.Marshal(usernames)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetListenersError(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonBytes)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetListenersError(err), http.StatusInternalServerError)
	}
}

func subscribeToNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		err = wrapSubscribeNotificationError(ws.WrapUpgradeConnectionError(err))
		utils.LogAndSendHTTPError(&w, err, http.StatusInternalServerError)
		return
	}

	claims, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapSubscribeNotificationError(err), http.StatusBadRequest)
		return
	}

	username := claims.Username
	if _, ok := userChannelsMap.Load(username); ok {
		err = wrapSubscribeNotificationError(newUserAlreadySubscribedError(username))
		utils.LogAndSendHTTPError(&w, err, http.StatusConflict)
		return
	}
	channels := userChannels{
		notificationChannel: make(chan *ws.WebsocketMsg, 5),
		finishChannel:       make(chan struct{}),
	}
	userChannelsMap.Store(username, channels)
	go handleUser(username, conn, channels, commsManager)
}

func unsubscribeToNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUnsubscribeNotificationError(err), http.StatusBadRequest)
		return
	}

	log.Infof("unsubscribing %s from notifications", authToken.Username)

	if _, ok := userChannelsMap.Load(authToken.Username); ok {
		userChannelsMap.Delete(authToken.Username)
	}
}

func handleUser(username string, conn *websocket.Conn, channels userChannels, writer ws.CommunicationManager) {
	log.Info("handling user ", username)

	ticker := time.NewTicker(ws.PingPeriod)
	defer closeUserListener(username, conn, ticker)

	_ = conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	})

	go func() {
		if _, _, err := conn.NextReader(); err != nil {
			return
		}
	}()

	for {
		select {
		case <-ticker.C:
			err := writer.WriteGenericMessageToConn(conn, ws.NewControlMsg(websocket.PingMessage))
			if err != nil {
				log.Warn(wrapHandleUserError(err, username))
				return
			}
		case msg := <-channels.notificationChannel:
			err := clients.Send(conn, msg, writer)
			if err != nil {
				log.Warn(wrapHandleUserError(err, username))
				return
			}
		}
	}
}

func closeUserListener(username string, conn *websocket.Conn, ticker *time.Ticker) {
	log.Info("removing user ", username)
	if err := conn.Close(); err != nil {
		log.Warn(err)
	}

	if _, ok := userChannelsMap.Load(username); ok {
		userChannelsMap.Delete(username)
	}

	ticker.Stop()
}
