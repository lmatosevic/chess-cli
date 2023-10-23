package cli

import (
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/cli/command"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

func Run() {
	var server string
	var username string
	var password string
	var token string
	var stateless bool

	app := &cli.App{
		Name:    "Chess CLI",
		Usage:   "Play a game of chess using command line interface",
		Version: "1.0.3",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "server",
				Aliases:     []string{"s"},
				Usage:       "chess server base URL",
				Destination: &server,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "players username or email",
				Destination: &username,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "players password",
				Destination: &password,
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "players access token",
				Destination: &token,
				Required:    false,
			},
			&cli.BoolFlag{
				Name:        "stateless",
				Aliases:     []string{"l"},
				Usage:       "do not use and set default configs from home directory",
				Destination: &stateless,
				Required:    false,
			},
		},
		Action: func(cCtx *cli.Context) error {
			if cCtx.NArg() > 0 {
				return errors.New("unknown arguments provided: " + strings.Join(cCtx.Args().Slice(), " "))
			}

			return StartInteractiveMode(server, username, password, token, stateless)
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "show server info",
				Action: func(cCtx *cli.Context) error {
					if err := StaticServerOnlyInput(server, stateless); err != nil {
						return err
					}

					info, err := command.ServerInfo()
					if err != nil {
						return err
					}

					ShowServerInfo(info)
					return nil
				},
			},
			{
				Name:    "register",
				Aliases: []string{"r"},
				Usage:   "register new player account",
				Action: func(cCtx *cli.Context) error {
					if err := StaticServerOnlyInput(server, stateless); err != nil {
						return err
					}

					_, err := command.Register(username, password, stateless)
					if err != nil {
						return err
					}

					ShowRegistrationMessage()
					return nil
				},
			},
			{
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "login into your account",
				Action: func(cCtx *cli.Context) error {
					if err := StaticServerOnlyInput(server, stateless); err != nil {
						return err
					}

					_, err := command.Login(username, password, stateless)
					if err != nil {
						return err
					}

					ShowLoginMessage()
					return nil
				},
			},
			{
				Name:    "logout",
				Aliases: []string{"o"},
				Usage:   "logout from the server",
				Action: func(cCtx *cli.Context) error {
					if err := StaticInputs(server, username, password, token, stateless); err != nil {
						return err
					}

					_, err := command.Logout()
					if err != nil {
						return err
					}

					ShowLogoutMessage()
					return nil
				},
			},
			{
				Name:    "changePassword",
				Aliases: []string{"c"},
				Usage:   "change your account's password",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "newPassword", Required: true},
				},
				Action: func(cCtx *cli.Context) error {
					if err := StaticInputs(server, username, password, token, stateless); err != nil {
						return err
					}

					_, err := command.ChangePassword(cCtx.String("newPassword"))
					if err != nil {
						return err
					}

					ShowPasswordChangeMessage()
					return nil
				},
			},
			{
				Name:    "whoami",
				Aliases: []string{"w"},
				Usage:   "show your account information",
				Action: func(cCtx *cli.Context) error {
					if err := StaticInputs(server, username, password, token, stateless); err != nil {
						return err
					}

					user, err := command.UserInfo()
					if err != nil {
						return err
					}

					ShowUserInfo(user)
					return nil
				},
			},
			{
				Name:    "events",
				Aliases: []string{"e"},
				Usage:   "subscribe to server sent events and show them in real-time",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{Name: "type", Required: true},
					&cli.Int64Flag{Name: "gameId"},
				},
				Action: func(cCtx *cli.Context) error {
					if err := StaticInputs(server, username, password, token, stateless); err != nil {
						return err
					}

					wait, _, err := command.ListenEvents(cCtx.StringSlice("type"), cCtx.Int64("gameId"),
						func(event *model.Event, end func()) {
							ShowEvent(event)
						})
					if err != nil {
						return err
					}

					wait()

					return nil
				},
			},
			{
				Name:     "game",
				Aliases:  []string{"g", "games"},
				Category: "games",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list all games",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "page"},
							&cli.IntFlag{Name: "size"},
							&cli.StringFlag{Name: "sort"},
							&cli.StringFlag{Name: "filter"},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							list, err := command.ListGames(cCtx.Int("page"), cCtx.Int("size"), cCtx.String("sort"),
								cCtx.String("filter"))
							if err != nil {
								return err
							}

							ShowGameList(list)
							return nil
						},
					},
					{
						Name:  "info",
						Usage: "show information about the game",
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "gameId", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							game, moves, err := command.GameInfo(cCtx.Int64("gameId"))
							if err != nil {
								return err
							}

							ShowGameInfo(game, moves)
							return nil
						},
					},
					{
						Name:  "create",
						Usage: "create new game",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "name", Required: true},
							&cli.StringFlag{Name: "password", Usage: "Make this game password protected"},
							&cli.IntFlag{Name: "turnDuration", Usage: "For unlimited duration use -1"},
							&cli.BoolFlag{Name: "white"},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							game, err := command.CreateGame(cCtx.String("name"), cCtx.String("password"),
								int32(cCtx.Int("turnDuration")), cCtx.Bool("white"))
							if err != nil {
								return err
							}

							ShowCreateGameMessage(game.Id)
							return nil
						},
					},
					{
						Name:  "join",
						Usage: "join existing game",
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "gameId", Required: true},
							&cli.StringFlag{Name: "password", Usage: "Use only when joining private game"},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							_, err := command.JoinGame(cCtx.Int64("gameId"), cCtx.String("password"))
							if err != nil {
								return err
							}

							ShowJoinGameMessage()
							return nil
						},
					},
					{
						Name:  "quit",
						Usage: "quit currently active game",
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "gameId", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							_, err := command.QuitGame(cCtx.Int64("gameId"))
							if err != nil {
								return err
							}

							ShowQuitGameMessage()
							return nil
						},
					},
					{
						Name:  "play",
						Usage: "play move in currently active game",
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "gameId", Required: true},
							&cli.StringFlag{Name: "move", Required: true, Usage: "Use portable game notation (e.g. Pa2a4)"},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							resp, err := command.PlayGameMove(cCtx.Int64("gameId"), cCtx.String("move"))
							if err != nil {
								return err
							}

							game, moves, err := command.GameInfo(cCtx.Int64("gameId"))
							if err != nil {
								return err
							}

							ShowGameMoveResult(resp.Data, game, moves)
							return nil
						},
					},
					{
						Name:  "manual",
						Usage: "Shows the instructions for all types of available moves",
						Action: func(cCtx *cli.Context) error {
							ShowGameMovesHelp()
							return nil
						},
					},
				},
			},
			{
				Name:     "player",
				Aliases:  []string{"p", "players"},
				Category: "players",
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list all players",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "page"},
							&cli.IntFlag{Name: "size"},
							&cli.StringFlag{Name: "sort"},
							&cli.StringFlag{Name: "filter"},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							list, err := command.ListPlayers(cCtx.Int("page"), cCtx.Int("size"), cCtx.String("sort"),
								cCtx.String("filter"))
							if err != nil {
								return err
							}

							ShowPlayerList(list)
							return nil
						},
					},
					{
						Name:  "info",
						Usage: "show information about the player",
						Flags: []cli.Flag{
							&cli.Int64Flag{Name: "playerId", Required: true},
						},
						Action: func(cCtx *cli.Context) error {
							if err := StaticInputs(server, username, password, token, stateless); err != nil {
								return err
							}

							player, err := command.PlayerInfo(cCtx.Int64("playerId"))
							if err != nil {
								return err
							}

							ShowPlayerInfo(player)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\nERROR: %s\n", err)
	}
}
