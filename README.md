# logevent

A collection of libraries/modules for Splunk/MQ integration.

## Run Tests
```bash
go test ./...

go test -coverprofile=coverage.out ./... \
  && go tool cover -html=coverage.out
```
