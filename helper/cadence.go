package helper

import (
	"github.com/uber-go/tally"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/worker"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
)

// NewService returns a new workflowserviceclient.Interface instance.
// This is your connection to the shared Cadence frontend server
func NewService(host string) (workflowserviceclient.Interface, error) {
	var cadenceService = "cadence-frontend"

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(cadenceService))
	if err != nil {
		return nil, err
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: cadenceService,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: ch.NewSingleOutbound(host)},
		},
	})
	if err := dispatcher.Start(); err != nil {
		return nil, err
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(cadenceService)), nil
}

// NewDomainManagmentClient returns a client that can be used to manage the domains on the service
func NewDomainManagmentClient(service workflowserviceclient.Interface) client.DomainClient {
	options := &client.Options{}
	return client.NewDomainClient(service, options)
}

// NewDomainClient returns a client that can be used to manage the tasks within a given domain
func NewDomainClient(service workflowserviceclient.Interface, domain string) client.Client {
	options := &client.Options{}
	return client.NewClient(service, domain, options)
}

// StartCadenceWorker starts up the cadence worker.
func StartCadenceWorker(client workflowserviceclient.Interface, domain string, taskListName string) {
	workerOptions := worker.Options{
		MetricsScope:                   tally.NewTestScope(domain, map[string]string{}),
		NonDeterministicWorkflowPolicy: worker.NonDeterministicWorkflowPolicyFailWorkflow,
	}

	worker := worker.New(client, domain, taskListName, workerOptions)
	worker.Start()
}
