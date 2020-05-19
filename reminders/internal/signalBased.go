package internal

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/cadence/activity"
	cadenceclient "go.uber.org/cadence/client"
	"go.uber.org/cadence/workflow"
)

const SignalDomain = "test-signals"

func init() {
	workflow.Register(signalWorkflow)
	activity.Register(sendReminderActivity)
}

type SignalReminder struct {
	CadenceClient cadenceclient.Client
}

func NewSignalReminder(cc cadenceclient.Client) *SignalReminder {
	return &SignalReminder{
		CadenceClient: cc,
	}
}

func (sr *SignalReminder) CreateReminder(ctx context.Context, m Event) error {
	err := sr.setReminder(ctx, m, 1*time.Minute)
	if err != nil {
		return err
	}
	//return sr.setReminder(ctx, m, 24*time.Hour)
	return nil
}

func (sr *SignalReminder) UpdateReminder(ctx context.Context, m Event) error {
	err := sr.setReminder(ctx, m, 15*time.Minute)
	if err != nil {
		return err
	}
	return sr.setReminder(ctx, m, 24*time.Hour)
}

func (sr *SignalReminder) CancelReminder(ctx context.Context, m Event) error {
	err := sr.setReminder(ctx, m, 15*time.Minute)
	if err != nil {
		return err
	}
	return sr.setReminder(ctx, m, 24*time.Hour)
}

func (sr *SignalReminder) setReminder(ctx context.Context, m Event, period time.Duration) error {
	fmt.Printf("Event in set reminder %#v", m)
	workflowID := fmt.Sprintf("reminder-for-%v-%v", m.ID, period)
	remindAt := m.Start.Add(-1 * period)
	if m.Cancelled || remindAt.Before(time.Now()) {
		err := sr.CadenceClient.CancelWorkflow(ctx, workflowID, "")
		return err
	} else {

		workflowOptions := cadenceclient.StartWorkflowOptions{
			ID:                           workflowID,
			TaskList:                     SignalDomain,
			ExecutionStartToCloseTimeout: time.Hour * 24 * 365 * 5, // Kill the task if it has not completed within 5 years
		}
		//Using SignalWithStartWorkflow instead of StartWorkflow will take care of doing a create or update
		_, err := sr.CadenceClient.SignalWithStartWorkflow(ctx, workflowID, "RemindAt", remindAt, workflowOptions, signalWorkflow, m.ID)
		return err
	}
}

func signalWorkflow(ctx workflow.Context, eventID string) error {
	logger := workflow.GetLogger(ctx)
	remindAtCh := workflow.GetSignalChannel(ctx, "RemindAt")

	var remindAt time.Time
	remindAtCh.Receive(ctx, &remindAt)

	timerFired := false
	for !timerFired {
		delay := remindAt.Sub(workflow.Now(ctx))

		selector := workflow.NewSelector(ctx)

		logger.Sugar().Infof("Setting up a timer to fire after: %v", delay)
		timerCancelCtx, cancelTimerHandler := workflow.WithCancel(ctx)
		timerFuture := workflow.NewTimer(timerCancelCtx, delay)
		selector.AddFuture(timerFuture, func(f workflow.Future) {
			logger.Info("Timer Fired.")
			timerFired = true
		})

		selector.AddReceive(remindAtCh, func(c workflow.Channel, more bool) {
			logger.Info("RemindAt signal received.")
			logger.Info("Cancel outstanding timer.")
			cancelTimerHandler()

			c.Receive(ctx, &remindAt)
			logger.Sugar().Infof("Update remind at to: %v", remindAt)
		})

		logger.Info("Waiting for timer to fire.")
		selector.Select(ctx)
	}

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToStartTimeout: 3 * time.Second,
		StartToCloseTimeout:    1 * time.Minute,
	})
	err := workflow.ExecuteActivity(ctx, sendReminderActivity, eventID).Get(ctx, nil)
	if err != nil {
		return err
	}

	workflow.GetLogger(ctx).Info("Workflow completed.")

	return nil
}

func sendReminderActivity(ctx context.Context, eventID string) error {
	fmt.Printf("Sending reminder for %v\n", eventID)
	return nil
}
