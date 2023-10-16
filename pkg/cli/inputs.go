package cli

import (
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/cli/command"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"os"
	"slices"
	"strings"
)

func InteractiveInputs(server string, username string, password string, token string, stateless bool) error {
	fmt.Println("Welcome to the Chess CLI!")

	if server == "" {
		serverBytes, err := utils.ReadFromFile(utils.HomeFilePath(command.ServerHostFile))
		if err == nil && len(serverBytes) > 0 && !stateless {
			server = strings.TrimSpace(string(serverBytes))
		} else {
			srv, err := utils.ReadStringFromStdin("Enter the chess server hostname: ")
			if err != nil {
				return err
			}
			server = srv
			if !stateless {
				err = utils.WriteToFile(utils.HomeFilePath(command.ServerHostFile), []byte(server))
				if err != nil {
					return err
				}
			}
		}
	}

	client.Init(server)

	_, err := command.ServerInfo()
	if err != nil {
		return err
	}

	if token != "" {
		client.SetAccessToken(token)
	} else if username != "" {
		if password == "" {
			pwd, err := utils.ReadPasswordFromStdin(fmt.Sprintf("Password for user %s: ", username))
			if err != nil {
				return err
			}
			password = pwd
		}
		_, err := command.Login(username, password, stateless)
		if err != nil {
			return err
		}
	} else {
		for {
			if !stateless {
				tokenBytes, err := utils.ReadFromFile(utils.HomeFilePath(command.AccessTokenFile))
				if err == nil && len(tokenBytes) > 0 {
					client.SetAccessToken(strings.TrimSpace(string(tokenBytes)))
					break
				}
			}

			option, err := utils.ReadStringFromStdin("\nSelect option: \n1 -> Login with existing player\n2 -> Register as a new player\n3 -> Exit\n\n")
			if !slices.Contains([]string{"1", "2", "3"}, option) {
				fmt.Println("Invalid option selected")
				continue
			}

			if option == "3" {
				os.Exit(0)
			}

			username, err = utils.ReadStringFromStdin("\nUsername: ")
			if err != nil {
				return err
			}

			password, err = utils.ReadPasswordFromStdin(fmt.Sprintf("Password for user %s: ", username))
			if err != nil {
				return err
			}

			if option == "1" {
				_, err := command.Login(username, password, stateless)
				if err != nil {
					fmt.Printf("\nERROR: %s\n", err.Error())
					continue
				}
			} else {
				_, err = command.Register(username, password, stateless)
				if err != nil {
					fmt.Printf("\nERROR: %s\n", err.Error())
					continue
				}
			}

			break
		}
	}

	return nil
}

func StaticInputs(server string, username string, password string, token string, stateless bool) error {
	err := StaticServerOnlyInput(server, stateless)
	if err != nil {
		return err
	}

	if token != "" {
		client.SetAccessToken(token)
	} else if username != "" {
		if password == "" {
			return errors.New("password is not provided")
		}
		_, err := command.Login(username, password, stateless)
		if err != nil {
			return err
		}
	} else {
		if !stateless {
			tokenBytes, err := utils.ReadFromFile(utils.HomeFilePath(command.AccessTokenFile))
			if err == nil && len(tokenBytes) > 0 {
				client.SetAccessToken(strings.TrimSpace(string(tokenBytes)))
			}
		} else {
			return errors.New("username and password or access token is not provided")
		}
	}

	return nil
}

func StaticServerOnlyInput(server string, stateless bool) error {
	if server == "" {
		serverBytes, err := utils.ReadFromFile(utils.HomeFilePath(command.ServerHostFile))
		if err == nil && len(serverBytes) > 0 && !stateless {
			server = strings.TrimSpace(string(serverBytes))
		} else {
			return errors.New("server hostname is not configured")
		}
	}

	client.Init(server)

	return nil
}
