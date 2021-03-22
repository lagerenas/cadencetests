package internal

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/cadence"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/activity"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/workflow"
)

const CadenceDomain = "test"

var CadenceClient client.Client

func init() {
	workflow.Register(reindexOrdersWorkflow)
	activity.Register(registerIndexActivity)
	activity.Register(reindexOrdersActivity)
}

func Processor(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Start processor\n")

	params := r.URL.Query()
	fmt.Fprintf(w, "Params: %+v\n", params)
	eventID := params.Get("eventID")
	minutes, _ := strconv.Atoi(params.Get("minutes"))
	err := ReindexOrders(r.Context(), CadenceClient)
	fmt.Fprintf(w, "%v: %v Error: %v", eventID, minutes, err)

}

func ReindexOrders(ctx context.Context, CadenceClient client.Client) error {

	options := client.StartWorkflowOptions{
		ID:                           "reindex-orders",
		TaskList:                     CadenceDomain,
		ExecutionStartToCloseTimeout: time.Hour * 24,
		WorkflowIDReusePolicy:        client.WorkflowIDReusePolicyAllowDuplicate,
	}
	exe, err := CadenceClient.StartWorkflow(ctx, options, reindexOrdersWorkflow)
	if err != nil {
		alreadyStarted, ok := err.(*shared.WorkflowExecutionAlreadyStartedError)
		if ok {
			fmt.Printf("Reindex workflow already started (%s) for %s:%s\n", *alreadyStarted.Message, *alreadyStarted.RunId, *alreadyStarted.StartRequestId)
			return nil
		}
		fmt.Printf("Error starting reindex cadence workflow: %s\n", err.Error())
		return err
	}
	fmt.Printf("Reindexing scheduled as cadence workflow: %s/%s\n", exe.ID, exe.RunID)
	return nil
}

func reindexOrdersWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: time.Minute,
		StartToCloseTimeout:    2 * time.Minute,
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    1 * time.Minute,
			// Once the vstore cursor expires we will need to start from the beginning
			ExpirationInterval:       1 * time.Minute,
			MaximumAttempts:          0,
			NonRetriableErrorReasons: nil,
		},
		HeartbeatTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	fmt.Printf("In reindexOrdersWorkflow\n")

	if err := workflow.ExecuteActivity(ctx, registerIndexActivity).Get(ctx, nil); err != nil {
		return err
	}

	if err := workflow.ExecuteActivity(ctx, reindexOrdersActivity).Get(ctx, nil); err != nil {
		return err
	}
	return nil
}

func registerIndexActivity(ctx context.Context) error {
	fmt.Printf("Running registerIndexActivity\n")
	return nil
}

func reindexOrdersActivity(ctx context.Context) error {
	fmt.Printf("Running reindexOrdersActivity\n")
	time.Sleep(15 * time.Second)
	fmt.Printf("Done Sleeping reindexOrdersActivity\n")
	return fmt.Errorf("fake error")
}
