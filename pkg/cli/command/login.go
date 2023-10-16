package command

import (
	"errors"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
)

func Login(username string, password string, stateless bool) (*model.AccessToken, error) {
	resp, err := client.SendRequest[model.AccessToken]("POST", "/v1/auth/login", nil,
		&model.PlayerRequest{Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Error.Error)
	}

	client.SetAccessToken(resp.Data.Token)

	if !stateless {
		err = utils.WriteToFile(utils.HomeFilePath(AccessTokenFile), []byte(client.GetAccessToken()))
		if err != nil {
			return nil, err
		}

		err = utils.WriteToFile(utils.HomeFilePath(ServerHostFile), []byte(client.GetBaseUrl()))
		if err != nil {
			return nil, err
		}
	}

	return &resp.Data, nil
}
