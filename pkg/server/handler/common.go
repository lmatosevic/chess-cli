package handler

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/lmatosevic/chess-cli/pkg/database/repository"
	"strconv"
	"strings"
)

func ParseQueryParams(c *gin.Context) (int, int, string, string) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.Query("size"))
	if err != nil || size < 1 {
		size = 20
	}

	sort := c.Query("sort")
	if sort == "" {
		sort = "-createdAt"
	}
	filter := c.Query("filter")

	return page, size, sort, filter
}

func ParseAuthorizationHeader(c *gin.Context) string {
	authParts := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authParts) > 1 {
		return authParts[1]
	} else {
		return authParts[0]
	}
}

func GetAccessToken(c *gin.Context) (*repository.AccessToken, error) {
	token := ParseAuthorizationHeader(c)

	at, err := repository.FindAccessToken(token)
	if err != nil {
		return nil, err
	}

	return at, nil
}

func GetAuthPlayer(c *gin.Context) (*repository.Player, error) {
	at, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}

	p, err := repository.FindPlayerById(at.PlayerId)
	if err != nil {
		return nil, err
	}

	return p, nil
}
