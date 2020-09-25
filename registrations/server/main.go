package main

import (
	"context"
	"fmt"
	"os"

	vendastacadence "github.com/vendasta/gosdks/cadence"
	"github.com/vendasta/gosdks/config"
	"go.uber.org/cadence/workflow"
)

func init() {
	workflow.RegisterWithOptions(fakeWorkflow, workflow.RegisterOptions{
		Name: "fakeWorkflow",
	})
}

func fakeWorkflow(ctx workflow.Context) error {
	return nil
}

func main() {
	ctx := context.Background()
	cadence, _, err := vendastacadence.NewService("test", config.Local)
	//defer closer
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}

	cadence.RegisterWorkflowWithAlias("TL1", fakeWorkflow, "fakeWorkflow")
	cadence.RegisterWorkflowWithAlias("TL2", fakeWorkflow, "fakeWorkflow")
	err = cadence.NewWorker(ctx, "test", "TL1")
	if err != nil {
		fmt.Printf("Error starting worker 1: %v\n", err)
		os.Exit(1)
	}
	err = cadence.NewWorker(ctx, "test", "TL2")
	if err != nil {
		fmt.Printf("Error starting worker 2: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ran ok")

}
