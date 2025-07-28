package main

import (
	"app/shared"
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

// @@@SNIPSTART money-transfer-project-template-go-start-workflow
func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := shared.RequestDetails{
		SegmentID: 1,
	}

	options := client.StartWorkflowOptions{
		ID:        "submit-segment-1",
		TaskQueue: shared.QueueName,
	}

	log.Printf(
		"Starting submit segment workflow with: %d",
		input.SegmentID,
	)

	we, err := c.ExecuteWorkflow(
		context.Background(),
		options,
		"SubmitSegmentWorkflow",
		input,
	)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}

// @@@SNIPEND
