package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/lagerenas/cadencetests/helper"
	"github.com/lagerenas/cadencetests/reminders/internal"
)

func main() {
	ctx := context.Background()
	fmt.Printf("Starting shared server\n")

	cadenceService, err := helper.NewService("localhost:7933")
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}

	cadenceClient := helper.NewDomainClient(cadenceService, internal.SignalDomain)
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}
	helper.StartCadenceWorker(cadenceService, internal.SignalDomain, internal.SignalDomain)

	fmt.Printf("Cadence running\n")

	signalReminder := internal.NewSignalReminder(cadenceClient)

	eventID := "1"

	signalReminder.CreateReminder(ctx, internal.Event{
		ID:          eventID,
		Start:       time.Now().Add(1 * time.Hour),
		End:         time.Now().Add(2 * time.Hour),
		Cancelled:   false,
		Description: "Event created",
	})

	signalReminder.UpdateReminder(ctx, internal.Event{
		ID:          eventID,
		Start:       time.Now().Add(3 * time.Hour),
		End:         time.Now().Add(4 * time.Hour),
		Cancelled:   false,
		Description: "Event move back 2 hours",
	})

	signalReminder.UpdateReminder(ctx, internal.Event{
		ID:          eventID,
		Start:       time.Now().Add(-1 * time.Hour),
		End:         time.Now().Add(2 * time.Hour),
		Cancelled:   false,
		Description: "move to the past",
	})

	signalReminder.UpdateReminder(ctx, internal.Event{
		ID:          eventID,
		Start:       time.Now().Add(1 * time.Hour),
		End:         time.Now().Add(2 * time.Hour),
		Cancelled:   false,
		Description: "Event moved back to original time",
	})

	signalReminder.CancelReminder(ctx, internal.Event{
		ID:          eventID,
		Start:       time.Now().Add(1 * time.Hour),
		End:         time.Now().Add(2 * time.Hour),
		Cancelled:   true,
		Description: "Event created",
	})

	http.HandleFunc("/", internal.Processor)
	http.ListenAndServe(":8090", nil)

	fmt.Printf("Stopping shared server\n")
}
