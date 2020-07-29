package sendsns

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/djschaap/logevent"
	"github.com/kr/pretty"
	"log"
)

type snsMessage struct {
	Message           string
	MessageAttributes map[string]*sns.MessageAttributeValue
}

// Sess stores sendsns session state.
type Sess struct {
	snsTopicArn string
	svc         *sns.SNS
	trace       bool
}

// CloseSvc closes the open session.
// CloseSvc must not be called when no session is open.
func (sender *Sess) CloseSvc() error {
	if sender.svc == nil {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	sender.svc = nil
	return nil
}

// OpenSvc opens a new session.
// OpenSvc must not be called when a session is already open.
func (sender *Sess) OpenSvc() error {
	if sender.svc != nil {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sender.svc = sns.New(sess)
	return nil
}

// SendMessage sends a LogEvent to Amazon Simple Notification Service.
func (sender *Sess) SendMessage(logEvent logevent.LogEvent) error {
	if sender.svc == nil {
		return errors.New("SendMessage() called before OpenSvc()")
	}
	snsMessage := sender.buildSnsMessage(logEvent)
	sender.tracePretty("TRACE_SNS MessageAttributes =", snsMessage.MessageAttributes)
	sender.tracePretty("TRACE_SNS Message =", snsMessage.Message)

	result, err := sender.svc.Publish(&sns.PublishInput{
		MessageAttributes: snsMessage.MessageAttributes,
		Message:           aws.String(snsMessage.Message),
		TopicArn:          &sender.snsTopicArn,
	})

	if err != nil {
		return err
	}

	sender.tracePrintln("TRACE_SNS Success", *result.MessageId)
	return nil
}

// SetTrace enables tracing, which dumps all messages to stderr.
func (sender *Sess) SetTrace(v bool) {
	sender.trace = v
}

func (sender *Sess) buildSnsMessage(logEvent logevent.LogEvent) snsMessage {
	attr := logEvent.Attributes
	messageAttributes := make(map[string]*sns.MessageAttributeValue)
	if attr.CustomerCode != "" {
		messageAttributes["customer_code"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attr.CustomerCode),
		}
	}
	if attr.Host != "" {
		messageAttributes["host"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attr.Host),
		}
	}
	if attr.Source != "" {
		messageAttributes["source"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attr.Source),
		}
	}
	if attr.SourceEnvironment != "" {
		messageAttributes["source_environment"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attr.SourceEnvironment),
		}
	}
	if attr.Sourcetype != "" {
		messageAttributes["sourcetype"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(attr.Sourcetype),
		}
	}
	messageJSONBytes, _ := json.Marshal(logEvent.Content)
	snsMsg := snsMessage{
		Message:           string(messageJSONBytes),
		MessageAttributes: messageAttributes,
	}
	return snsMsg
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

// New creates a new sendsns object/session.
// It requires an SNS topic ARN.
func New(snsTopicArn string) *Sess {
	sess := Sess{
		snsTopicArn: snsTopicArn,
	}
	return &sess
}
