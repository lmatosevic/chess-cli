package command

import (
	"errors"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
)

func ListGames(page int, size int, sort string, filter string) (*model.GameListResponse, error) {
	params := BuildQueryParams(page, size, sort, filter)

	resp, err := client.SendRequest[model.GameListResponse]("GET", "/v1/games", &params, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
