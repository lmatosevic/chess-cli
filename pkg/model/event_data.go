package model

type EventData struct {
	GameId   int64  `json:"gameId"`
	PlayerId int64  `json:"playerId"`
	Payload  string `json:"payload"`
}
