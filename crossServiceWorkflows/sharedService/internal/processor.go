package internal

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"go.uber.org/cadence/client"
	"go.uber.org/cadence/workflow"
)

func init() {
	workflow.Register(Process)
	workflow.RegisterWithOptions(t, workflow.RegisterOptions{Name:"SharedWorkflow"})
}

func StartProcessor(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Start processor\n")

	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	fmt.Fprintf(w, "path is, %s!", getFunctionName(Process))
	var c client.Client 
	c.ExecuteWorkflow()
}

func Process(ctx workflow.Context) error {
	workflow.ChildWorkflowOptions{Domain:"yext"}
	workflow.ExecuteChildWorkflow
	return nil
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
