package database

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gitlab.com/lmatosevic/chess-cli/configs"
	"log"
)

var db *sql.DB

func Init() {
	// Initialize database connection only once
	if db != nil {
		return
	}

	conf := *configs.GetConfig()

	var err error

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.Username, conf.Database.Password, conf.Database.Name,
		conf.Database.Schema)
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	log.Println("Database connection is successful")

	if conf.Database.AutoMigrate {
		Migrate()
	}
}

func GetConnection() *sql.DB {
	if db == nil {
		log.Println("Database connection is not initialized")
		return nil
	}
	return db
}

func Migrate() {
	conf := *configs.GetConfig()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", conf.Database.MigrationsDir),
		conf.Database.Name, driver)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("Nothing to migrate")
			return
		}
		panic(err)
	}

	log.Println("Database migrations finished successfully")
}
