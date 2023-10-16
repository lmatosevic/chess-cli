package cli

import (
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/cli/command"
	"github.com/lmatosevic/chess-cli/pkg/game"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/server/handler"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

type NavigationResult struct {
	Page       int
	Sort       string
	Filter     string
	ResourceId int
	IsEnd      bool
}

func StartInteractiveMode(server string, username string, password string, token string, stateless bool) error {
	if err := InteractiveInputs(server, username, password, token, stateless); err != nil {
		return err
	}

out:
	for {
		option, err := utils.ReadStringFromStdin("\nSelect option: \n1 -> Show players\n2 -> Show games\n" +
			"3 -> Resume game\n4 -> Join game\n5 -> Create game\n6 -> Logout\n7 -> Exit\n\n")
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch option {
		case "1":
			ShowPlayers()
		case "2":
			ShowGames()
		case "3":
			ResumeGame()
		case "4":
			JoinGame()
		case "5":
			CreateGame()
		case "6":
			_, err = command.Logout()
			if err != nil {
				fmt.Println(err)
			} else {
				ShowLogoutMessage()
			}
			break out
		case "7":
			break out
		default:
			fmt.Println("Invalid option")
			continue
		}
	}

	return nil
}

func ShowPlayers() {
	page := 1
	size := 20
	sort := "-lastPlayedAt"
	filter := ""

	for {
		players, err := command.ListPlayers(page, size, sort, filter)
		if err != nil {
			fmt.Println(err)
			break
		}

		ShowPlayerList(players)

		resp := pageNavigation(players.TotalCount, players.ResultCount, page, size, sort, filter, true, "Show player details")
		if resp.IsEnd {
			break
		}
		page = resp.Page
		sort = resp.Sort
		filter = resp.Filter

		if resp.ResourceId > 0 {
			player, err := command.PlayerInfo(int64(resp.ResourceId))
			if err != nil {
				fmt.Println(err)
				break
			}
			ShowPlayerInfo(player)

			err = utils.WaitForAnyKey()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func ShowGames() {
	page := 1
	size := 20
	sort := "-createdAt"
	filter := ""

	for {
		games, err := command.ListGames(page, size, sort, filter)
		if err != nil {
			fmt.Println(err)
			break
		}

		ShowGameList(games)

		resp := pageNavigation(games.TotalCount, games.ResultCount, page, size, sort, filter, true, "Show game details")
		if resp.IsEnd {
			break
		}
		page = resp.Page
		sort = resp.Sort
		filter = resp.Filter

		if resp.ResourceId > 0 {
			g, moves, err := command.GameInfo(int64(resp.ResourceId))
			if err != nil {
				fmt.Println(err)
				break
			}
			ShowGameInfo(g, moves)

			err = utils.WaitForAnyKey()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
	}
}

func ResumeGame() {
	user, err := command.UserInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	page := 1
	size := 20
	sort := "-createdAt"
	filter := fmt.Sprintf("whitePlayerId=%d;and;inProgress=true;or;blackPlayerId=%d;and;inProgress=true;or;creatorId=%d;and;startedAt=null",
		user.Id, user.Id, user.Id)

	for {
		games, err := command.ListGames(page, size, sort, filter)
		if err != nil {
			fmt.Println(err)
			break
		}

		if games.TotalCount == 0 {
			fmt.Println("You have no active games in progress. Please join other games or create a new one.")
			return
		}

		ShowGameList(games)

		resp := pageNavigation(games.TotalCount, games.ResultCount, page, size, sort, filter, false, "Select game to resume")
		if resp.IsEnd {
			break
		}
		page = resp.Page
		sort = resp.Sort

		if resp.ResourceId > 0 {
			playGame(int64(resp.ResourceId), user)
		}
	}
}

func JoinGame() {
	user, err := command.UserInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	page := 1
	size := 20
	sort := "-createdAt"
	filter := fmt.Sprintf("whitePlayerId!=%d;and;blackPlayerId!=%d;and;inProgress=false;and;endedAt=null", user.Id, user.Id)

	for {
		games, err := command.ListGames(page, size, sort, filter)
		if err != nil {
			fmt.Println(err)
			break
		}

		if games.TotalCount == 0 {
			fmt.Println("There are no active games for you to join. Please create a new one.")
			return
		}

		ShowGameList(games)

		resp := pageNavigation(games.TotalCount, games.ResultCount, page, size, sort, filter, false, "Select game to Join")
		if resp.IsEnd {
			break
		}
		page = resp.Page
		sort = resp.Sort

		if resp.ResourceId > 0 {
			g, _, err := command.GameInfo(int64(resp.ResourceId))
			if err != nil {
				fmt.Println(err)
				break
			}

			password := ""
			if !g.Public {
				pwd, err := utils.ReadPasswordFromStdin("Enter password: ")
				if err != nil {
					fmt.Println(err)
					break
				}
				password = pwd
			}

			_, err = command.JoinGame(g.Id, password)
			if err != nil {
				fmt.Println(err)
				break
			}

			playGame(g.Id, user)
		}
	}
}

func CreateGame() {
	user, err := command.UserInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

out:
	for {
		name, err := utils.ReadStringFromStdin("Enter game name: ")
		if err != nil {
			fmt.Println(err)
			break
		}
		if strings.TrimSpace(name) == "" {
			fmt.Println("Game name is required")
			continue
		}

		password, err := utils.ReadStringFromStdin("Enter password (optional): ")
		if err != nil {
			fmt.Println(err)
			break
		}

		turnDuration := 0
		for {
			duration, err := utils.ReadStringFromStdin("Enter turn duration in seconds (optional, -1 for unlimited): ")
			if err != nil {
				fmt.Println(err)
				break out
			}

			if duration != "" {
				durationNumber, err := strconv.Atoi(duration)
				if err != nil {
					fmt.Println(err)
					continue
				}
				turnDuration = durationNumber
			}
			break
		}

		white := "1"
		for {
			white, err = utils.ReadStringFromStdin("Choose side:\n1 -> White\n2 -> Black:\n\n")
			if err != nil {
				fmt.Println(err)
				break out
			}
			if !slices.Contains([]string{"1", "2"}, white) {
				fmt.Println("Invalid option")
				continue
			}
			break
		}

		g, err := command.CreateGame(name, strings.TrimSpace(password), int32(turnDuration), strings.ToLower(white) == "1")
		if err != nil {
			fmt.Println(err)
			break
		}

		playGame(g.Id, user)
		break
	}
}

func playGame(gameId int64, player *model.Player) {
	joinChan := make(chan bool)
	turnChan := make(chan bool)
	sigtermChan := make(chan os.Signal)

	signal.Notify(sigtermChan, os.Interrupt, syscall.SIGTERM)

	var opponent *model.Player

	var cancelListener func()

	go func(cancel *func()) {
		wait, c, err := command.ListenEvents([]string{handler.GameAnyEvent}, gameId,
			func(event *model.Event, end func()) {
				if event.Type == handler.GameEndEvent {
					end()
					turnChan <- true
				}
				if event.Type == handler.GameMoveEvent && event.Data.PlayerId != player.Id {
					moveDesc := moveDescription(event.Data.Payload)
					fmt.Printf("\nOpponent played move: %s %s\n", event.Data.Payload, moveDesc)
					turnChan <- true
				}
				if event.Type == handler.GameJoinEvent && event.Data.PlayerId != player.Id {
					username := fmt.Sprintf("ID: %d", event.Data.PlayerId)
					p, err := command.PlayerInfo(event.Data.PlayerId)
					if err != nil {
						fmt.Println(err)
					} else {
						username = p.Username
						opponent = p
					}
					fmt.Printf("\nOpponent %s has joined the game as %s\n", username, event.Data.Payload)
					joinChan <- true
				}
				if event.Type == handler.GameQuitEvent && event.Data.PlayerId != player.Id {
					fmt.Println("Opponent has quit the game")
					turnChan <- true
				}
			})
		if err != nil {
			fmt.Println(err)
			return
		}
		cancel = &c
		wait()
	}(&cancelListener)

out:
	for {
		g, moves, err := command.GameInfo(gameId)
		if err != nil {
			fmt.Println(err)
			break
		}

		side := "white"
		if g.WhitePlayerId != player.Id {
			side = "black"
		}

		if g.StartedAt != "" && opponent == nil {
			oppId := g.BlackPlayerId
			if g.WhitePlayerId != player.Id {
				oppId = g.WhitePlayerId
			}
			p, err := command.PlayerInfo(oppId)
			if err != nil {
				fmt.Println(err)
				break
			}
			opponent = p
		}

		if g.EndedAt != "" {
			fmt.Println(gameEndStatus(side, g))
			err = utils.WaitForAnyKey()
			if err != nil {
				fmt.Println(err)
			}
			break
		}

		if (len(moves.Items) > 0 && moves.Items[len(moves.Items)-1].PlayerId != player.Id) ||
			(len(moves.Items) == 0 && side == "white" && g.InProgress) {
			go func() {
				turnChan <- true
			}()
		} else {
			fmt.Println()
			utils.PrintChessBoard(g.Tiles)
			fmt.Println()

			if g.InProgress {
				fmt.Printf("Waiting for %s's move...\n", opponent.Username)
			}
		}

		if g.StartedAt == "" {
			fmt.Println("Waiting for opponent to join...")
		}

		select {
		case <-joinChan:
			{
				continue
			}
		case <-turnChan:
			{
				g, _, err = command.GameInfo(gameId)
				if err != nil {
					fmt.Println(err)
					break out
				}

				fmt.Println()
				utils.PrintChessBoard(g.Tiles)
				fmt.Println()

				if g.EndedAt != "" {
					fmt.Println(gameEndStatus(side, g))
					err = utils.WaitForAnyKey()
					if err != nil {
						fmt.Println(err)
					}
					break out
				}

				for {
					move, err := utils.ReadStringFromStdin(fmt.Sprintf("Enter move (%s): ", side))
					if err != nil {
						if err.Error() == "EOF" {
							go func() {
								sigtermChan <- syscall.SIGTERM
							}()
							break
						}
						fmt.Println(err)
						break out
					}

					_, err = command.PlayGame(gameId, move)
					if err != nil {
						fmt.Println(err)
					} else {
						break
					}
				}
			}
		case <-sigtermChan:
			{
			ctrl:
				for {
					option, err := utils.ReadStringFromStdin("\nSelect option:\n1 -> Go back\n2 -> Continue playing\n" +
						"3 -> Show help\n4 -> Surrender\n\n")
					if err != nil {
						fmt.Println(err)
						break out
					}

					switch option {
					case "1":
						break out
					case "2":
						break ctrl
					case "3":
						ShowGameMovesHelp()
						err = utils.WaitForAnyKey()
						if err != nil {
							fmt.Println(err)
							break out
						}
						break ctrl
					case "4":
						_, err = command.QuitGame(gameId)
						if err != nil {
							fmt.Println(err)
							break out
						}
						fmt.Println("You have surrendered")
						break ctrl
					default:
						fmt.Println("Invalid option")
					}
				}
				continue
			}
		}
	}

	signal.Reset(os.Interrupt, syscall.SIGTERM)

	if cancelListener != nil {
		cancelListener()
	}
}

func moveDescription(move string) string {
	moveDesc := ""
	if move == game.KingSideCastligMove {
		moveDesc = "(king side castling)"
	}
	if move == game.QueenSideCastligMove {
		moveDesc = "(queen side castling)"
	}
	if move == game.DrawOfferMove {
		moveDesc = "(draw offer)"
	}
	if move == game.DrawOfferRejectMove {
		moveDesc = "(draw offer rejected)"
	}
	return moveDesc
}

func gameEndStatus(side string, game *model.Game) string {
	if game.EndedAt != "" && game.WinnerId > 0 {
		winnerSide := "white"
		if game.WinnerId == game.BlackPlayerId {
			winnerSide = "black"
		}
		if winnerSide == side {
			return "Congratulations! You have won the game"
		} else {
			return "You have lost the game. Better luck next time."
		}
	} else {
		return "Game ended in a draw"
	}
}

func pageNavigation(totalCount int, resultCount int, page int, size int, sort string, filter string, filterEnabled bool,
	firstOption string) *NavigationResult {
	if totalCount == 0 {
		return &NavigationResult{page, sort, filter, 0, true}
	}

	nextPage := ""
	previousPage := ""
	if page*size < totalCount && totalCount != resultCount {
		nextPage = "2 -> Next page\n"
	}
	if page > 1 && totalCount != resultCount {
		previousPage = "3 -> Previous page\n"
	}

	if firstOption != "" {
		firstOption = fmt.Sprintf("1 -> %s\n", firstOption)
	}

	filterOption := ""
	if filterEnabled {
		filterOption = "5 -> Filter table\n"
	}

	option, err := utils.ReadStringFromStdin(
		fmt.Sprintf("\nSelect option: \n%s%s%s4 -> Sort table\n%s6 -> Back\n\n", firstOption, nextPage,
			previousPage, filterOption))
	if err != nil {
		fmt.Println(err)
		return &NavigationResult{page, sort, filter, 0, true}
	}

	switch option {
	case "1":
		idStr, err := utils.ReadStringFromStdin("Enter ID: ")
		if err != nil {
			fmt.Println(err)
			return &NavigationResult{page, sort, filter, 0, true}
		}
		id, err := strconv.Atoi(idStr)
		return &NavigationResult{page, sort, filter, id, false}
	case "2":
		if nextPage == "" {
			fmt.Println("Invalid option")
			return &NavigationResult{page, sort, filter, 0, false}
		}
		page += 1
	case "3":
		if previousPage == "" {
			fmt.Println("Invalid option")
			return &NavigationResult{page, sort, filter, 0, false}
		}
		page -= 1
	case "4":
		sortColumn, err := utils.ReadStringFromStdin("Enter column to sort by: ")
		if err != nil {
			fmt.Println(err)
			return &NavigationResult{page, sort, filter, 0, true}
		}
		return &NavigationResult{page, sortColumn, filter, 0, false}
	case "5":
		if !filterEnabled {
			fmt.Println("Invalid option")
			return &NavigationResult{page, sort, filter, 0, false}
		}
		filterQuery, err := utils.ReadStringFromStdin("Enter query filter: ")
		if err != nil {
			fmt.Println(err)
			return &NavigationResult{page, sort, filter, 0, true}
		}
		return &NavigationResult{page, sort, filterQuery, 0, false}
	case "6":
		return &NavigationResult{page, sort, filter, 0, true}
	}

	return &NavigationResult{page, sort, filter, 0, false}
}
