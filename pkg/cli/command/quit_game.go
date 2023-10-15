package command

import (
	"errors"
	"fmt"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func QuitGame(gameId int64) (*model.GenericResponse, error) {
	resp, err := client.SendRequest[model.GenericResponse]("POST", fmt.Sprintf("/v1/games/%d/quit", gameId), nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
