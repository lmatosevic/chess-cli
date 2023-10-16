package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lmatosevic/chess-cli/pkg/database/repository"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

// Login godoc
// @Summary Login registered player
// @Description Login registered player
// @Tags auth
// @Accept json
// @Produce json
// @Param player body model.PlayerRequest true "Login player"
// @Success 200 {object} model.AccessToken "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /v1/auth/login [post]
func Login(c *gin.Context) {
	pr, err := utils.ParseJson[model.PlayerRequest](c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	username := strings.TrimSpace(pr.Username)

	p, err := repository.FindPlayerByUsername(username)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(pr.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Success: false, Error: "Invalid player password provided"})
		return
	}

	at, err := repository.CreateAccessToken(p.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.AccessToken{Token: at.Token})
}

// AuthPlayer godoc
// @Summary Get authorized player
// @Description Get authorized player
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.Player "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/auth/player [get]
func AuthPlayer(c *gin.Context) {
	p, err := GetAuthPlayer(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, makePlayerDTO(p))
}

// Logout godoc
// @Summary Logout authorized player
// @Description Logout authorized player
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} model.GenericResponse "Ok"
// @Failure 400 {object} model.ErrorResponse "Bad Request"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /v1/auth/logout [post]
func Logout(c *gin.Context) {
	token := ParseAuthorizationHeader(c)
	err := repository.RevokeAccessToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Success: false, Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, model.GenericResponse{Success: false})
}
