# Cross Server Cadence Workflows
This experiment is to see how we can split the processing of a workflow between multiple servers.

In the mocked senario clientA is one of several microserivces that make use of common processing within sharedService.
When clientA asks sharedService to do some processing it passes a "function pointer" to a function in clientA that should be run in the middle of the workflow by sharedService. 




## To run test

### Start the cadence server
cd to your copy of github.com/uber/cadence/docker
run `docker-compose up`

### Create a domain called `test-domain`
cd to your copy of github.com/uber/cadence/docker
run `docker run --network=host --rm ubercadence/cli:master --do test-domain domain register -rd 1`

View the workflow in Cadence UI at http://localhost:8088/domain/test-domain/workflows