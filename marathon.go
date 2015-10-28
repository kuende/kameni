package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kuende/kameni/Godeps/_workspace/src/github.com/mailgun/log"
)

const (
	// StatusUpdateEvent is called when a task changes state
	StatusUpdateEvent = "status_update_event"
)

var (
	addStatuses    = []string{"TASK_RUNNING"}
	removeStatuses = []string{"TASK_FINISHED", "TASK_FAILED", "TASK_KILLED", "TASK_LOST"}
)

// MarathonEvent keeps a marathon event data
type MarathonEvent struct {
	EventType  string `json:"eventType"`
	TaskID     string `json:"taskId"`
	TaskStatus string `json:"taskStatus"`
	AppID      string `json:"appId"`
	Host       string `json:"host"`
	Ports      []int  `json:"ports"`
}

func decodeMarathonEvent(req *http.Request) (*MarathonEvent, error) {
	var event MarathonEvent
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&event)

	return &event, err
}

func handleMarathonEvent(req *http.Request) {
	event, err := decodeMarathonEvent(req)

	if err != nil {
		log.Errorf("Error decoding data callback: %v", err)
		return
	}

	if err = handleEvent(event); err != nil {
		log.Errorf("%v", err)
	}
}

func handleEvent(event *MarathonEvent) error {
	if event.EventType == StatusUpdateEvent {
		return handleStatusUpdate(event)
	}

	log.Infof("No action for event type: %s", event.EventType)
	return nil
}

func handleStatusUpdate(event *MarathonEvent) error {
	log.Infof("Received status update: %s", event.TaskStatus)

	for _, status := range addStatuses {
		if event.TaskStatus == status {
			server := getServer(event)

			if server != nil {
				return addVulcandServer(event.AppID, *server)
			}

			return nil
		}
	}

	for _, status := range removeStatuses {
		if event.TaskStatus == status {
			server := getServer(event)
			if server != nil {
				return removeVulcandServer(event.AppID, *server)
			}

			return nil
		}
	}

	log.Infof("No action for status: %s", event.TaskStatus)
	return nil
}

func getServer(event *MarathonEvent) *VulcandServer {
	if len(event.Ports) == 0 {
		log.Infof("No ports for app %s", event.AppID)
		return nil
	}

	port := event.Ports[0]
	url := fmt.Sprintf("http://%s:%d", event.Host, port)
	hostPort := fmt.Sprintf("%s:%d", event.Host, port)

	return &VulcandServer{
		ID:       event.TaskID,
		URL:      url,
		HostPort: hostPort,
	}
}
