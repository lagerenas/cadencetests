package main

import (
	"net/http"

	"github.com/lagerenas/cadencetests/sharedServices/internal"
)

func main() {
	http.HandleFunc("/start", internal.StartProcessor)
	http.ListenAndServe(":8090", nil)
}
