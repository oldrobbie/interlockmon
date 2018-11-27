package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZachtimusPrime/Go-Splunk-HTTP/splunk"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	// curl -k -u admin:yourpassword https://localhost:8089/servicesNS/admin/splunk_httpinput/data/inputs/http -d name=mytoken
	// application/x-www-form-urlencoded

	enableHEC()
	token := createHECToken()

	splunk := splunk.NewClient(
		nil,
		os.Getenv("SPLUNK_HEC_URL")+"/services/collector",
		token,
		"events",
		"json",
		"main",
	)

	ctx := context.Background()

	// Create Docker API client
	client, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("error creating docker client: %s\n", err.Error())
	}

	// Listen for Docker events
	eventFilter := filters.NewArgs()
	//eventFilter.Add("type", events.ContainerEventType)
	eventChan, errorChan := client.Events(ctx, types.EventsOptions{
		Filters: eventFilter,
	})

	// Receive SIGTERM no channel
	stopSignalChan := make(chan os.Signal, 1)
	signal.Notify(stopSignalChan, os.Interrupt, syscall.SIGTERM)

	// Main event loop
	for {
		select {
		case <-stopSignalChan:
			return
		case err := <-errorChan:
			log.Printf("error receiving engine events: %s\n", err.Error())
			return
		case event := <-eventChan:
			err := splunk.Log(event)
			if err != nil {
				log.Fatalf("error sending log: %s\n", err.Error())
			}
			//fmt.Println(string(e))
			//log.Printf("ID=\"%s\" From=\"%s\" Type=\"%s\" Action=\"%s\" Actor=\"%v\"\n", event.ID, event.From, event.Type, event.Action, event.Actor)
		}
	}
}
