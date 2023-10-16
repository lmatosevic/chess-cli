package command

import (
	"errors"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/server/handler"
	"sync"
)

func ListenEvents(events []string, gameId int64, onEvent func(event *model.Event, end func())) (func(), func(), error) {
	var wg sync.WaitGroup

	waitFn := func() {
		wg.Wait()
	}

	cancelFn := func() {
		for i := 0; i < len(events); i++ {
			wg.Done()
		}
	}

	for _, event := range events {
		if !handler.IsValidEventType(event) {
			return nil, nil, errors.New("invalid event type: " + event)
		}
		wg.Add(1)
		go handleEvent(event, gameId, &wg, cancelFn, onEvent)
	}

	return waitFn, cancelFn, nil
}

func handleEvent(eventType string, gameId int64, wg *sync.WaitGroup, cancel func(),
	onEvent func(event *model.Event, end func())) {
	defer wg.Done()
	err := client.SubscribeOnEvent(eventType, gameId, cancel, onEvent)
	if err != nil {
		cancel()
	}
}
