package v1

import (
	"bytes"
	"fmt"

	"k8s.io/client-go/pkg/api/v1"
)

// EventType is to label cache events.
type EventType string

const (
	// AddEvent means add event for watcher
	AddEvent = "add"
	// DeleteEvent means delete event for watcher
	DeleteEvent = "delete"
	// UpdateEvent means update event for watcher
	UpdateEvent = "update"
)

// PodAsyncEvent is an abstraction of cache events for asynchrous handling.
type PodAsyncEvent struct {
	name      string
	pod       *v1.Pod
	newPod    *v1.Pod
	eventType EventType
}

// String returns the contents of PodAsyncEvent as a string.
func (event *PodAsyncEvent) String() string {
	var outBuffer bytes.Buffer

	outBuffer.WriteString(fmt.Sprintf("event is \"%s\"\n", string(event.eventType)))
	outBuffer.WriteString(fmt.Sprintf("\tthe pod is \"%s\" in \"%s\" namespace\n", event.pod.Name, event.pod.GetNamespace()))
	outBuffer.WriteString(fmt.Sprintf("\tthe pod(\"%s\") is %s on \"%s\" host", event.pod.Status.PodIP, event.pod.Status.Phase, event.pod.Status.HostIP))
	for i, c := range event.pod.Status.ContainerStatuses {
		outBuffer.WriteString(fmt.Sprintf("\n\t\t%d: %s(%s)", i, c.Name, c.ContainerID))
	}
	return outBuffer.String()
}

// ShortString returns the contents of PodAsyncEvent as a string briefly.
func (event *PodAsyncEvent) ShortString() string {
	var outBuffer bytes.Buffer

	outBuffer.WriteString(fmt.Sprintf("\"%s\" ", string(event.eventType)))
	outBuffer.WriteString(fmt.Sprintf("pod(\"%s: %s\" of \"%s\" namespace is \"%s\" in \"%s\")", event.pod.Name, event.pod.Status.PodIP, event.pod.GetNamespace(), event.pod.Status.Phase, event.pod.Status.HostIP))
	return outBuffer.String()
}

// NewPodAsyncEvent makes an instance for asynchronous handling for each event.
func NewPodAsyncEvent(eventType EventType, args ...interface{}) *PodAsyncEvent {
	var event PodAsyncEvent

	event.eventType = eventType
	event.pod = args[0].(*v1.Pod)
	if event.eventType == UpdateEvent {
		event.newPod = args[1].(*v1.Pod)
	}
	event.name = event.pod.Name
	return &event
}
