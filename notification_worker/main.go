package main

import (
	"app/shared"
	"context"
	"errors"
	"log"
	"math/rand"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, shared.NotificationQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterActivity(SendNotification)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func SendNotification(ctx context.Context, data shared.RequestDetails) (string, error) {
	log.Printf("Sending notification for segment %d", data.SegmentID)
	isError := rand.Intn(2)
	if isError%2 == 0 {
		return "", errors.New("notification failed")
	}
	return "Notification sent", nil
}
