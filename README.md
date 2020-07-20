# logevent

[![Build Status](https://travis-ci.com/djschaap/logevent.svg?branch=master)](https://travis-ci.com/djschaap/logevent)

A collection of packages for Splunk/message queue integration.

## Run Tests
```bash
# standalone unit tests
go test ./...

# coverage
go test -coverprofile=coverage.out ./... \
  && go tool cover -html=coverage.out
```

## send CLI

The send executable is included as a sample tool to send messages.

Environment variables are used for output/package configuration.
Command-line arguments are used for message/event-specific properties.

Setting boolean variables (TRACE, etc.) to ANYTHING other than the
empty string will be interpreted as true.

### senddump Package

Dump message to stderr for debug purposes.

```bash
TRACE=x go run cmd/send/main.go \
  "bare message"

TRACE=x go run cmd/send/main.go \
  -customer abc -host h1 -index main \
  -source s -sourceenvironment se -sourcetype st \
  -epoch $(date +%s) -field a=A -field b="indexed event field B" \
  "with integer time and indexed event fields"

TRACE=x go run cmd/send/main.go \
  -time 2020-01-01T00:00:00Z \
  "message with UTC time"

TRACE=x go run cmd/send/main.go \
  -time 2020-01-01T12:00:00+06:00 \
  "message with time offset"
```

### sendhec Package

Send message directly to Splunk HTTP Event Collector (HEC).

```bash
export HEC_URL=https://localhost:8088
export HEC_TOKEN=00000000-0000-0000-0000-000000000000
export HEC_INSECURE=true
PACKAGE=sendhec TRACE=x go run cmd/send/main.go \
  -host h2 \
  "message with host"
```

### sendsns Package

Send message to Amazon SNS topic.

```bash
export AWS_ACCESS_KEY_ID=xxx
export AWS_REGION=us-east-1
export AWS_SECRET_ACCESS_KEY=xxx
export TOPIC=arn:xxx
PACKAGE=sendsns TRACE=x go run cmd/send/main.go \
  -host h2 \
  "message with host"
```
