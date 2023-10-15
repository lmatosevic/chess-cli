package command

import (
	"errors"
	"fmt"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func JoinGame(gameId int64, password string) (*model.GenericResponse, error) {
	resp, err := client.SendRequest[model.GenericResponse]("POST", fmt.Sprintf("/v1/games/%d/join", gameId), nil,
		&model.GameJoin{Password: password})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
