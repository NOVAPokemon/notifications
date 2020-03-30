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

func (n *UserNotificationChannels) Get(username string) chan []byte {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()
	return n.channels[username]
}

func (n *UserNotificationChannels) Has(username string) bool {
	n.RWMutex.RLock()
	defer n.RWMutex.RUnlock()
	return n.channels[username] != nil
}
