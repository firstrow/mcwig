package mcwig

import "sync"

type Position struct {
	Line int
	Char int
}

type Range struct {
	Start Position
	End   Position
}

type EventTextChange struct {
	Buf     *Buffer
	Start   Position
	End     Position
	NewEnd  Position
	Text    string
	OldText string

	wg *sync.WaitGroup
}

type EventsManager struct {
	source         chan any
	listeners      []chan Event
	newListener    chan chan Event
	removeListener chan (<-chan any)
}

func NewEventsManager() *EventsManager {
	e := &EventsManager{
		source:         make(chan any, 32),
		listeners:      make([]chan Event, 32),
		newListener:    make(chan chan Event, 32),
		removeListener: make(chan (<-chan any)),
	}
	go e.start()
	return e
}

func (e *EventsManager) Subscribe() <-chan Event {
	c := make(chan Event, 0)
	e.newListener <- c
	return c
}

func (e *EventsManager) Unsubscribe(ch <-chan any) {
	// TODO
}

type Event struct {
	Msg any
	Wg  *sync.WaitGroup
}

// this is very quick and dirty implementation.
func (e *EventsManager) Broadcast(msg any) {
	wg := sync.WaitGroup{}
	for _, l := range e.listeners {
		if l == nil {
			continue
		}
		wg.Add(1)
		l <- Event{msg, &wg}
		wg.Wait()
	}
}

func (e *EventsManager) start() {
	for {
		select {
		case l := <-e.newListener:
			e.listeners = append(e.listeners, l)
			// case msg := <-e.source:
			// wg := sync.WaitGroup{}
			// for _, l := range e.listeners {
			// if l == nil {
			// continue
			// }
			// wg.Add(1)
			// l <- Event{msg, &wg}
			// wg.Wait()
			// }
		}
	}
}

