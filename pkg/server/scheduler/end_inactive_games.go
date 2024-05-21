package scheduler

import (
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/database/repository"
	"github.com/lmatosevic/chess-cli/pkg/server/handler"
	"log"
)

func EndInactiveGames() {
	inactiveGames, err := repository.FindInactiveGames()
	if err != nil {
		log.Printf("Error while querying inactive games: %s", err.Error())
		return
	}

	for _, game := range *inactiveGames {
		if game.LastMovePlayedAt.Valid {
			moves, e := repository.QueryGameMoves(fmt.Sprintf("gameId=%d", game.Id), 1, 10000, "")
			if e != nil {
				log.Printf("Error while querying game moves: %s", e.Error())
				continue
			}

			winner, e := repository.FindPlayerById((*moves)[len(*moves)-1].PlayerId.Int64)
			if e != nil {
				log.Printf("Error while querying player: %s", e.Error())
				continue
			}

			loserId := game.WhitePlayerId
			if game.WhitePlayerId.Int64 == winner.Id {
				loserId = game.BlackPlayerId
			}

			loser, e := repository.FindPlayerById(loserId.Int64)
			if e != nil {
				log.Printf("Error while querying player: %s", e.Error())
				continue
			}

			err = handler.UpdateEndGameState(&game, winner, loser, false)
		} else {
			err = repository.DeleteGame(game.Id)
		}

		if err != nil {
			log.Printf("Error while ending inactive game with ID: %d", game.Id)
			log.Println(err)
		}
	}

	if len(*inactiveGames) > 0 {
		log.Printf("Ended %d inactive games", len(*inactiveGames))
	}
}
