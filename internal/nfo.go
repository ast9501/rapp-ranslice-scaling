package internal

import (
	"log"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ScaleRequest struct {
	ScaleType string `json:"scaleType"`
	ScaleOp   string `json:"scaleOp"`
}

func scaleNs (id string) {
	// Prepare payload
	payload := ScaleRequest{
		ScaleType: "SCALE_NS",
		ScaleOp:   "scale-out",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	// Set url and payload
	var c Conf
	c.ReadConf()
	url := fmt.Sprintf("http://%s/nslcm/v1/ns_instances/%s/scale/", c.NfvoUrl, id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}

	// Set header
	req.Header.Set("Content-Type", "application/json")

	// Send POST Request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP request error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {
		log.Println("Scale out success")
	} else {
		log.Println("Failed to scale out")
	}
}