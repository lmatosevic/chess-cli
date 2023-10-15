package repository

import (
	"database/sql"
	"fmt"
	"gitlab.com/lmatosevic/chess-cli/pkg/database"
	"gitlab.com/lmatosevic/chess-cli/pkg/utils"
	"time"
)

type GameMove struct {
	Id        int64
	GameId    int64
	PlayerId  sql.NullInt64
	Move      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (gm *GameMove) FormatCreatedAt() string {
	return utils.ISODate(gm.CreatedAt)
}

func (gm *GameMove) FormatUpdatedAt() string {
	return utils.ISODate(gm.UpdatedAt)
}

func CreateGameMove(gameId int64, playerId int64, move string) error {
	_, err := database.GetConnection().Exec(
		`INSERT INTO game_move ("gameId", "playerId", "move") VALUES ($1, $2, $3)`, gameId, playerId, move)
	return err
}

func QueryGameMoves(filter string, page int, size int, sort string) (*[]GameMove, error) {
	where, sort, order, args := PrepareQueryParams(filter, page, size, sort)
	rows, err := database.GetConnection().Query(
		fmt.Sprintf(`SELECT * FROM game_move %s ORDER BY "%s" %s NULLS LAST LIMIT $%d OFFSET $%d`, where, sort, order,
			len(args)-1, len(args)), args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	gameMoves := make([]GameMove, 0)

	for rows.Next() {
		p := GameMove{}
		err := scanGameMoveRows(rows, &p)
		if err != nil {
			return nil, err
		}
		gameMoves = append(gameMoves, p)
	}

	return &gameMoves, nil
}

func CountGameMoves(filter string) (int, error) {
	where, _, _, args := PrepareQueryParams(filter, 0, 0, "")
	row := database.GetConnection().QueryRow(fmt.Sprintf(`SELECT count(*) FROM game_move %s`, where), args[:len(args)-2]...)

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func scanGameMoveRows(rows *sql.Rows, gm *GameMove) error {
	return rows.Scan(&gm.Id, &gm.GameId, &gm.PlayerId, &gm.Move, &gm.CreatedAt, &gm.UpdatedAt)
}
