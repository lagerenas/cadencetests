# Reminders
These tests look at options for sending reminders for a meeting and handeling changes to the reminder time.


## To run test

### Start the cadence server
cd to your copy of github.com/uber/cadence/docker
run `docker-compose up`

### Create a domain called `test-signals`
cd to your copy of github.com/uber/cadence/docker
run `docker run --network=host --rm ubercadence/cli:master --do test-signals domain register -rd 1`

View the workflow in Cadence UI at http://localhost:8088/domain/test-signals/workflows

### Run the server
From this directory run `go run server/main.go`

### Signal the workflow
http://localhost:8090/?eventID=3&minutes=17
