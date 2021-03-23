package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lagerenas/cadencetests/activityPanics/internal"
	"github.com/lagerenas/cadencetests/helper"
)

func main() {
	fmt.Printf("Starting shared server\n")

	cadenceService, err := helper.NewService("localhost:7933")
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}

	cadenceClient := helper.NewDomainClient(cadenceService, internal.CadenceDomain)
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}
	helper.StartCadenceWorker(cadenceService, internal.CadenceDomain, internal.CadenceDomain)

	fmt.Printf("Cadence running\n")

	internal.CadenceClient = cadenceClient

	http.HandleFunc("/workflow", internal.Processor)
	http.ListenAndServe(":8090", nil)

	fmt.Printf("Stopping shared server\n")
}
