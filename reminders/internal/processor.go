package internal

import (
	"fmt"
	"net/http"
)

func init() {
}

func Processor(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Start processor\n")

	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])

}
