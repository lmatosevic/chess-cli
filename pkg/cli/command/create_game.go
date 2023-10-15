package command

import (
	"errors"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func CreateGame(name string, password string, turnDuration int32, isWhite bool) (*model.Game, error) {
	resp, err := client.SendRequest[model.Game]("POST", "/v1/games/create", nil,
		&model.GameCreate{Name: name, Password: password, TurnDurationSeconds: turnDuration, IsWhite: isWhite})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
