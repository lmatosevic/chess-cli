package model

type Player struct {
	Id           int64   `json:"id"`
	Username     string  `json:"username"`
	Wins         int32   `json:"wins"`
	Losses       int32   `json:"losses"`
	Draws        int32   `json:"draws"`
	Rate         float32 `json:"rate"`
	Elo          int32   `json:"elo"`
	IsPlaying    bool    `json:"isPlaying"`
	LastPlayedAt string  `json:"lastPlayedAt"`
	CreatedAt    string  `json:"createdAt"`
}

type PlayerListResponse ListResponse[Player]
