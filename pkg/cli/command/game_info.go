package command

import (
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
)

func GameInfo(gameId int64) (*model.Game, *model.GameMoveListResponse, error) {
	resp, err := client.SendRequest[model.Game]("GET", fmt.Sprintf("/v1/games/%d", gameId), nil, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, errors.New(resp.Error.Error)
	}

	params := BuildQueryParams(1, 10000, "createdAt", fmt.Sprintf("gameId=%d", gameId))
	respMoves, err := client.SendRequest[model.GameMoveListResponse]("GET", fmt.Sprintf("/v1/games/%d/moves", gameId),
		&params, nil)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return nil, nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, &respMoves.Data, nil
}
