package model

type GameMove struct {
	Id        int64  `json:"id"`
	GameId    int64  `json:"gameId"`
	PlayerId  int64  `json:"playerId"`
	Move      string `json:"move"`
	CreatedAt string `json:"createdAt"`
}
