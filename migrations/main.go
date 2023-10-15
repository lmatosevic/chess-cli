package main

import (
	"gitlab.com/lmatosevic/chess-cli/configs"
	"gitlab.com/lmatosevic/chess-cli/pkg/database"
)

func main() {
	conf := configs.GetConfig()

	database.Init()

	if !conf.Database.AutoMigrate {
		database.Migrate()
	}
}
