package command

import (
	"errors"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func Register(username string, password string, stateless bool) (*model.Player, error) {
	resp, err := client.SendRequest[model.Player]("POST", "/v1/players/register", nil,
		&model.PlayerRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	_, err = Login(username, password, stateless)
	if err != nil {
		return nil, err
	}

	return &resp.Data, nil
}
