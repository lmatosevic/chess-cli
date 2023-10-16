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

type Player struct {
	Id           int64
	Username     string
	PasswordHash string
	Wins         int32
	Losses       int32
	Draws        int32
	Rate         float32
	Elo          int32
	LastPlayedAt sql.NullTime
	IsPlaying    bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (p *Player) FormatLastPlayedAt() string {
	if p.LastPlayedAt.Valid {
		return utils.ISODate(p.LastPlayedAt.Time)
	} else {
		return ""
	}
}

func (p *Player) FormatCreatedAt() string {
	return utils.ISODate(p.CreatedAt)
}

func (p *Player) FormatUpdatedAt() string {
	return utils.ISODate(p.UpdatedAt)
}

func (p *Player) RefreshIsPlaying() {
	gameCount, _ := CountGames(fmt.Sprintf("whitePlayerId=%d;and;inProgress=true;or;blackPlayerId=%d;and;inProgress=true", p.Id, p.Id))
	if gameCount > 0 {
		p.IsPlaying = true
	} else {
		p.IsPlaying = false
	}
}

func FindPlayerById(id int64) (*Player, error) {
	rows, err := database.GetConnection().Query(
		`SELECT * FROM player WHERE id = $1 LIMIT 1`, id)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	p := Player{}

	for rows.Next() {
		err := scanPlayerRows(rows, &p)
		if err != nil {
			return nil, err
		}
	}

	if p.Id == 0 {
		return nil, errors.New("player does not exist")
	}

	return &p, nil
}

func FindPlayerByUsername(username string) (*Player, error) {
	rows, err := database.GetConnection().Query(
		`SELECT * FROM player WHERE username = $1 LIMIT 1`, username)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	p := Player{}

	for rows.Next() {
		err := scanPlayerRows(rows, &p)
		if err != nil {
			return nil, err
		}
	}

	if p.Id == 0 {
		return nil, errors.New("player does not exist")
	}

	return &p, nil
}

func QueryPlayers(filter string, page int, size int, sort string) (*[]Player, error) {
	where, sort, order, args := PrepareQueryParams(filter, page, size, sort)
	rows, err := database.GetConnection().Query(
		fmt.Sprintf(`SELECT * FROM player %s ORDER BY "%s" %s NULLS LAST LIMIT $%d OFFSET $%d`, where, sort, order,
			len(args)-1, len(args)), args...)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	players := make([]Player, 0)

	for rows.Next() {
		p := Player{}
		err := scanPlayerRows(rows, &p)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return &players, nil
}

func CountPlayers(filter string) (int, error) {
	where, _, _, args := PrepareQueryParams(filter, 0, 0, "")
	row := database.GetConnection().QueryRow(fmt.Sprintf(`SELECT count(*) FROM player %s`, where), args[:len(args)-2]...)

	var totalCount int
	err := row.Scan(&totalCount)
	if err != nil {
		return 0, err
	}

	return totalCount, nil
}

func CreatePlayer(username string, password string) (*Player, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return nil, err
	}

	row := database.GetConnection().QueryRow(
		`INSERT INTO player ("username", "passwordHash") VALUES ($1, $2) RETURNING id`, username, string(passwordHash))

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return FindPlayerById(id)
}

func UpdatePlayer(player *Player) error {
	res, err := database.GetConnection().Exec(`UPDATE player SET "passwordHash" = $2, "wins" = $3, "losses" = $4, 
                  "draws" = $5, "rate" = $6, "elo" = $7, "lastPlayedAt" = $8, "isPlaying" = $9, "updatedAt" = $10 
              WHERE id = $1`,
		player.Id, player.PasswordHash, player.Wins, player.Losses, player.Draws, player.Rate, player.Elo,
		SqlDateFormat(player.LastPlayedAt), player.IsPlaying, utils.ISODateNow())
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("player does not exist")
	}
	return err
}

func DeletePlayer(id int64) error {
	res, err := database.GetConnection().Exec(`DELETE FROM player WHERE id = $1`, id)
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("player does not exist")
	}
	return err
}

func scanPlayerRows(rows *sql.Rows, p *Player) error {
	return rows.Scan(&p.Id, &p.Username, &p.PasswordHash, &p.Wins, &p.Losses, &p.Draws, &p.Rate, &p.Elo,
		&p.LastPlayedAt, &p.CreatedAt, &p.UpdatedAt, &p.IsPlaying)
}
