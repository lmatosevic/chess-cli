package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lmatosevic/chess-cli/pkg/database/repository"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
)

const (
	GameAnyEvent             = "GameAnyEvent"
	GameMoveEvent            = "GameMoveEvent"
	GameJoinEvent            = "GameJoinEvent"
	GameQuitEvent            = "GameQuitEvent"
	GameStartEvent           = "GameStartEvent"
	GameEndEvent             = "GameEndEvent"
	GameWhitePlayerMoveEvent = "GameWhitePlayerMoveEvent"
	GameBlackPlayerMoveEvent = "GameWhitePlayerMoveEvent"
	PlayerMessage            = "PlayerMessage"
)

var eventChannels = make(map[string]chan model.Event)

func SendEvent(eventType string, gameId int64, playerId int64, payload string) {
	event := model.Event{Type: eventType, Timestamp: utils.ISODateNow(),
		Data: model.EventData{GameId: gameId, PlayerId: playerId, Payload: payload}}

	for _, eventChan := range eventChannels {
		eventChan <- event
	}
}

// SubscribeToEvent godoc
// @Summary Subscribe to server sent events
// @Description Subscribe to server sent events
// @Tags events
// @Accept json
// @Produce text/event-stream
// @Param token query string true "Access token"
// @Param event query string true "Event type" Enums(GameAnyEvent, GameMoveEvent, GameJoinEvent, GameQuitEvent, GameStartEvent, GameEndEvent, GameWhitePlayerMoveEvent, GameWhitePlayerMoveEvent)
// @Param gameId query int false "Game ID"
// @Success 200 {object} model.Event "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 403 {object} model.ErrorResponse "Forbidden"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/events/subscribe [get]
func SubscribeToEvent(c *gin.Context) {
	accessToken := c.Query("token")
	at, err := repository.FindAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: "Invalid access token"})
		return
	}
	player, err := repository.FindPlayerById(at.PlayerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	eventType := c.Query("event")
	if eventType == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: "Event type is required"})
		return
	}
	if !IsValidEventType(eventType) {
		c.JSON(http.StatusBadRequest,
			model.ErrorResponse{Success: false, Error: fmt.Sprintf("Invalid event type: %s", eventType)})
		return
	}

	gameId := c.Query("gameId")
	var gid int64
	var game *repository.Game

	// Check access of player to the game
	if strings.HasPrefix(eventType, "Game") {
		id, e := strconv.Atoi(c.Query("gameId"))
		if e != nil || id == 0 {
			c.JSON(http.StatusBadRequest,
				model.ErrorResponse{Success: false, Error: fmt.Sprintf("Invalid required gameId: %s", gameId)})
			return
		}
		gid = int64(id)

		game, e = repository.FindGameById(gid)
		if e != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: e.Error()})
			return
		}

		if game.PasswordHash.String != "" && game.WhitePlayerId.Int64 != player.Id &&
			game.BlackPlayerId.Int64 != player.Id {
			c.JSON(http.StatusForbidden,
				model.ErrorResponse{Success: false, Error: "The game is private and player has not joined this game"})
			return
		}
	}

	// Generate uuid and create new subscriber channel
	requestId := uuid.New().String()
	eventChannels[requestId] = make(chan model.Event)

	// Infinite loop function which handles messages from event channel and removes client subscription on disconnect
	c.Stream(func(w io.Writer) bool {
		select {
		case <-c.Request.Context().Done():
			delete(eventChannels, requestId)
			return true

		case event := <-eventChannels[requestId]:
			if shouldReceiveEvent(event, eventType, gid, player, game) {
				c.SSEvent("message", event)
			}
			return true
		}
	})
}

func IsValidEventType(eventType string) bool {
	return slices.Contains([]string{GameAnyEvent, GameMoveEvent, GameJoinEvent, GameQuitEvent, GameStartEvent,
		GameEndEvent, GameWhitePlayerMoveEvent, GameBlackPlayerMoveEvent, PlayerMessage}, eventType)
}

func shouldReceiveEvent(event model.Event, eventType string, gameId int64, player *repository.Player,
	game *repository.Game) bool {

	if event.Type != eventType && (eventType != GameAnyEvent || !strings.HasPrefix(event.Type, "Game")) {
		return false
	}

	if gameId > 0 && gameId != game.Id {
		return false
	}

	if eventType == PlayerMessage && player.Id != event.Data.PlayerId {
		return false
	}

	return true
}
