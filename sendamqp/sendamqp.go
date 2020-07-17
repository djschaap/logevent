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

type sess struct {
	amqpChan       *amqp.Channel
	amqpConn       *amqp.Connection
	amqpExchange   string
	amqpRoutingKey string
	amqpUrl        string
	trace          bool
}

func (self *sess) CloseSvc() error {
	if self.amqpConn == nil {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	self.amqpChan = nil
	self.amqpConn.Close()
	self.amqpConn = nil
	return nil
}

func (self *sess) OpenSvc() error {
	if self.amqpConn != nil || self.amqpChan != nil {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	conn, err := amqp.Dial(self.amqpUrl)
	if err != nil {
		return fmt.Errorf("amqp.Dial() failed: %v", err)
	}
	self.amqpConn = conn
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp.Connection.Channel() failed: %v", err)
	}
	self.amqpChan = ch

	// TODO create exchange here, if desired

	return nil
}

func (self *sess) SendMessage(logEvent logevent.LogEvent) error {
	if self.amqpChan == nil {
		return errors.New("SendMessage() called before OpenSvc()")
	}

	amqpMessage := self.buildAmqpMessage(logEvent)
	self.amqpChan.Publish(
		self.amqpExchange,
		self.amqpRoutingKey,
		false, // mandatory
		false, // immediate
		amqpMessage,
	)
	self.tracePretty("TRACE_SENDAMQP amqpMessage:", amqpMessage,
		"\nBody:", string(amqpMessage.Body))
	return nil
}

func (self *sess) SetTrace(v bool) {
	self.trace = v
}

func (self *sess) buildAmqpMessage(logEvent logevent.LogEvent) amqp.Publishing {
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
	messageJson_bytes, _ := json.Marshal(logEvent.Content)
	amqpMessage := amqp.Publishing{
		Headers:         headers,
		ContentType:     "application/json",
		ContentEncoding: "",
		Body:            messageJson_bytes,
		DeliveryMode:    amqp.Persistent,
		Priority:        0,
	}
	return amqpMessage
}

func (self *sess) tracePretty(
	args ...interface{},
) {
	if self.trace {
		pretty.Log(args...)
	}
}

func (self *sess) tracePrintln(
	args ...interface{},
) {
	if self.trace {
		log.Println(args...)
	}
}

func New(amqpUrl, amqpExchange, amqpRoutingKey string) *sess {
	sess := sess{
		amqpExchange:   amqpExchange,
		amqpRoutingKey: amqpRoutingKey,
		amqpUrl:        amqpUrl,
	}
	return &sess
}
