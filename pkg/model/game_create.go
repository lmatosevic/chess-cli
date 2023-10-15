package model

type GameCreate struct {
	Name                string `json:"name"`
	Password            string `json:"password"`
	TurnDurationSeconds int32  `json:"turnDurationSeconds"`
	IsWhite             bool   `json:"isWhite"`
}
