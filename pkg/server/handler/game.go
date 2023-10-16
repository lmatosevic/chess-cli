package handler

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lmatosevic/chess-cli/configs"
	"github.com/lmatosevic/chess-cli/pkg/database/repository"
	"github.com/lmatosevic/chess-cli/pkg/game"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"math"
	"net/http"
	"strconv"
	"time"
)

// ListGames godoc
// @Summary Query and list games
// @Description Query and list games
// @Tags games
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param sort query string false "Sort"
// @Param filter query string false "Filter"
// @Success 200 {object} model.GameListResponse "Ok"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games [get]
func ListGames(c *gin.Context) {
	_, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	page, size, sort, filter := ParseQueryParams(c)

	games, err := repository.QueryGames(filter, page, size, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	totalCount, err := repository.CountGames(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	gamesDTO := make([]model.Game, 0)
	for _, g := range *games {
		gamesDTO = append(gamesDTO, makeGameDTO(&g))
	}

	c.JSON(http.StatusOK, model.ListResponse[model.Game]{
		Items:       gamesDTO,
		ResultCount: len(gamesDTO),
		TotalCount:  totalCount,
	})
}

// FindOneGame godoc
// @Summary Find one game
// @Description Find one game
// @Tags games
// @Produce json
// @Param id path int true "Game ID"
// @Success 200 {object} model.Game "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/{id} [get]
func FindOneGame(c *gin.Context) {
	_, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	g, err := repository.FindGameById(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makeGameDTO(g))
}

// CreateGame godoc
// @Summary Create new game
// @Description Create new game
// @Tags games
// @Accept json
// @Produce json
// @Param game body model.GameCreate true "Create game"
// @Success 200 {object} model.Game "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/create [post]
func CreateGame(c *gin.Context) {
	conf := *configs.GetConfig()
	player, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	gc, err := utils.ParseJson[model.GameCreate](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	maxCreatedGames := int(conf.Rules.MaxCreatedGames)
	createdGames, err := repository.QueryGames(fmt.Sprintf("creatorId=%d;and;endedAt=null", player.Id), 1,
		maxCreatedGames, "")
	if err != nil || len(*createdGames) == maxCreatedGames {
		errMsg := fmt.Sprintf("Maximum number of created active games reached (%d)", maxCreatedGames)
		if err != nil {
			errMsg = err.Error()
		}
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: errMsg})
		return
	}

	turnDuration := gc.TurnDurationSeconds
	if turnDuration == 0 {
		turnDuration = conf.Rules.DefaultTurnDurationSeconds
	}

	g, err := repository.CreateGame(gc.Name, gc.Password, turnDuration, player.Id, gc.IsWhite,
		game.MakeStartingBoard())
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makeGameDTO(g))
}

// JoinGame godoc
// @Summary Join existing game
// @Description Join existing game
// @Tags games
// @Accept json
// @Produce json
// @Param id path int true "Game ID"
// @Param game body model.GameJoin false "Join game"
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/{id}/join [post]
func JoinGame(c *gin.Context) {
	conf := *configs.GetConfig()
	player, g, err, code := getPlayerAndGame(c)
	if err != nil {
		c.JSON(code, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	jg, err := utils.ParseJson[model.GameJoin](c.Request.Body)
	if err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	maxJoinedGames := int(conf.Rules.MaxJoinedGames)
	joinedGames, err := repository.QueryGames(
		fmt.Sprintf("whitePlayerId=%d;and;endedAt=null;or;blackPlayerId=%d;and;endedAt=null", player.Id, player.Id), 1,
		maxJoinedGames, "")
	if err != nil || len(*joinedGames) == maxJoinedGames {
		errMsg := fmt.Sprintf("Maximum number of joined active games reached (%d)", maxJoinedGames)
		if err != nil {
			errMsg = err.Error()
		}
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: errMsg})
		return
	}

	if g.PasswordHash.Valid {
		if jg.Password == "" {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Game password is required"})
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(g.PasswordHash.String), []byte(jg.Password))
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: "Invalid game password provided"})
			return
		}
	}

	if g.InProgress {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Game is already in progress"})
		return
	}

	if g.EndedAt.Valid {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Game has already finished"})
		return
	}

	if g.WhitePlayerId.Int64 == player.Id || g.BlackPlayerId.Int64 == player.Id {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "You have already joined this game"})
		return
	}

	g.StartedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	g.InProgress = true

	var side string
	playerId := sql.NullInt64{Int64: player.Id, Valid: true}
	if g.WhitePlayerId.Int64 != 0 {
		g.BlackPlayerId = playerId
		side = "black"
	} else {
		g.WhitePlayerId = playerId
		side = "white"
	}

	err = repository.UpdateGame(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	SendEvent(GameJoinEvent, g.Id, player.Id, side)

	c.JSON(http.StatusOK, model.GenericResponse{Success: true})
}

