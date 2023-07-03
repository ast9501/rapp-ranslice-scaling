package internal

import (
	"encoding/json"
	"log"
)

func ProcessContent(content []string) {
	for _, event := range content {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(event), &m)
		if err != nil {
			log.Println("Failed to unmarshal to map[string]interface")
		}

		switch m["EventType"].(string) {
		case "SliceRegistration":
			log.Println("Assert into SliceRegistrationEvent")
			HandleRanRegistration()
		case "VnfMetrics":
			log.Println("Assert into VnfMetricsEvent")
			HandleVnfMetrics()
		default:
			log.Println("Unsupport event type, event body: ", event)
		}
	}
}

func HandleRanRegistration() {
	log.Println("Processing SliceRegistration Event")
}

func HandleVnfMetrics() {
	log.Println("Processing VnfMetrics Event")
}
