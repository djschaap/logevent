# logevent

[![Build Status](https://travis-ci.com/djschaap/logevent.svg?branch=master)](https://travis-ci.com/djschaap/logevent)

A collection of libraries/modules for Splunk/message queue integration.

## Run Tests
```bash
go test ./...

go test -coverprofile=coverage.out ./... \
  && go tool cover -html=coverage.out
```