// QuitGame godoc
// @Summary Quit joined game
// @Description Quit joined game
// @Tags games
// @Accept json
// @Produce json
// @Param id path int true "Game ID"
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ErrorResponse "Forbidden"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/{id}/quit [post]
func QuitGame(c *gin.Context) {
	player, g, err, code := getPlayerAndGame(c)
	if err != nil {
		c.JSON(code, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	if g.WhitePlayerId.Int64 != player.Id && g.BlackPlayerId.Int64 != player.Id {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Forbidden access to not joined game"})
		return
	}

	// Delete game if not joined by both players
	if !g.WhitePlayerId.Valid || !g.BlackPlayerId.Valid {
		err = repository.DeleteGame(g.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, model.GenericResponse{Success: true})
	}

	side := "black"
	winnerId := g.WhitePlayerId
	if g.WhitePlayerId.Int64 == player.Id {
		winnerId = g.BlackPlayerId
		side = "white"
	}

	winnerPlayer, err := repository.FindPlayerById(winnerId.Int64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	err = UpdateEndGameState(g, winnerPlayer, player, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	SendEvent(GameQuitEvent, g.Id, player.Id, side)

	c.JSON(http.StatusOK, model.GenericResponse{Success: true})
}

// MakeGameMove godoc
// @Summary Make game move
// @Description Make game move
// @Tags games
// @Accept json
// @Produce json
// @Param id path int true "Game ID"
// @Param game body model.GameMakeMove true "Game move"
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ErrorResponse "Forbidden"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/{id}/move [post]
func MakeGameMove(c *gin.Context) {
	conf := *configs.GetConfig()
	player, g, err, code := getPlayerAndGame(c)
	if err != nil {
		c.JSON(code, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	if g.WhitePlayerId.Int64 != player.Id && g.BlackPlayerId.Int64 != player.Id {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Forbidden access to not joined game"})
		return
	}

	if !g.InProgress {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Cannot make a move to not started game"})
		return
	}

	gameMoves, err := repository.QueryGameMoves(fmt.Sprintf(`gameId=%d`, g.Id), 1, 10000, "createdAt")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	movesCount := len(*gameMoves)
	var lastMove *repository.GameMove
	if movesCount > 0 {
		lastMove = &(*gameMoves)[movesCount-1]
		if lastMove.PlayerId.Int64 == player.Id {
			c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "Its the other players turn"})
			return
		}
	} else if player.Id != g.WhitePlayerId.Int64 {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Success: false, Error: "The white player is first on turn"})
		return
	}

	var moves []string
	for _, m := range *gameMoves {
		moves = append(moves, m.Move)
	}

	gameModel, err := game.MakeGame(g.Tiles, moves)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	gm, err := utils.ParseJson[model.GameMakeMove](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	isDraw := false
	if lastMove != nil && lastMove.Move == game.DrawOfferMove {
		if gm.Move == game.DrawOfferMove {
			isDraw = true
		} else if gm.Move != game.DrawOfferRejectMove {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false,
				Error: fmt.Sprintf("You must respond to opponents draw request by either accepting (%s) or declining "+
					"(%s) request", game.DrawOfferMove, game.DrawOfferRejectMove)})
			return
		}
	} else if gm.Move == game.DrawOfferMove {
		timeoutTurns := int(conf.Rules.DrawRequestTimeoutTurns)
		if movesCount < timeoutTurns {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false,
				Error: fmt.Sprintf("It must pass at least %d turns before draw can be requested", timeoutTurns)})
			return
		}
		turnsLeft := timeoutTurns
		for i := movesCount - 1; i > movesCount-timeoutTurns; i-- {
			if (*gameMoves)[i].Move == game.DrawOfferMove {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false,
					Error: fmt.Sprintf("It must pass %d more turn/s before draw can be requested again", turnsLeft)})
				return
			}
			turnsLeft--
		}
	} else if gm.Move == game.DrawOfferRejectMove {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false,
			Error: fmt.Sprintf("There is no draw offer from opponent to reject")})
		return
	}

	move, isWin, err := gameModel.MakeMove(gm.Move, g.WhitePlayerId.Int64 == player.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	err = repository.CreateGameMove(g.Id, player.Id, move)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	g.Tiles = gameModel.GetTiles()
	g.LastMovePlayedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	player.LastPlayedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	if isWin || isDraw {
		otherPlayerId := g.WhitePlayerId
		if g.WhitePlayerId.Int64 == player.Id {
			otherPlayerId = g.BlackPlayerId
		}

		otherPlayer, e := repository.FindPlayerById(otherPlayerId.Int64)
		if e != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: e.Error()})
			return
		}

		e = UpdateEndGameState(g, player, otherPlayer, isDraw)
		if e != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: e.Error()})
			return
		}
	} else {
		err = repository.UpdatePlayer(player)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
			return
		}
		err = repository.UpdateGame(g)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
			return
		}
	}

	SendEvent(GameMoveEvent, g.Id, player.Id, move)

	if g.WhitePlayerId.Int64 == player.Id {
		SendEvent(GameWhitePlayerMoveEvent, g.Id, player.Id, move)
	} else {
		SendEvent(GameBlackPlayerMoveEvent, g.Id, player.Id, move)
	}

	c.JSON(http.StatusOK, model.GenericResponse{Success: true, Data: move})
}

