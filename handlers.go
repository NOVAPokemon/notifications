package main

import (
	"encoding/json"
	"fmt"
	"github.com/NOVAPokemon/notifications/kafka"
	"github.com/NOVAPokemon/utils"
	"github.com/NOVAPokemon/utils/api"
	"github.com/NOVAPokemon/utils/clients"
	notificationdb "github.com/NOVAPokemon/utils/database/notification"
	"github.com/NOVAPokemon/utils/tokens"
	ws "github.com/NOVAPokemon/utils/websockets"
	"github.com/NOVAPokemon/utils/websockets/notifications"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
	"sync"
	"time"
)

type keyType = string
type valueType = chan ws.Serializable

var userChannels = sync.Map{}
var kafkaUrl string

func init() {
	kafkaUrlAux, exists := os.LookupEnv(utils.KafkaEnvVar)
	if !exists {
		panic(fmt.Sprintf("missing: %s", utils.KafkaEnvVar))
	}
	kafkaUrl = kafkaUrlAux
}

func AddNotificationHandler(w http.ResponseWriter, r *http.Request) {
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

	notificationMsg.Notification.Id = primitive.NewObjectID()
	err = notificationdb.AddNotification(notificationMsg.Notification)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapAddNotificationError(err), http.StatusBadRequest)
		return
	}

	value, ok := userChannels.Load(username)
	if !ok {
		// user is not listening to this server
		before := ws.MakeTimestamp()
		producer := kafka.NotificationsProducer{
			Username: username,
			KafkaUrl: kafkaUrl,
		}
		err := producer.IssueOneNotification(notificationMsg)
		if err != nil {
			utils.LogAndSendHTTPError(&w, wrapAddNotificationError(err), http.StatusInternalServerError)
		}

		after := ws.MakeTimestamp()

		log.Infof("issue notification %s to kafka: %d ms", notificationMsg.Id, after-before)

		return
	}
	channel := value.(valueType)
	log.Infof("got notification from %s to %s", claims.Username, username)
	channel <- notificationMsg
}

// Possibly useless endpoint since users dont need explicitly to delete
// the notifications read.
func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
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

func GetOtherListenersHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapGetListenersError(err), http.StatusBadRequest)
		return
	}

	username := authToken.Username

	var usernames []string
	userChannels.Range(func(key, value interface{}) bool {
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

func SubscribeToNotificationsHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, ok := userChannels.Load(username); ok {
		err = wrapSubscribeNotificationError(newUserAlreadySubscribedError(username))
		utils.LogAndSendHTTPError(&w, err, http.StatusConflict)
		return
	}
	channel := make(chan ws.Serializable)
	userChannels.Store(username, channel)
	go handleUser(username, conn, channel)
}

func UnsubscribeToNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	authToken, err := tokens.ExtractAndVerifyAuthToken(r.Header)
	if err != nil {
		utils.LogAndSendHTTPError(&w, wrapUnsubscribeNotificationError(err), http.StatusBadRequest)
		return
	}

	log.Infof("unsubscribing %s from notifications", authToken.Username)

	if _, ok := userChannels.Load(authToken.Username); ok {
		userChannels.Delete(authToken.Username)
	}
}

func handleUser(username string, conn *websocket.Conn, channel chan ws.Serializable) {
	log.Info("handling user ", username)

	kafkaFinishChan := make(chan struct{})
	consumer := kafka.NotificationsConsumer{
		Username:             username,
		KafkaUrl:             kafkaUrl,
		FinishChan:           kafkaFinishChan,
		NotificationsChannel: channel,
	}
	go consumer.PipeMessagesFromTopic()
	ticker := time.NewTicker(ws.PingPeriod)
	defer closeUserListener(username, conn, channel, kafkaFinishChan, ticker)

	_ = conn.SetReadDeadline(time.Now().Add(ws.PongWait))
	conn.SetPongHandler(func(string) error {
		// log.Warn("Received pong")
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
			// log.Warn("Pinging")
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(wrapHandleUserError(err, username))
				return
			}
		case msg := <-channel:
			msgString := msg.SerializeToWSMessage().Serialize()
			err := clients.Send(conn, &msgString)
			if err != nil {
				log.Error(wrapHandleUserError(err, username))
				return
			}
		}
	}
}

func closeUserListener(username string, conn *websocket.Conn, channel chan ws.Serializable, kafkaFinishChan chan struct{}, ticker *time.Ticker) {
	log.Info("removing user ", username)
	ws.CloseConnection(conn)
	close(channel)
	close(kafkaFinishChan)

	if _, ok := userChannels.Load(username); ok {
		userChannels.Delete(username)
	}

	ticker.Stop()
}
