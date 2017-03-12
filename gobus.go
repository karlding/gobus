package gobus

import (
	"fmt"
	"reflect"
	"sync"
)

// GoBus is a Message Bus
type GoBus struct {
	handlers map[string][]*messageHandler
	lock     sync.Mutex
	wg       sync.WaitGroup
}

type messageHandler struct {
	callBack reflect.Value
}

// New creates a new GoBus with no Message Handlers
func New() *GoBus {
	return &GoBus{
		make(map[string][]*messageHandler),
		sync.Mutex{},
		sync.WaitGroup{},
	}
}

func (bus *GoBus) subscribe(topic string, fn interface{}, handler *messageHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	bus.handlers[topic] = append(bus.handlers[topic], handler)
	return nil
}

func (bus *GoBus) publish(topic string, handler *messageHandler, args ...interface{}) {
	defer bus.wg.Done()

	passedArguments := make([]reflect.Value, 0)
	for _, arg := range args {
		passedArguments = append(passedArguments, reflect.ValueOf(arg))
	}

	handler.callBack.Call(passedArguments)
}

// Publish publishes messages on a topic
func (bus *GoBus) Publish(topic string, args ...interface{}) {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if handlers, ok := bus.handlers[topic]; ok {
		for _, handler := range handlers {
			bus.wg.Add(1)
			go bus.publish(topic, handler, args...)
		}
	}
}

// Subscribe subscribes to a topic
func (bus *GoBus) Subscribe(topic string, fn interface{}) error {
	return bus.subscribe(topic, fn, &messageHandler{
		reflect.ValueOf(fn),
	})
}

func (bus *GoBus) removeHandler(topic string, callback reflect.Value) {
	if _, ok := bus.handlers[topic]; ok {
		for i, handler := range bus.handlers[topic] {
			if handler.callBack == callback {
				bus.handlers[topic] = append(bus.handlers[topic][:i], bus.handlers[topic][i+1:]...)
			}
		}
	}
}

// Unsubscribe removes the associated handler from a topic
// If there are no callbacks on the topic, an error is returned
func (bus *GoBus) Unsubscribe(topic string, handler interface{}) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	if _, ok := bus.handlers[topic]; ok && len(bus.handlers[topic]) > 0 {
		bus.removeHandler(topic, reflect.ValueOf(handler))
		return nil
	}
	return fmt.Errorf("topic %s doesn't exist", topic)
}

// Wait blocks until all async callbacks are run
func (bus *GoBus) Wait() {
	bus.wg.Wait()
}
