package sendamqp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/djschaap/logevent"
	"github.com/kr/pretty"
	"github.com/streadway/amqp"
	"log"
)

// Sess stores sendamqp session state.
type Sess struct {
	amqpChan       *amqp.Channel
	amqpConn       *amqp.Connection
	amqpError      chan *amqp.Error
	amqpExchange   string
	amqpRoutingKey string
	amqpURL        string
	trace          bool
}

// CloseSvc closes the open session.
// CloseSvc must not be called when no session is open.
func (sender *Sess) CloseSvc() error {
	if sender.amqpConn == nil {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	sender.amqpChan = nil
	sender.amqpConn.Close()
	sender.amqpConn = nil
	return nil
}

// OpenSvc opens a new session.
// OpenSvc must not be called when a session is already open.
func (sender *Sess) OpenSvc() error {
	if sender.amqpConn != nil || sender.amqpChan != nil {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	conn, err := amqp.Dial(sender.amqpURL)
	if err != nil {
		return fmt.Errorf("amqp.Dial() failed: %v", err)
	}
	sender.amqpConn = conn
	sender.amqpError = conn.NotifyClose(make(chan *amqp.Error))

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp.Connection.Channel() failed: %v", err)
	}
	sender.amqpChan = ch

	// TODO create exchange here, if desired

	return nil
}

// SendMessage sends a LogEvent to a RabbitMQ (AMQP) exchange.
func (sender *Sess) SendMessage(logEvent logevent.LogEvent) error {
	if sender.amqpChan == nil {
		return errors.New("SendMessage() called before OpenSvc()")
	}

	select {
	case err := <-sender.amqpError:
		return fmt.Errorf("AMQP connection closed unexpectedly: %s", err)
	default:
	}

	amqpMessage := sender.buildAmqpMessage(logEvent)
	sender.amqpChan.Publish(
		sender.amqpExchange,
		sender.amqpRoutingKey,
		false, // mandatory
		false, // immediate
		amqpMessage,
	)
	sender.tracePretty("TRACE_SENDAMQP amqpMessage:", amqpMessage,
		"\nBody:", string(amqpMessage.Body))
	return nil
}

// SetTrace enables tracing, which dumps all messages to stderr.
func (sender *Sess) SetTrace(v bool) {
	sender.trace = v
}

func (sender *Sess) buildAmqpMessage(logEvent logevent.LogEvent) amqp.Publishing {
	attr := logEvent.Attributes
	headers := make(map[string]interface{})
	if attr.CustomerCode != "" {
		headers["customer_code"] = attr.CustomerCode
	}
	if attr.Host != "" {
		headers["host"] = attr.Host
	}
	if attr.Source != "" {
		headers["source"] = attr.Source
	}
	if attr.SourceEnvironment != "" {
		headers["source_environment"] = attr.SourceEnvironment
	}
	if attr.Sourcetype != "" {
		headers["sourcetype"] = attr.Sourcetype
	}
	if attr.Type != "" {
		headers["type"] = attr.Type
	}
	messageJSONBytes, _ := json.Marshal(logEvent.Content)
	amqpMessage := amqp.Publishing{
		Headers:         headers,
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            messageJSONBytes,
		DeliveryMode:    amqp.Persistent,
		Priority:        0,
	}
	if attr.Type != "" {
		amqpMessage.Type = attr.Type
	}
	return amqpMessage
}

func (sender *Sess) tracePretty(
	args ...interface{},
) {
	if sender.trace {
		pretty.Log(args...)
	}
}

func (sender *Sess) tracePrintln(
	args ...interface{},
) {
	if sender.trace {
		log.Println(args...)
	}
}

// New creates a new sendhec object/session.
// It requires an AMQP URL/URI (which may contain username, password, host, port, and/or vhost), exchange name, and routing key.
func New(amqpURL, amqpExchange, amqpRoutingKey string) *Sess {
	sess := Sess{
		amqpExchange:   amqpExchange,
		amqpRoutingKey: amqpRoutingKey,
		amqpURL:        amqpURL,
	}
	return &sess
}
