package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lmatosevic/chess-cli/pkg/database"
	"time"
)

type AccessToken struct {
	Id        int64
	PlayerId  int64
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateAccessToken(playerId int64) (*AccessToken, error) {
	token := uuid.New().String()

	_, err := database.GetConnection().Exec(
		`INSERT INTO access_token ("playerId", "token") VALUES ($1, $2)`, playerId, token)
	if err != nil {
		return nil, err
	}

	return FindAccessToken(token)
}

func FindAccessToken(token string) (*AccessToken, error) {
	rows, err := database.GetConnection().Query(
		`SELECT * FROM access_token WHERE token = $1 LIMIT 1`, token)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	at := AccessToken{}

	for rows.Next() {
		err := scanAccessTokenRows(rows, &at)
		if err != nil {
			return nil, err
		}
	}

	if at.Id == 0 {
		return nil, errors.New("access token does not exist")
	}

	return &at, nil
}

func RevokeAccessToken(token string) error {
	res, err := database.GetConnection().Exec(`DELETE FROM access_token WHERE token = $1`, token)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("access token does not exist")
	}

	return nil
}

func scanAccessTokenRows(rows *sql.Rows, at *AccessToken) error {
	return rows.Scan(&at.Id, &at.PlayerId, &at.Token, &at.CreatedAt, &at.UpdatedAt)
}
