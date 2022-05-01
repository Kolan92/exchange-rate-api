# Rate Exchange API

Databases stores all rates as nullable `NUMERIC(15, 6)`. I chose that type as all exchange rates from \*.csv are covered by it. If need it can be easily changed to hold values with higher precision.

User can insert null values to database. Null exchange rates are also included in endpoints returning more than one exchange rate.
Only when user request last exchange rate, the last non-null value is returned.

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
