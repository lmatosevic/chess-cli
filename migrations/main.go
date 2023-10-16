package main

import (
	"github.com/lmatosevic/chess-cli/configs"
	"github.com/lmatosevic/chess-cli/pkg/database"
)

func main() {
	conf := configs.GetConfig()

	database.Init()

	if !conf.Database.AutoMigrate {
		database.Migrate()
	}
}
