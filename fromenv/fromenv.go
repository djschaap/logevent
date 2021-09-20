package fromenv

import (
	"errors"
	"fmt"
	"github.com/djschaap/logevent"
	"github.com/djschaap/logevent/internal/sendamqp"
	"github.com/djschaap/logevent/internal/senddump"
	"github.com/djschaap/logevent/internal/sendhec"
	"github.com/djschaap/logevent/internal/sendsns"
	"log"
	"os"
	"regexp"
)

// os.Getenv mocking concept from alexellis
// https://gist.github.com/alexellis/adc67eb022b7fdca31afc0de6529e5ea
type anyEnv interface {
	Getenv(string) string
	Setenv(string, string)
	Unsetenv(string)
}

type realEnv struct{}

var env anyEnv

func (realEnv) Getenv(k string) string {
	return os.Getenv(k)
}

func (realEnv) Setenv(k string, v string) {
	os.Setenv(k, v)
}

func (realEnv) Unsetenv(k string) {
	os.Unsetenv(k)
}

func GetMessageSenderFromEnv() (logevent.MessageSender, error) {
	initEnv()
	senderPackage := env.Getenv("SENDER_PACKAGE")
	var traceOutput bool
	if len(env.Getenv("SENDER_TRACE")) > 0 {
		fmt.Println("*** SENDER_TRACE is enabled ***")
		traceOutput = true
	}

	var sender logevent.MessageSender
	if senderPackage == "sendamqp" {
		amqpURL := buildAmqpUrl()
		amqpExchange := env.Getenv("AMQP_EXCHANGE")
		amqpRoutingKey := env.Getenv("AMQP_ROUTING_KEY")
		amqpTtl := env.Getenv("AMQP_TTL")
		if len(amqpRoutingKey) <= 0 {
			log.Println("WARNING: sendamqp requires AMQP_ROUTING_KEY; continuing anyway")
		}
		amqpSender := sendamqp.New(amqpURL, amqpExchange, amqpRoutingKey, amqpTtl)
		sender = amqpSender
	} else if senderPackage == "sendhec" {
		hecURL := env.Getenv("HEC_URL")
		hecToken := env.Getenv("HEC_TOKEN")
		if len(hecToken) <= 0 {
			return nil, errors.New("FATAL: sendhec requires HEC_TOKEN")
		}
		hecSender := sendhec.New(hecURL, hecToken)
		if len(env.Getenv("HEC_INSECURE")) > 0 {
			hecSender.SetHecInsecure(true)
		}
		sender = hecSender
	} else if senderPackage == "sendsns" {
		// github.com/aws/aws-sdk-go/aws reads env vars itself
		//aws_access_key_id := env.Getenv("AWS_ACCESS_KEY_ID")
		//aws_region := env.Getenv("AWS_REGION")
		//aws_secret_access_key := env.Getenv("AWS_SECRET_ACCESS_KEY")
		topicString := env.Getenv("AWS_SNS_TOPIC")
		hasQueue, _ := regexp.MatchString(`^arn:`, topicString)
		if !hasQueue {
			log.Println("WARNING: sendsns requires AWS_SNS_TOPIC; continuing anyway")
		}
		sender = sendsns.New(topicString)
	} else if senderPackage == "senddump" || senderPackage == "" {
		sender = senddump.New()
	} else {
		return nil, errors.New("FATAL: SENDER_PACKAGE " + senderPackage + " is not valid")
	}

	if traceOutput {
		sender.SetTrace(true)
	}
	return sender, nil
}

func buildAmqpUrl() string {
	var amqpUrl string
	amqpUrl = env.Getenv("AMQP_URL")
	if len(amqpUrl) > 0 {
		return amqpUrl
	}

	amqpHost := env.Getenv("AMQP_HOST")
	amqpPassword := env.Getenv("AMQP_PASSWORD")
	amqpPort := env.Getenv("AMQP_PORT")
	amqpProtocol := env.Getenv("AMQP_PROTOCOL")
	amqpUsername := env.Getenv("AMQP_USERNAME")
	amqpVhost := env.Getenv("AMQP_VHOST")

	var userPass string
	if len(amqpPassword) > 0 {
		userPass = fmt.Sprintf("%s:%s@", amqpUsername, amqpPassword)
	} else {
		userPass = fmt.Sprintf("%s@", amqpUsername)
	}

	var hostPort string
	if len(amqpPort) > 0 {
		hostPort = fmt.Sprintf("%s:%s", amqpHost, amqpPort)
	} else {
		hostPort = amqpHost
	}
	amqpUrl = fmt.Sprintf("%s://%s%s", amqpProtocol, userPass, hostPort)
	if len(amqpVhost) > 0 {
		amqpUrl = amqpUrl + "/" + amqpVhost
	}

	return amqpUrl
}

func getenvBool(k string) bool {
	initEnv()
	v := env.Getenv(k)
	if len(v) > 0 {
		return true
	}
	return false
}

func initEnv() {
	if env == nil {
		env = realEnv{}
	}
}
