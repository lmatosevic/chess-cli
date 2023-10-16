package cli

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lmatosevic/chess-cli/pkg/game"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
)

func ShowServerInfo(status *model.Status) {
	utils.PrintStruct(status)
}

func ShowRegistrationMessage() {
	fmt.Println("registration successful")
}

func ShowLoginMessage() {
	fmt.Println("login successful")
}

func ShowLogoutMessage() {
	fmt.Println("logout successful")
}

func ShowPasswordChangeMessage() {
	fmt.Println("password changed successful")
}

func ShowUserInfo(user *model.Player) {
	utils.PrintStruct(user)
}

func ShowEvent(event *model.Event) {
	utils.PrintStruct(event)
}

func ShowGameList(list *model.GameListResponse) {
	title := fmt.Sprintf("Games | Total: %d | Results: %d", list.TotalCount, list.ResultCount)
	headers := table.Row{"ID", "Name", "Public", "White Player ID", "Black Player ID", "In progress", "Created at"}
	rows := make([]table.Row, 0)
	for _, g := range list.Items {
		rows = append(rows, table.Row{g.Id, g.Name, g.Public, g.WhitePlayerId, g.BlackPlayerId, g.InProgress,
			utils.ToLocalDate(g.CreatedAt)})
	}

	utils.PrintTable(title, headers, rows)
}

func ShowGameInfo(game *model.Game, moves *model.GameMoveListResponse) {
	utils.PrintStruct(game)

	fmt.Println()
	utils.PrintChessBoard(game.Tiles)
	fmt.Println()

	if len(moves.Items) > 0 {
		fmt.Print("Moves history: ")
		for _, move := range moves.Items {
			fmt.Printf("%s ", move.Move)
		}
		fmt.Print("\n\n")
	}

	side := "white"
	if game.WhitePlayerId != 0 {
		side = "black"
	}
	status := fmt.Sprintf("waiting for %s player to join the game", side)
	if game.InProgress {
		side = "white"
		if len(moves.Items) > 0 && moves.Items[len(moves.Items)-1].PlayerId == game.BlackPlayerId {
			side = "black"
		}
		status = fmt.Sprintf("%s player is on turn", side)
	} else if game.EndedAt != "" && game.WinnerId > 0 {
		side = "white"
		if game.WinnerId == game.BlackPlayerId {
			side = "black"
		}
		status = fmt.Sprintf("%s player has won the game", side)
	} else if game.EndedAt != "" {
		status = "game ended in a draw"
	}

	fmt.Printf("Status: %s\n", status)
}

func ShowCreateGameMessage(gameId int64) {
	fmt.Println("game created with ID: ", gameId)
}

func ShowJoinGameMessage() {
	fmt.Println("game joined")
}

func ShowQuitGameMessage() {
	fmt.Println("game quit")
}

func ShowGameMoveResult(move string, game *model.Game, moves *model.GameMoveListResponse) {
	fmt.Printf("played move %s", move)

	fmt.Println()
	utils.PrintChessBoard(game.Tiles)
	fmt.Println()

	if len(moves.Items) > 0 {
		fmt.Print("Moves history: ")
		for _, m := range moves.Items {
			fmt.Printf("%s ", m.Move)
		}
	}
}

func ShowGameMovesHelp() {
	fmt.Print("The move input should be formatted according to the Public Game Notation chess standard: " +
		"(figure)(file*)(rank*)(dest_file)(dest_rank)(figure_to_promote*)\n")
	fmt.Print("* - marks the optional parts of the move string\n")
	fmt.Print("Example valid moves: Paa3, Qa3, Nbf3, Bf1c4, Ph7h8Q\n")
	fmt.Printf("King side castling move is marked as %s and queen side castling as %s string\n",
		game.KingSideCastligMove, game.QueenSideCastligMove)
	fmt.Printf("To make a draw request, use the following sign: %s, and to accept the draw request use also the same "+
		"sign: %s, or to reject it use: %s\n", game.DrawOfferMove, game.DrawOfferMove, game.DrawOfferRejectMove)
}

func ShowPlayerList(list *model.PlayerListResponse) {
	title := fmt.Sprintf("Players | Total: %d | Results: %d", list.TotalCount, list.ResultCount)
	headers := table.Row{"ID", "Username", "Wins", "Losses", "Draws", "Rate", "Elo", "Is playing", "Last Played At"}
	rows := make([]table.Row, 0)
	for _, p := range list.Items {
		rows = append(rows, table.Row{p.Id, p.Username, p.Wins, p.Losses, p.Draws, fmt.Sprintf("%.2f%%", p.Rate*100),
			p.Elo, p.IsPlaying, utils.ToLocalDate(p.LastPlayedAt)})
	}

	utils.PrintTable(title, headers, rows)
}

func ShowPlayerInfo(player *model.Player) {
	utils.PrintStruct(player)
}
