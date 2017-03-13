package v1

import "plugin-to-SONA/log"
import "time"
import "plugin-to-SONA/config"

// AsyncHandler is for asyncronous handling for events
type AsyncHandler struct {
	events map[string]*PodAsyncEvent
}

var asyncHandler AsyncHandler

func init() {
	handler := GetAsyncHandler()
	handler.events = make(map[string]*PodAsyncEvent)
}

// GetAsyncHandler returns an instance of AsyncHandler
func GetAsyncHandler() *AsyncHandler {
	return &asyncHandler
}

// Run initiates a helper for asynchronous handling
func (handler *AsyncHandler) Run(podEvent *PodAsyncEvent) bool {
	log.Info("a new event has arrived for asynch handling")
	log.Info("\t", podEvent.ShortString())
	if handler.isRunning() {
		handler.events[podEvent.pod.Name] = podEvent
	} else {
		handler.events[podEvent.pod.Name] = podEvent
		go handler.run()
	}
	return true
}

func (handler *AsyncHandler) isRunning() bool {
	return len(handler.events) > 0
}

func (handler *AsyncHandler) run() {
	for len(handler.events) > 0 {
		time.Sleep(time.Second * config.EventHandlingInterval)
		for name, event := range handler.events {
			log.Infof("pending %s event: %s", event.eventType, name)
		}
	}
}
