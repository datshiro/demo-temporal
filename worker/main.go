package main

import (
	"app/worker/activity"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()
	queueName := "export-segment"

	w := worker.New(c, queueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterActivity(activity.RefreshSegment)
	w.RegisterActivity(activity.ConvertSegment)
	w.RegisterActivity(activity.ExportSegment)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
