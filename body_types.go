package main

type AddNotificationRequest struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	Content  []byte `json:"content"`
}
