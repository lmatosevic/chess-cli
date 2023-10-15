package scheduler

import (
	"github.com/go-co-op/gocron"
	"log"
	"time"
)

func Start() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(10).Seconds().Do(EndInactiveGames)
	if err != nil {
		log.Fatalf("Error while starting the scheduled job: %s", err.Error())
	}

	s.StartAsync()
}
