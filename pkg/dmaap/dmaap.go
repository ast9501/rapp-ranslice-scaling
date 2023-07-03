package dmaap

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type DmaapClient interface {
	CreateTopic() (int, error)
	PublishEvent(string) (int, error)
	ConsumeTopic() ([]string, error)
}

type DmaapTopic struct {
	TopicName          string `json:"topicName"`
	TopicDescrip       string `json:"topicDescription"`
	PartitionCount     int    `json:"partitionCount"`
	ReplicationCount   int    `json:"replicationCount"`
	TransactionEnabled string `json:"transactionEnabled"`
}

type SmoDmaap struct {
	DmaapUrl string
	Topic    *DmaapTopic
}

type TestEvent struct {
	EventType string `json:"EventType"`
	EventName string `json:"eventName"`
	Value     int    `json:"value"`
}

type SliceRegistrationEvent struct {
	EventType string `json:"EventType"`
	NsId      string `json:"NsId"`
	NsPkgId   string `json:"NsPkgId"`
	ExpDataDl string `json:"expDataRateDl"`
	ExpDataUl string `json:"expDataRateUl"`
}

type VnfMetricsEvent struct {
	EventType string `json:"EventType"`
	VnfId     string `json:"VnfId"`
	VnfType   string `json:"VnfType"`
	Cpu       string `json:"Cpu"`
	Ram       string `json:"Ram"`
	NetOut    string `Json:"NetOut"`
	NetIn     string `Json:"NetIn"`
}

func (s SmoDmaap) CreateTopic() (err error) {
	data, _ := json.Marshal(s.Topic)
	url := s.DmaapUrl + "/topics/create"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Failed to init CreateTopicRequest")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Failed to send request: ", err.Error())
		return
	} else {
		log.Println("Create topic with status: ", resp.StatusCode)
	}
	defer resp.Body.Close()

	return
}

func (s SmoDmaap) PublishEvent(e interface{}) (StatueCode int, err error) {
	var data []byte

	switch e.(type) {
	case *TestEvent:
		log.Println("Prepare to publish TestEvent...")
		//log.Println(e.(*TestEvent).EventName)
		data, err = json.Marshal(e.(*TestEvent))
		if err != nil {
			log.Println("Failed to marshal json!")
			return 0, err
		}

	case *SliceRegistrationEvent:
		log.Println("Prepare to publish SliceRegistrationEvent")
		data, err = json.Marshal(e.(*SliceRegistrationEvent))
		if err != nil {
			log.Println("Failed to marshal json!")
			return 0, err
		}

	case *VnfMetricsEvent:
		log.Println("Prepare to publish VnfMetricsEvent")
		data, err = json.Marshal(e.(*VnfMetricsEvent))
		if err != nil {
			log.Println("Failed to marshal json!")
			return 0, err
		}

	default:
		log.Println("No match event type!")
	}

	url := s.DmaapUrl + "/events/" + s.Topic.TopicName + "/"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Failed to init request")
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request")
		return 0, err
	} else {
		log.Println("Publish event with status: ", resp.StatusCode)
	}
	defer resp.Body.Close()

	return resp.StatusCode, err
}

func (s SmoDmaap) ConsumeTopic(user, timeout, limit string) (content []string, err error) {
	log.Println("Consume Topic")

	url := s.DmaapUrl + "/events/" + s.Topic.TopicName + "/users/" + user + "?timeout=" + timeout + "&limit=" + limit
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println("Failed to init request")
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to send request")
		return
	} else {
		log.Println("Consume topic with status: ", resp.StatusCode)
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body")
		return
	}

	// unmarshal response into string slice
	// i.e.:
	// string(response) is ["{\"eventName\":\"test\",\"value\":1}","{\"eventName\":\"test\",\"value\":1}"]
	// [{"eventName":"test","value":1} {"eventName":"test","value":1}]

	var strSlice []string
	err = json.Unmarshal(response, &strSlice)
	if err != nil {
		log.Println("Failed to unmarshal to strSlice, ", err.Error())
	}

	content = strSlice
	return
}

func NewTopic(TopicName string, TopicDescrip string, PartitionCount int, ReplicationCount int, TransactionEnabled string) (t *DmaapTopic) {
	topic := DmaapTopic{
		TopicName:          TopicName,
		TopicDescrip:       TopicDescrip,
		PartitionCount:     PartitionCount,
		ReplicationCount:   ReplicationCount,
		TransactionEnabled: TransactionEnabled,
	}
	t = &topic
	return
}

func NewTestEvent(EventName string, Value int) (e *TestEvent) {
	e = &TestEvent{
		EventType: "Test",
		EventName: EventName,
		Value:     Value,
	}

	return
}

func NewSliceRegistrationEvent(NsId, NsPkgId, ExpDataDl, ExpDataUl string) (e *SliceRegistrationEvent) {
	e = &SliceRegistrationEvent{
		EventType: "SliceRegistration",
		NsId:      NsId,
		NsPkgId:   NsPkgId,
		ExpDataDl: ExpDataDl,
		ExpDataUl: ExpDataUl,
	}

	return e
}

func NewVnfMetricEvent(VnfId, VnfType, CpuRate, Ram, NetIn, NetOut string) (e *VnfMetricsEvent) {
	e = &VnfMetricsEvent{
		EventType: "VnfMetric",
		VnfId:     VnfId,
		VnfType:   VnfType,
		Cpu:       CpuRate,
		Ram:       Ram,
		NetIn:     NetIn,
		NetOut:    NetOut,
	}

	return e
}
