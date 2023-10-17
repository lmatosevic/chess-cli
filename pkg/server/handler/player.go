package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lmatosevic/chess-cli/pkg/database/repository"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

// ListPlayers godoc
// @Summary Query and list players
// @Description Query and list players
// @Tags players
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param sort query string false "Sort"
// @Param filter query string false "Filter"
// @Success 200 {object} model.PlayerListResponse "Ok"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/players [get]
func ListPlayers(c *gin.Context) {
	_, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	page, size, sort, filter := ParseQueryParams(c)

	players, err := repository.QueryPlayers(filter, page, size, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	totalCount, err := repository.CountPlayers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	playersDTO := make([]model.Player, 0)
	for _, p := range *players {
		playersDTO = append(playersDTO, makePlayerDTO(&p))
	}

	c.JSON(http.StatusOK, model.ListResponse[model.Player]{
		Items:       playersDTO,
		ResultCount: len(playersDTO),
		TotalCount:  totalCount,
	})
}

// FindOnePlayer godoc
// @Summary Find one player
// @Description Find one player
// @Tags players
// @Produce json
// @Param id path int true "Player ID"
// @Success 200 {object} model.Player "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/players/{id} [get]
func FindOnePlayer(c *gin.Context) {
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

	p, err := repository.FindPlayerById(int64(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makePlayerDTO(p))
}

// RegisterPlayer godoc
// @Summary Register new player
// @Description Register new player
// @Tags players
// @Accept json
// @Produce json
// @Param player body model.PlayerRequest true "Register player"
// @Success 200 {object} model.Player "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/players/register [post]
func RegisterPlayer(c *gin.Context) {
	pr, err := utils.ParseJson[model.PlayerRequest](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	username := strings.TrimSpace(pr.Username)

	_, err = repository.FindPlayerByUsername(username)
	if err == nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: "Player with username already exists"})
		return
	}

	p, err := repository.CreatePlayer(username, pr.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makePlayerDTO(p))
}

// UpdatePlayer godoc
// @Summary Update player account
// @Description Update player account
// @Tags players
// @Accept json
// @Produce json
// @Param player body model.PlayerRequest true "Update player data"
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/players/update [put]
func UpdatePlayer(c *gin.Context) {
	player, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	pr, err := utils.ParseJson[model.PlayerRequest](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	username := strings.TrimSpace(pr.Username)

	if username != "" {
		player.Username = username
	}

	if pr.Password != "" {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(pr.Password), 6)
		if err == nil {
			player.PasswordHash = string(passwordHash)
		}
	}

	err = repository.UpdatePlayer(player)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	if username != "" {
		err = repository.UpdateGamePlayerUsername(player.Id, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, model.GenericResponse{Success: true})
}

// DeletePlayer godoc
// @Summary Delete player account
// @Description Delete player account
// @Tags players
// @Accept json
// @Produce json
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 401 {object} model.ErrorResponse "Unauthorized"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/players/delete [delete]
func DeletePlayer(c *gin.Context) {
	at, err := GetAccessToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	err = repository.DeletePlayer(at.PlayerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GenericResponse{Success: true})
}

func makePlayerDTO(p *repository.Player) model.Player {
	return model.Player{Id: p.Id, Username: p.Username, Wins: p.Wins, Losses: p.Losses, Draws: p.Draws, Rate: p.Rate,
		Elo: p.Elo, IsPlaying: p.IsPlaying, LastPlayedAt: p.FormatLastPlayedAt(), CreatedAt: p.FormatCreatedAt()}
}
