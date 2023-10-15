package command

import (
	"errors"
	"gitlab.com/lmatosevic/chess-cli/pkg/client"
	"gitlab.com/lmatosevic/chess-cli/pkg/model"
	"gitlab.com/lmatosevic/chess-cli/pkg/utils"
)

func Logout() (*model.GenericResponse, error) {
	resp, err := client.SendRequest[model.GenericResponse]("POST", "/v1/auth/logout", nil, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	tokenFilePath := utils.HomeFilePath(AccessTokenFile)

	if utils.FileExists(tokenFilePath) {
		err = utils.DeleteFile(tokenFilePath)
		if err != nil {
			return nil, err
		}
	}

	return &resp.Data, nil
}
