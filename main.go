package main

import (
	"log"

	internal "github.com/ast9501/rapp-ranslice-scaling/internal"
	dmaapService "github.com/ast9501/rapp-ranslice-scaling/pkg/dmaap"
	"context"
	"os"
	"os/signal"
	"fmt"
	"syscall"
	"time"
)

func listenOnReg(client dmaapService.SmoDmaap, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Registration Listener is shutting down...")
			return
		default:
			content, err := client.ConsumeTopic("nrtric", "5000", "100")
			if err != nil {
				log.Println("Failed to consume topic")
			}
			//log.Println(content)
			internal.ProcessContent(content)
			time.Sleep(2 * time.Second)
		}
	}
}

func listenOnFault(vesClient dmaapService.SmoDmaap, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Fault Listener is shutting down...")
			return
		default:
			vesContent, err := vesClient.ConsumeTopic("nrtric", "5000", "100")
			if err != nil {
				log.Println("Failed to consume topic")
			}
			//log.Println(vesContent)
			internal.Process3gppEventContent(vesContent)
			time.Sleep(2 * time.Second)
		}
	}
}

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

	/*
	fakeRegistrationEvent := dmaapService.NewSliceRegistrationEvent("NsId", "NsPkgId", "ExpDataDl", "ExpDataUl")
	status, err = client.PublishEvent(fakeRegistrationEvent)
	if err != nil {
		log.Println("Failed to publish evnet: ", status)
	}
	*/

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Create Ran Slice Registration Listener
	go listenOnReg(client, ctx)

	var vesClient dmaapService.SmoDmaap
	vesClient.Topic = dmaapService.NewTopic("unauthenticated.SEC_FAULT_OUTPUT", "VES Fault Event topic", 1, 1, "false")
	vesClient.DmaapUrl = c.DmaapUrl

	// Create 3GPP Fault Event Listener
	go listenOnFault(vesClient, ctx)


	// Listen on shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// wait
	<-sigCh

	// Send shutdown signal to goroutines
	cancel()
	fmt.Println("Shutting down gracefully...")
	time.Sleep(3 * time.Second)
	
	//log.Printf("%s", content)
}
