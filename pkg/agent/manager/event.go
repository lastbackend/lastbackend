package manager

import (
	"github.com/lastbackend/lastbackend/pkg/agent/context"
	"github.com/lastbackend/lastbackend/pkg/apis/types"
)

type EventManager struct {
	update chan types.Event
	close  chan bool
}

func NewEventManager() *EventManager {
	ctx := context.Get()
	ctx.Log.Info("Create new event Manager")
	var em = new(EventManager)

	em.update = make(chan types.Event)
	em.close = make(chan bool)

	return em
}

func ReleaseEventManager(em *EventManager) error {
	ctx := context.Get()
	ctx.Log.Info("Release event Manager")
	close(em.update)
	close(em.close)
	return nil
}

func (em *EventManager) watch() error {
	ctx := context.Get()
	ctx.Log.Info("start event watcher")

	for {
		select {
		case _ = <-em.close:
			return ReleaseEventManager(em)
		case event := <-em.update:
			ctx := context.Get()
			ctx.Log.Infof("handle event %s", event.ToJson())
		}
	}

	return nil
}
