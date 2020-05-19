package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lagerenas/cadencetests/helper"
	"github.com/lagerenas/cadencetests/reminders/internal"
)

func main() {
	//ctx := context.Background()
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
	internal.RS = signalReminder
	/*
		eventID := "2"

		err = signalReminder.CreateReminder(ctx, internal.Event{
			ID:          eventID,
			Start:       time.Now().Add(1 * time.Hour),
			End:         time.Now().Add(2 * time.Hour),
			Cancelled:   false,
			Description: "Event created",
		})
		if err != nil {
			fmt.Printf("error creating reminder: %v\n", err)
		}

		err = signalReminder.UpdateReminder(ctx, internal.Event{
			ID:          eventID,
			Start:       time.Now().Add(3 * time.Hour),
			End:         time.Now().Add(4 * time.Hour),
			Cancelled:   false,
			Description: "Event move back 2 hours",
		})
		if err != nil {
			fmt.Printf("error creating reminder: %v\n", err)
		}

		err = signalReminder.UpdateReminder(ctx, internal.Event{
			ID:          eventID,
			Start:       time.Now().Add(-1 * time.Hour),
			End:         time.Now().Add(2 * time.Hour),
			Cancelled:   false,
			Description: "move to the past",
		})
		if err != nil {
			fmt.Printf("error creating reminder: %v\n", err)
		}

		err = signalReminder.UpdateReminder(ctx, internal.Event{
			ID:          eventID,
			Start:       time.Now().Add(1 * time.Hour),
			End:         time.Now().Add(2 * time.Hour),
			Cancelled:   false,
			Description: "Event moved back to original time",
		})
		if err != nil {
			fmt.Printf("error creating reminder: %v\n", err)
		}

		err = signalReminder.CancelReminder(ctx, internal.Event{
			ID:          eventID,
			Start:       time.Now().Add(1 * time.Hour),
			End:         time.Now().Add(2 * time.Hour),
			Cancelled:   true,
			Description: "Event created",
		})
		if err != nil {
			fmt.Printf("error creating reminder: %v\n", err)
		}
	*/

	http.HandleFunc("/", internal.Processor)
	http.ListenAndServe(":8090", nil)

	fmt.Printf("Stopping shared server\n")
}
