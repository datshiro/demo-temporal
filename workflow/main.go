package main

import (
	"app/shared"
	"fmt"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
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
	w.RegisterWorkflow(SubmitSegmentWorkflow)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func SubmitSegmentWorkflow(ctx workflow.Context, input shared.RequestDetails) (string, error) {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500, // 0 is unlimited retries
		NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"},
	}

	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)
	var segmentOutput string

	segmentErr := workflow.ExecuteActivity(ctx, "RefreshSegment", input).Get(ctx, &segmentOutput)

	if segmentErr != nil {
		return "", segmentErr
	}
	segmentErr = workflow.ExecuteActivity(ctx, "ConvertSegment", input).Get(ctx, &segmentOutput)

	if segmentErr != nil {
		return "", segmentErr
	}

	var result string

	err := workflow.ExecuteActivity(ctx, "ExportSegment", input).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	log.Printf("Submit Segment complete: %s", result)

	err = workflow.ExecuteActivity(ctx, "SendNotification", input).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	log.Printf("Send Notification complete: %s", result)

	return fmt.Sprintf("Submit Segment complete: %s", result), nil
}
