package tools

import (
	"fmt"
)

type EventController struct {
	consumed   bool
	suppressed bool
}

func NewEventController() *EventController {
	return &EventController{
		consumed:   false,
		suppressed: false,
	}
}

type AnyEventEmitter struct {
	events map[string][]func(args ...interface{})
}

func NewAnyEventEmitter() *AnyEventEmitter {
	return &AnyEventEmitter{
		events: make(map[string][]func(args ...interface{})),
	}
}

func (emitter *AnyEventEmitter) On(event string, listener func(args ...interface{})) {
	if _, exists := emitter.events[event]; !exists {
		emitter.events[event] = []func(args ...interface{}){}
	}
	emitter.events[event] = append(emitter.events[event], listener)
}

func (emitter *AnyEventEmitter) Off(event string, listener func(args ...interface{})) {
	if _, exists := emitter.events[event]; exists {
		listeners := emitter.events[event]
		for i, l := range listeners {
			if fmt.Sprintf("%v", l) == fmt.Sprintf("%v", listener) {
				emitter.events[event] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

func (emitter *AnyEventEmitter) Emit(event string, args ...interface{}) {
	if listeners, exists := emitter.events[event]; exists {
		controller := NewEventController()
		for _, listener := range listeners {
			if !controller.suppressed {
				listener(append(args, controller)...)
			}
		}
	}
}

var Aee = NewAnyEventEmitter()
