package command

import (
	"errors"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
)

func ChangePassword(password string) (*model.GenericResponse, error) {
	resp, err := client.SendRequest[model.GenericResponse]("PUT", "/v1/players/update/", nil,
		&model.PlayerRequest{Password: password})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	return &resp.Data, nil
}
