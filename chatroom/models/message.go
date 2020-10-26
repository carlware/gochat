package models

import "time"

type Message struct {
	ID      string    `json:"id"`
	RID     string    `json:"rid"`
	UID     string    `json:"uid"`
	Created time.Time `json:"created"`
	Message string    `json:"message"`
}
