package command

import (
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
)

func PlayGame(gameId int64, move string) (*model.GenericResponse, error) {
	resp, err := client.SendRequest[model.GenericResponse]("POST", fmt.Sprintf("/v1/games/%d/move", gameId), nil,
		&model.GameMakeMove{Move: move})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
