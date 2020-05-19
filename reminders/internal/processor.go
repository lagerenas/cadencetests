package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func init() {
}

var RS ReminderSender

func Processor(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("Start processor\n")

	params := r.URL.Query()
	fmt.Fprintf(w, "Params: %+v\n", params)
	eventID := params.Get("eventID")
	minutes, _ := strconv.Atoi(params.Get("minutes"))

	e := Event{
		ID:          eventID,
		Start:       time.Now().Add(time.Duration(minutes) * time.Minute),
		End:         time.Now().Add(time.Duration(minutes+1) * time.Minute),
		Cancelled:   false,
		Description: "Event created",
	}
	fmt.Fprintf(w, "Event: %+v\n", e)

	err := RS.CreateReminder(r.Context(), e)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}

}
