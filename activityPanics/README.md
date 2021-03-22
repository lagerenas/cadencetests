# Panic tests
Find out what happens when a workflow panics


## To run test

### Start the cadence server
cd to your copy of github.com/uber/cadence/docker
run `docker-compose up`

### Create a domain called `test-domain`
cd to your copy of github.com/uber/cadence/docker
run `docker run --network=host --rm ubercadence/cli:master --do test-domain domain register -rd 1`

View the workflow in Cadence UI at http://localhost:8088/domain/test-domain/workflows

Start the server
run `go run server/*.go`

Trigger new workflows to start by opening a web browser to this URL. Adjust params as needed.
http://localhost:8090/?asyncWorkflow=true&panicWorkflow=false&panicActivity=200&localActivity=false
