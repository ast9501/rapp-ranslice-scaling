package internal

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type registrationRequestBody struct {
	AppVer   string `json:"version"`
	Name     string `json:"display_name"`
	Descript string `json:"description"`
}

type registrationResponse struct {
	Detail string `json:"detail"`
	Status int32  `json:"status"`
}

// Register rapp to nonrt-ric catalogue service
func Register(CatalogueServiceUrl string) {
	url := CatalogueServiceUrl
	data, err := encapRegistrationRequest("1.0.0", "RAN-slice-scaling-control", "Analyze RAN slice metrics and determine the scale of slice")
	if err != nil {
		log.Printf("Failed to generate registration request")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to init request to catalogue service")
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error occured while PUT rapp registration")
		log.Printf(err.Error())
	} else {
		if resp.StatusCode != 201 {
			//TODO: handle the unexpected status code
			respBody, _ := ioutil.ReadAll(resp.Body)
			var response registrationResponse
			json.Unmarshal(respBody, &response)
			log.Printf("Registration failed: %v", response)
		} else {
			log.Printf("Register rapp with status %d", resp.StatusCode)
			log.Println("The rapp was registered at url: ", resp.Header["Location"])
		}
	}

	defer resp.Body.Close()
}

func encapRegistrationRequest(appVersion, displayName, description string) (jsonReq []byte, err error) {
	req := registrationRequestBody{
		AppVer:   appVersion,
		Name:     displayName,
		Descript: description,
	}

	jsonReq, err = json.Marshal(req)
	return
}
