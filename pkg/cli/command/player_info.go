package command

import (
	"errors"
	"fmt"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
)

func PlayerInfo(playerId int64) (*model.Player, error) {
	resp, err := client.SendRequest[model.Player]("GET", fmt.Sprintf("/v1/players/%d", playerId), nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
