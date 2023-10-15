package command

import (
	"errors"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func UserInfo() (*model.Player, error) {
	resp, err := client.SendRequest[model.Player]("GET", "/v1/auth/player", nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
