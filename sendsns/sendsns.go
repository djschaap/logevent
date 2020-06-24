package sendsns

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/djschaap/logevent"
	"github.com/kr/pretty"
	"log"
)

// structs

type SnsMessage struct {
	Message           string
	MessageAttributes map[string]*sns.MessageAttributeValue
}

type sess struct {
	svc   *sns.SNS
	trace bool
}

// function(s)

func (self *sess) OpenSvc() error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	self.svc = sns.New(sess)
	return nil
}

func (self *sess) SendMessage(topicArn string, logEvent logevent.LogEvent) error {
	self.tracePretty("TRACE_SNS logEvent =", logEvent)
	snsMessage := self.buildSnsMessage(logEvent)
	self.tracePretty("TRACE_SNS MessageAttributes =", snsMessage.MessageAttributes)
	self.tracePretty("TRACE_SNS Message =", snsMessage.Message)

	result, err := self.svc.Publish(&sns.PublishInput{
		MessageAttributes: snsMessage.MessageAttributes,
		Message:           aws.String(snsMessage.Message),
		TopicArn:          &topicArn,
	})

	if err != nil {
		return err
	}

	self.tracePrintln("TRACE_SNS Success", *result.MessageId)
	return nil
}

func (self *sess) SetTrace(v bool) {
	self.trace = v
}

func (self *sess) buildSnsMessage(logEvent logevent.LogEvent) SnsMessage {
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
	messageJson_bytes, _ := json.Marshal(logEvent.Content)
	snsMessage := SnsMessage{
		Message:           string(messageJson_bytes),
		MessageAttributes: messageAttributes,
	}
	return snsMessage
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

func New() *sess {
	sess := sess{}
	return &sess
}