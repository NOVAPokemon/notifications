package main

import (
	"sync"
)

type UserNotificationChannels struct {
	sync.RWMutex
	channels map[string]chan []byte
}

func (n *UserNotificationChannels) Add(username string, channel chan []byte) {
	n.RWMutex.Lock()
	defer n.RWMutex.Unlock()
	n.channels[username] = channel
}

func (n *UserNotificationChannels) Remove(username string) {
	n.RWMutex.Lock()
	defer n.RWMutex.Unlock()
	delete(n.channels, username)
}

func (n *UserNotificationChannels) Get(username string) (channel chan []byte, ok bool) {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()

	channel, ok = n.channels[username]
	return channel, ok
}

func (n *UserNotificationChannels) Has(username string) bool {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()
	return n.channels[username] != nil
}

func (n *UserNotificationChannels) GetOthers(myUsername string) []string {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()

	var usernames []string
	for username, _ := range n.channels {
		if username != myUsername{
			usernames = append(usernames, username)
		}
	}

	return usernames
}
