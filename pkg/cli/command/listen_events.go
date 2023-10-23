package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/lmatosevic/chess-cli/pkg/client"
	"github.com/lmatosevic/chess-cli/pkg/model"
	"github.com/lmatosevic/chess-cli/pkg/server/handler"
	"sync"
)

func ListenEvents(events []string, gameId int64, onEvent func(event *model.Event, end func())) (func(), func(), error) {
	var wg sync.WaitGroup
	ctxCancelFns := make(map[string]func())

	waitFn := func() {
		wg.Wait()
	}

	cancelFn := func() {
		for i := 0; i < len(events); i++ {
			ctxCancelFns[events[i]]()
			wg.Done()
		}
	}

	for _, event := range events {
		if !handler.IsValidEventType(event) {
			return nil, nil, errors.New("invalid event type: " + event)
		}
		wg.Add(1)
		ctx, c := context.WithCancel(context.Background())
		ctxCancelFns[event] = c
		go handleEvent(event, gameId, &wg, ctx, cancelFn, onEvent)
	}

	return waitFn, cancelFn, nil
}

func handleEvent(eventType string, gameId int64, wg *sync.WaitGroup, ctx context.Context, cancel func(),
	onEvent func(event *model.Event, end func())) {
	err := client.SubscribeOnEvent(eventType, gameId, ctx, cancel, onEvent)
	if err != nil {
		fmt.Println(err.Error())
		defer wg.Done()
	}
}
