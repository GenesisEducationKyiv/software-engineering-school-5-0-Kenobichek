# Weather-Forecast-API
Weather Subscription API – A simple API that lets users subscribe to weather updates for their city. Built for Genesis Software Engineering School 5.0.

---

### Build

```shell script
docker-compose up --build      # compiles the project and starts required services
```


(The compiled binaries live inside the container image; adjust the compose file if you need to mount or copy them out.)

---

### Running Tests

Prerequisites: Go ≥ 1.21 installed.

Command | What it runs
------- | ------------
`go test -v -short ./...` | Unit tests (fast, in-memory)
`go test -v -tags=integration ./tests/integration/...` | Integration tests (needs deps)
`go test -v -tags=e2e ./tests/e2e/...` | End-to-End tests
`go test -v ./...` | Everything

CI runs the same sequence: unit → integration → e2e.