package main

import (
	"flag"
	"fmt"
	"github.com/djschaap/logevent"
	"github.com/djschaap/logevent/flagarray"
	"github.com/djschaap/logevent/sendamqp"
	"github.com/djschaap/logevent/senddump"
	"github.com/djschaap/logevent/sendhec"
	"github.com/djschaap/logevent/sendsns"
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	buildDt string
	commit  string
	version string
)

func getenvBool(k string) bool {
	v := os.Getenv(k)
	if len(v) > 0 {
		return true
	}
	return false
}

func printVersion() {
	fmt.Println("logevent send  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}

func main() {
	printVersion()

	customerCode := flag.String("customer", "", "set customer code attribute")
	epochAttr := flag.Int64("epoch", 0, "time_t/epoch, as 64-bit int")
	var fieldArgs flagarray.StringArray
	flag.Var(&fieldArgs, "field", "field value, as fieldName=value, may be repeated")
	hostAttr := flag.String("host", "", "set host attribute")
	indexAttr := flag.String("index", "", "set index attribute")
	sourceAttr := flag.String("source", "", "source attribute")
	sourceEnvironmentAttr := flag.String("sourceenvironment", "", "sourceenvironment attribute")
	sourcetypeAttr := flag.String("sourcetype", "", "sourcetype attribute")
	timeAttr := flag.String("time", "", "time, as ISO 8601/RFC 3339")
	printVersion := flag.Bool("v", false, "print version and exit")
	flag.Parse()
	if *printVersion {
		os.Exit(0)
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	senderPackage := os.Getenv("PACKAGE")
	var traceOutput bool
	if len(os.Getenv("TRACE")) > 0 {
		fmt.Println("*** TRACE is enabled ***")
		traceOutput = true
	}

	var sender logevent.MessageSender
	if senderPackage == "sendamqp" {
		amqpURL := os.Getenv("AMQP_URL")
		amqpExchange := os.Getenv("AMQP_EXCHANGE")
		amqpRoutingKey := os.Getenv("AMQP_ROUTING_KEY")
		if len(amqpRoutingKey) <= 0 {
			log.Println("WARNING: sendamqp requires AMQP_ROUTING_KEY; continuing anyway")
		}
		amqpSender := sendamqp.New(amqpURL, amqpExchange, amqpRoutingKey)
		sender = amqpSender
	} else if senderPackage == "senddump" || senderPackage == "" {
		sender = senddump.New()
	} else if senderPackage == "sendhec" {
		hecURL := os.Getenv("HEC_URL")
		hecToken := os.Getenv("HEC_TOKEN")
		if len(hecToken) <= 0 {
			log.Fatal("Splunk HEC_TOKEN must be specified")
		}
		hecSender := sendhec.New(hecURL, hecToken)
		if len(os.Getenv("HEC_INSECURE")) > 0 {
			hecSender.SetHecInsecure(true)
		}
		sender = hecSender
	} else if senderPackage == "sendsns" {
		// github.com/aws/aws-sdk-go/aws reads env vars itself
		//aws_access_key_id := os.Getenv("AWS_ACCESS_KEY_ID")
		//aws_region := os.Getenv("AWS_REGION")
		//aws_secret_access_key := os.Getenv("AWS_SECRET_ACCESS_KEY")
		topicString := os.Getenv("TOPIC")
		hasQueue, _ := regexp.MatchString(`^arn:`, topicString)
		if !hasQueue {
			log.Println("WARNING: sendsns requires TOPIC; continuing anyway")
		}
		sender = sendsns.New(topicString)
	} else {
		log.Fatal("package ", senderPackage, " is not valid")
	}

	if traceOutput {
		sender.SetTrace(true)
	}
	messageContent := flag.Arg(0)
	logEvent := logevent.LogEvent{
		Content: logevent.MessageContent{
			Event: messageContent,
		},
	}

	if *customerCode != "" {
		logEvent.Attributes.CustomerCode = *customerCode
	}
	if *hostAttr != "" {
		logEvent.Attributes.Host = *hostAttr
		logEvent.Content.Host = *hostAttr
	}
	if *indexAttr != "" {
		logEvent.Content.Index = *indexAttr
	}
	if *sourceAttr != "" {
		logEvent.Attributes.Source = *sourceAttr
		logEvent.Content.Source = *sourceAttr
	}
	if *sourceEnvironmentAttr != "" {
		logEvent.Attributes.SourceEnvironment = *sourceEnvironmentAttr
	}
	if *sourcetypeAttr != "" {
		logEvent.Attributes.Sourcetype = *sourcetypeAttr
		logEvent.Content.Sourcetype = *sourcetypeAttr
	}
	if *epochAttr > 0 {
		t := time.Unix(*epochAttr, 0)
		logEvent.Content.Time = t
	} else if *timeAttr != "" {
		t, err := time.Parse(time.RFC3339, *timeAttr)
		if err != nil {
			log.Fatal(err)
		}
		logEvent.Content.Time = t
	}

	logEvent.Content.Fields = make(map[string]interface{})
	for _, rawPair := range fieldArgs {
		re := regexp.MustCompile(`(\S+?)=(.+)`)
		kv := re.FindStringSubmatch(rawPair)
		if len(kv) < 2 {
			log.Fatal("unable to parse field/value: ", rawPair)
		}
		//log.Printf("field: k=%s v=%s\n", kv[1], kv[2]) // DEBUG
		logEvent.Content.Fields[kv[1]] = kv[2]
	}

	err = sender.OpenSvc()
	if err != nil {
		log.Fatal("Error from OpenSvc:", err)
	}
	defer sender.CloseSvc()
	err = sender.SendMessage(logEvent)
	if err != nil {
		log.Fatal("Error from SendMessage:", err)
	}
}
