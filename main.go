package main

import (
	"log"

	internal "github.com/ast9501/rapp-ranslice-scaling/internal"
	dmaapService "github.com/ast9501/rapp-ranslice-scaling/pkg/dmaap"
)

func main() {
	var c internal.Conf
	c.ReadConf()
	internal.Register(c.CatalogueServiceUrl)

	var client dmaapService.SmoDmaap
	client.Topic = dmaapService.NewTopic(c.DmaapTopic, "RAN Slice metrics and state sync", 1, 1, "false")
	client.DmaapUrl = c.DmaapUrl
	err := client.CreateTopic()
	if err != nil {
		log.Println("Failed to create topic: ", err.Error())
	}

	event := dmaapService.NewTestEvent("test", 1)
	status, err := client.PublishEvent(event)
	if err != nil {
		log.Println("Failed to publish evnet: ", status)
	}

	fakeRegistrationEvent := dmaapService.NewSliceRegistrationEvent("NsId", "NsPkgId", "ExpDataDl", "ExpDataUl")
	status, err = client.PublishEvent(fakeRegistrationEvent)
	if err != nil {
		log.Println("Failed to publish evnet: ", status)
	}

	content, err := client.ConsumeTopic("nrtric", "3000", "100")
	if err != nil {
		log.Println("Failed to consume topic")
	}

	internal.ProcessContent(content)
	//log.Printf("%s", content)
}
