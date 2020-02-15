package internal

import (
	"fmt"
	"net/http"
)

func StartProcessor(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
