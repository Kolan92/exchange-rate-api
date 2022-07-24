# Rate Exchange API

Simple API written in GO which allows to retrive currency exchange rates between curencies.

## Run app

App can run in docker compose or in standalone mode.
For both ways swagger UI is available here: http://localhost:8081/swagger/index.html

### Docker compose

Just run `docker compose up` in root directory

### Standalone

Those env variables are required to run the app. If postgres container port 5432 is exposed, api can connect to it when it is executed outside of docker.

```bash
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=root
export DB_HOST=0.0.0.0
export DB_PORT=5432
```

Command: `go run .` from exchange-rate-api directory

## Run tests

Run those commands from exchange-rate-api directory

- Unit tests: `go test ./... -short -v`, skips integration tests.
- Integration tests: `go test github.com/kolan92/exchange-rate-api/integration-tests` - requires to have running api on localhost:8081
