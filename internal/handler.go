package internal

import (
	"encoding/json"
	"log"
)

//TODO: keep the ran registration data at external DB
var regSlices []string

// Process 3GPP define Event type
func Process3gppEventContent(content []string) {
	for _, event := range content {
		var m map[string]interface{}
		err := json.Unmarshal([]byte(event), &m)
		if err != nil {
			log.Println("Failed to unmarshal to map[string]interface")
		}
		log.Println(m)
		switch m["event"].(map[string]interface{})["commonEventHeader"].(map[string]interface{})["domain"].(string) {
		case "fault":
			log.Println("Assert into fault event")
			Handle3gppFault(m)
		default:
			log.Println("Unsupport event type")
		}
	}
}

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
			HandleRanRegistration(m)
		case "VnfMetrics":
			log.Println("Assert into VnfMetricsEvent")
			HandleVnfMetrics()
		default:
			log.Println("Unsupport event type, event body: ", event)
		}
	}
}

func HandleRanRegistration(m map[string]interface{}) {
	log.Println("Processing SliceRegistration Event")
	log.Println(m["NsId"].(string))
	regSlices = append(regSlices, m["NsId"].(string))
}

func HandleVnfMetrics() {
	log.Println("Processing VnfMetrics Event")
}

func Handle3gppFault(m map[string]interface{}) {
	log.Println("Processing 3GPP Fault Event")
	sourceName := m["event"].(map[string]interface{})["commonEventHeader"].(map[string]interface{})["sourceName"].(string)
	log.Println(sourceName)
	problem := m["event"].(map[string]interface{})["faultFields"].(map[string]interface{})["specificProblem"].(string)
	log.Println(problem)

	// TODO: Select NS to scale-out
	log.Println("Determine to Scale-out the slice: ", regSlices[0])
	scaleNs(regSlices[0])
}