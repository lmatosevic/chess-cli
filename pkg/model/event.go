package model

type Event struct {
	Type      string    `json:"type"`
	Data      EventData `json:"data"`
	Timestamp string    `json:"timestamp"`
}
