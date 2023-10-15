package main

import (
	"gitlab.com/lmatosevic/chess-cli/pkg/database"
	"gitlab.com/lmatosevic/chess-cli/pkg/server"
	"gitlab.com/lmatosevic/chess-cli/pkg/server/scheduler"
)

//	@contact.name	Luka Matošević
//	@contact.url	https://lukamatosevic.com
//	@contact.email	lukamatosevic5@gmail.com

//	@license.name	MIT 2023
//	@license.url	https://www.mit.edu/~amini/LICENSE.md

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				The access token obtained from /login endpoint, required for accessing protected routes

func main() {
	database.Init()
	scheduler.Start()
	server.Run()
}
