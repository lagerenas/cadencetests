package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/lagerenas/cadencetests/helper"
	"github.com/lagerenas/cadencetests/sharedService/internal"
)

func main() {
	fmt.Printf("Starting shared server\n")

	cadenceClient, err := helper.NewService("localhost:7933")
	if err != nil {
		fmt.Printf("Error starting cadence client: %v\n", err)
		os.Exit(1)
	}
	helper.StartCadenceWorker(cadenceClient, "test-domain", "messanger")

	fmt.Printf("Cadence running\n")

	http.HandleFunc("/start", internal.StartProcessor)
	http.ListenAndServe(":8090", nil)

	fmt.Printf("Stopping shared server\n")
}
