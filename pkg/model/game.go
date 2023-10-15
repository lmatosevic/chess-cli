package model

type Game struct {
	Id                  int64  `json:"id"`
	Name                string `json:"name"`
	TurnDurationSeconds int32  `json:"turnDurationSeconds"`
	Public              bool   `json:"public"`
	WhitePlayerId       int64  `json:"whitePlayerId"`
	BlackPlayerId       int64  `json:"blackPlayerId"`
	WinnerId            int64  `json:"winnerId"`
	CreatorId           int64  `json:"creatorId"`
	InProgress          bool   `json:"inProgress"`
	Tiles               string `json:"tiles"`
	LastMovePlayedAt    string `json:"lastMovePlayedAt"`
	StartedAt           string `json:"startedAt"`
	EndedAt             string `json:"endedAt"`
	CreatedAt           string `json:"createdAt"`
}

type GameListResponse ListResponse[Game]

type GameMoveListResponse ListResponse[GameMove]
