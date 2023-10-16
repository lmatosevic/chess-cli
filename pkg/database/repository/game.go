package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/database"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Game struct {
	Id                  int64
	Name                string
	PasswordHash        sql.NullString
	TurnDurationSeconds sql.NullInt32
	WhitePlayerId       sql.NullInt64
	BlackPlayerId       sql.NullInt64
	CreatorId           sql.NullInt64
	WinnerId            sql.NullInt64
	Tiles               string
	InProgress          bool
	LastMovePlayedAt    sql.NullTime
	StartedAt           sql.NullTime
	EndedAt             sql.NullTime
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (g *Game) FormatLastMovePlayedAt() string {
	if g.LastMovePlayedAt.Valid {
		return utils.ISODate(g.LastMovePlayedAt.Time)
	} else {
		return ""
	}
}

func (g *Game) FormatStartedAt() string {
	if g.StartedAt.Valid {
		return utils.ISODate(g.StartedAt.Time)
	} else {
		return ""
	}
}

func (g *Game) FormatEndedAt() string {
	if g.EndedAt.Valid {
		return utils.ISODate(g.EndedAt.Time)
	} else {
		return ""
	}
}

func (g *Game) FormatCreatedAt() string {
	return utils.ISODate(g.CreatedAt)
}

func (g *Game) FormatUpdatedAt() string {
	return utils.ISODate(g.UpdatedAt)
}

func FindGameById(id int64) (*Game, error) {
	rows, err := database.GetConnection().Query(
		`SELECT * FROM game WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	g := Game{}

	for rows.Next() {
		err := scanGameRows(rows, &g)
		if err != nil {
			return nil, err
		}
	}

	if g.Id == 0 {
		return nil, errors.New("game does not exist")
	}

	return &g, nil
}

func QueryGames(filter string, page int, size int, sort string) (*[]Game, error) {
	where, sort, order, args := PrepareQueryParams(filter, page, size, sort)
	rows, err := database.GetConnection().Query(
		fmt.Sprintf(`SELECT * FROM game %s ORDER BY "%s" %s NULLS LAST LIMIT $%d OFFSET $%d`, where, sort, order,
			len(args)-1, len(args)), args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	games := make([]Game, 0)

	for rows.Next() {
		g := Game{}
		err := scanGameRows(rows, &g)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	return &games, nil
}

func CountGames(filter string) (int, error) {
	where, _, _, args := PrepareQueryParams(filter, 0, 0, "")
	row := database.GetConnection().QueryRow(fmt.Sprintf(`SELECT count(*) FROM game %s`, where), args[:len(args)-2]...)

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func CreateGame(name string, password string, turnDurationSeconds int32, creatorId int64, white bool, tiles string) (*Game, error) {
	var passwordHash sql.NullString
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
		if err != nil {
			return nil, err
		}
		passwordHash = sql.NullString{String: string(hash), Valid: true}
	}

	blackPlayerId := sql.NullInt64{}
	whitePlayerId := sql.NullInt64{}
	if white {
		whitePlayerId = sql.NullInt64{Int64: creatorId, Valid: true}
	} else {
		blackPlayerId = sql.NullInt64{Int64: creatorId, Valid: true}
	}

	turnDuration := sql.NullInt32{}
	if turnDurationSeconds > 0 {
		turnDuration = sql.NullInt32{Int32: turnDurationSeconds, Valid: true}
	}

	row := database.GetConnection().QueryRow(
		`INSERT INTO game ("name", "passwordHash", "turnDurationSeconds", "tiles", "whitePlayerId", "blackPlayerId", "creatorId") VALUES 
    	 ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, name, passwordHash, turnDuration, tiles, whitePlayerId, blackPlayerId, creatorId)

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return FindGameById(id)
}

func UpdateGame(game *Game) error {
	res, err := database.GetConnection().Exec(`UPDATE game SET "name" = $2, "passwordHash" = $3, "turnDurationSeconds" = $4, 
                "whitePlayerId" = $5, "blackPlayerId" = $6, "creatorId" = $7, "winnerId" = $8, "tiles" = $9, 
                "inProgress" = $10, "lastMovePlayedAt" = $11, "startedAt" = $12, "endedAt" = $13, "updatedAt" = $14 WHERE id = $1`,
		game.Id, game.Name, game.PasswordHash, game.TurnDurationSeconds, game.WhitePlayerId, game.BlackPlayerId, game.CreatorId,
		game.WinnerId, game.Tiles, game.InProgress, SqlDateFormat(game.LastMovePlayedAt), SqlDateFormat(game.StartedAt),
		SqlDateFormat(game.EndedAt), utils.ISODateNow())
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("game does not exist")
	}
	return err
}

func DeleteGame(id int64) error {
	res, err := database.GetConnection().Exec(`DELETE FROM game WHERE id = $1`, id)
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("game does not exist")
	}
	return err
}

func FindInactiveGames() (*[]Game, error) {
	rows, err := database.GetConnection().Query(`SELECT * FROM game WHERE "turnDurationSeconds" IS NOT NULL
  AND "inProgress" IS TRUE AND (("lastMovePlayedAt" IS NOT NULL AND
        "lastMovePlayedAt" <= (now() at time zone 'utc') - concat("turnDurationSeconds"::text, ' seconds')::interval)
    OR ("lastMovePlayedAt" IS NULL AND "startedAt" <= (now() at time zone 'utc') - concat("turnDurationSeconds"::text, ' seconds')::interval))`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	games := make([]Game, 0)

	for rows.Next() {
		g := Game{}
		err := scanGameRows(rows, &g)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	return &games, nil
}

func scanGameRows(rows *sql.Rows, g *Game) error {
	return rows.Scan(&g.Id, &g.Name, &g.PasswordHash, &g.TurnDurationSeconds, &g.WhitePlayerId, &g.BlackPlayerId,
		&g.CreatorId, &g.WinnerId, &g.Tiles, &g.InProgress, &g.LastMovePlayedAt, &g.StartedAt, &g.EndedAt, &g.CreatedAt,
		&g.UpdatedAt)
}
