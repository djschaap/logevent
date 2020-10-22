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

# sendamqp integration tests
docker run --rm -p 5672:5672 -p 15672:15672 rabbitmq:management-alpine
AMQP_URL=amqp://guest:guest@localhost:5672 \
  go test -tags amqp ./...
```

## send CLI

The send executable is included as a sample tool to send messages.

Environment variables are used for output/package configuration.
Command-line arguments are used for message/event-specific properties.

Currently, setting boolean variables (`SENDER_TRACE`, etc.) to ANYTHING
other than the empty string will be interpreted as true.
This may change without notice; using zero, "N", or similar to
represent true is NOT recommended.

### sendamqp Package

Send message to RabbitMQ exchange.
Exchange must already exist.
`AMQP_ROUTING_KEY` may be meaningless when using a headers exchange,
but some value must still be provided.
`AMQP_TTL` is specified in seconds (default is no TTL).

```bash
export AMQP_URL=amqp://guest:guest@localhost:5672
export AMQP_EXCHANGE=amq.headers
export AMQP_ROUTING_KEY=the_weather
export AMQP_TTL=60
SENDER_PACKAGE=sendamqp SENDER_TRACE=x go run cmd/send/main.go \
  -host h2 \
  "message with host"
```

### senddump Package

Dump message to stderr for debug purposes.
Default package when `SENDER_PACKAGE` is not set.

```bash
SENDER_TRACE=x go run cmd/send/main.go \
  "bare message"

SENDER_TRACE=x go run cmd/send/main.go \
  -customer abc -host h1 -index main \
  -source s -sourceenvironment se -sourcetype st \
  -epoch $(date +%s) -field a=A -field b="indexed event field B" \
  "with integer time and indexed event fields"

SENDER_TRACE=x go run cmd/send/main.go \
  -time 2020-01-01T00:00:00Z \
  "message with UTC time"

SENDER_TRACE=x go run cmd/send/main.go \
  -time 2020-01-01T12:00:00+06:00 \
  "message with time offset"
```

### sendhec Package

Send message directly to Splunk HTTP Event Collector (HEC).
`HEC_TOKEN` is required.

```bash
export HEC_URL=https://localhost:8088
export HEC_TOKEN=00000000-0000-0000-0000-000000000000
export HEC_INSECURE=true
SENDER_PACKAGE=sendhec SENDER_TRACE=x go run cmd/send/main.go \
  -host h2 \
  "message with host"
```

### sendsns Package

Send message to Amazon SNS topic.
`AWS_SNS_TOPIC` is required.

```bash
export AWS_ACCESS_KEY_ID=xxx
export AWS_REGION=us-east-1
export AWS_SECRET_ACCESS_KEY=xxx
export AWS_SNS_TOPIC=arn:xxx
SENDER_PACKAGE=sendsns SENDER_TRACE=x go run cmd/send/main.go \
  -host h2 \
  "message with host"
```