// ListGameMoves godoc
// @Summary Query and list game moves
// @Description Query and list game moves
// @Tags games
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param sort query string false "Sort"
// @Param filter query string false "Filter"
// @Success 200 {object} model.GameListResponse "Ok"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/games/{id}/moves [get]
func ListGameMoves(c *gin.Context) {
	_, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	page, size, sort, filter := ParseQueryParams(c)

	gameMoves, err := repository.QueryGameMoves(filter, page, size, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	totalCount, err := repository.CountGameMoves(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	gameMovesDTO := make([]model.GameMove, 0)
	for _, gm := range *gameMoves {
		gameMovesDTO = append(gameMovesDTO, makeGameMoveDTO(&gm))
	}

	c.JSON(http.StatusOK, model.ListResponse[model.GameMove]{
		Items:       gameMovesDTO,
		ResultCount: len(gameMovesDTO),
		TotalCount:  totalCount,
	})
}

func UpdateEndGameState(game *repository.Game, winner *repository.Player, loser *repository.Player, isDraw bool) error {
	winnerElo := winner.Elo
	loserElo := loser.Elo

	if !isDraw {
		// Update winner player data
		winner.Wins = winner.Wins + 1
		winner.Rate = float32(winner.Wins) / float32(winner.Losses+winner.Wins)
		winner.Elo = calculateElo(winnerElo, loserElo, true, false)
	} else {
		winner.Draws = winner.Draws + 1
		winner.Elo = calculateElo(winnerElo, loserElo, false, true)
	}

	err := repository.UpdatePlayer(winner)
	if err != nil {
		return err
	}

	if !isDraw {
		// Update loser player data
		loser.Losses = loser.Losses + 1
		loser.Rate = float32(loser.Wins) / float32(loser.Losses+loser.Wins)
		loser.Elo = calculateElo(loserElo, winnerElo, false, false)
	} else {
		loser.Draws = loser.Draws + 1
		loser.Elo = calculateElo(loserElo, winnerElo, false, true)
	}

	err = repository.UpdatePlayer(loser)
	if err != nil {
		return err
	}

	// Update game data
	if !isDraw {
		game.WinnerId = sql.NullInt64{Int64: winner.Id, Valid: true}
	}
	game.EndedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	game.InProgress = false

	err = repository.UpdateGame(game)
	if err != nil {
		return err
	}

	status := "win"
	if isDraw {
		status = "draw"
	}

	SendEvent(GameEndEvent, game.Id, winner.Id, status)

	return nil
}

// The elo after the match (R'a) is calculated using following formula:
// R'a = Ra + k*(Sa — Ea)
// Ea = Qa /(Qa + Qb)
// Qa = 10^(Ra/c)
// Qb = 10^(Rb/c)
// k = 32, c = 400
func calculateElo(playerElo int32, opponentElo int32, isWinner bool, isDraw bool) int32 {
	k := 32
	c := 400

	qa := math.Pow(10, float64(playerElo)/float64(c))
	qb := math.Pow(10, float64(opponentElo)/float64(c))

	ea := qa / (qa + qb)

	var sa float64
	if isWinner {
		sa = 1
	} else if !isDraw {
		sa = 0
	} else {
		sa = 0.5
	}

	elo := float64(playerElo) + float64(k)*(sa-ea)

	return int32(math.Round(elo))
}

func getPlayerAndGame(c *gin.Context) (*repository.Player, *repository.Game, error, int) {
	player, err := GetAuthPlayer(c)
	if err != nil {
		return nil, nil, err, http.StatusUnauthorized
	}

	idParam, _ := c.Params.Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, nil, err, http.StatusBadRequest
	}

	g, err := repository.FindGameById(int64(id))
	if err != nil {
		return nil, nil, err, http.StatusInternalServerError
	}

	return player, g, nil, http.StatusOK
}

func makeGameDTO(g *repository.Game) model.Game {
	return model.Game{Id: g.Id, Name: g.Name, TurnDurationSeconds: g.TurnDurationSeconds.Int32,
		Public: !g.PasswordHash.Valid, WhitePlayerId: g.WhitePlayerId.Int64, BlackPlayerId: g.BlackPlayerId.Int64,
		WinnerId: g.WinnerId.Int64, CreatorId: g.CreatorId.Int64, InProgress: g.InProgress, Tiles: g.Tiles,
		LastMovePlayedAt: g.FormatLastMovePlayedAt(), StartedAt: g.FormatStartedAt(), EndedAt: g.FormatEndedAt(),
		CreatedAt: g.FormatCreatedAt()}
}

func makeGameMoveDTO(gm *repository.GameMove) model.GameMove {
	return model.GameMove{Id: gm.Id, GameId: gm.GameId, PlayerId: gm.PlayerId.Int64, Move: gm.Move,
		CreatedAt: gm.FormatCreatedAt()}
}
