package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lmatosevic/chess-cli/configs"
	"github.com/lmatosevic/chess-cli/docs"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/server/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Run() {
	conf := *configs.GetConfig()

	version := "1.0.3"

	docs.SwaggerInfo.Title = conf.General.AppName
	docs.SwaggerInfo.Description = conf.General.Description
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = func() string {
		if strings.HasPrefix(conf.Server.Hostname, "http") {
			return strings.Split(conf.Server.Hostname, "://")[1]
		} else {
			return conf.Server.Hostname
		}
	}()

	if !conf.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Enable CORS
	r.Use(cors.Default())

	// Configure status and info route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.Status{Name: conf.General.AppName, Version: version, Status: "ok",
			SwaggerURL: fmt.Sprintf("%s/swagger/index.html", conf.Server.Hostname)})
	})

	// Define API endpoints and handlers
	v1 := r.Group("/v1")
	{
		players := v1.Group("/players")
		{
			players.GET("/", handler.ListPlayers)
			players.GET("/:id", handler.FindOnePlayer)
			players.POST("/register", handler.RegisterPlayer)
			players.PUT("/update", handler.UpdatePlayer)
			players.DELETE("/delete", handler.DeletePlayer)
		}

		games := v1.Group("/games")
		{
			games.GET("/", handler.ListGames)
			games.GET("/:id", handler.FindOneGame)
			games.POST("/create", handler.CreateGame)
			games.POST("/:id/join", handler.JoinGame)
			games.POST("/:id/quit", handler.QuitGame)
			games.GET("/:id/moves", handler.ListGameMoves)
			games.POST("/:id/move", handler.MakeGameMove)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.GET("/player", handler.AuthPlayer)
			auth.POST("/logout", handler.Logout)
		}

		events := v1.Group("/events")
		{
			events.GET("/subscribe", handler.SubscribeToEvent)
		}
	}

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := r.Run(conf.Server.Host + ":" + strconv.Itoa(int(conf.Server.Port)))
	if err != nil {
		log.Fatalf("Error while starting server: %s", err.Error())
	}
}
