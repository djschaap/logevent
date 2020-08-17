package main

import (
	"flag"
	"fmt"
	"github.com/djschaap/logevent"
	"github.com/djschaap/logevent/flagarray"
	"github.com/djschaap/logevent/fromenv"
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

func printVersion() {
	fmt.Println("logevent send  Version:",
		version, " Commit:", commit,
		" Built at:", buildDt)
}

func main() {
	printVersion()

	eventCount := flag.Int("count", 1, "send N events")
	repeatDelay := flag.Int("delay", 1, "delay N seconds between events")
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
		log.Fatal("Error loading .env file:", err)
	}
	sender, err := fromenv.GetMessageSenderFromEnv()
	if err != nil {
		log.Fatal("Error initializing output:", err)
	}
	err = sender.OpenSvc()
	if err != nil {
		log.Fatal("Error from OpenSvc:", err)
	}
	defer sender.CloseSvc()

	for i := 0; i < *eventCount; i++ {
		if i > 0 {
			// delay before any additional events
			//log.Println("sleeping...") // DEBUG
			time.Sleep(time.Duration(*repeatDelay) * time.Second)
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

		err = sender.SendMessage(logEvent)
		if err != nil {
			log.Fatal("Error from SendMessage:", err)
		}
	}
}
