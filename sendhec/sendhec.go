package sendhec

import (
	"crypto/tls"
	"errors"
	"github.com/djschaap/logevent"
	"github.com/fuyufjh/splunk-hec-go" // hec
	"github.com/kr/pretty"
	"log"
	"net/http"
	"time"
)

type sess struct {
	hecClient   hec.HEC
	hecInsecure bool
	hecToken    string
	hecUrl      string
	trace       bool
}

func (self *sess) CloseSvc() error {
	if self.hecClient == nil {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	self.hecClient = nil
	return nil
}

func (self *sess) OpenSvc() error {
	if self.hecClient != nil {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	client := hec.NewCluster(
		[]string{self.hecUrl},
		self.hecToken,
	)
	if self.hecInsecure {
		client.SetHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}})
	}
	self.hecClient = client
	return nil
}

func (self *sess) SendMessage(logEvent logevent.LogEvent) error {
	if self.hecClient == nil {
		return errors.New("SendMessage() called before OpenSvc()")
	}
	hecEvents := []*hec.Event{
		self.formatLogEvent(logEvent),
	}
	self.tracePretty("TRACE_SENDHEC time =",
		logEvent.Content.Time.UTC().Format(time.RFC3339),
		" hecEvents =", hecEvents)
	err := self.hecClient.WriteBatch(hecEvents)
	return err
}

func (self *sess) SetHecInsecure(v bool) {
	self.hecInsecure = v
}

func (self *sess) SetTrace(v bool) {
	self.trace = v
}

func (self *sess) formatLogEvent(logEvent logevent.LogEvent) *hec.Event {
	var hecEvent *hec.Event
	hecEvent = hec.NewEvent(logEvent.Content.Event)
	if logEvent.Content.Host != "" {
		hecEvent.SetHost(logEvent.Content.Host)
	}
	if logEvent.Content.Index != "" {
		hecEvent.SetIndex(logEvent.Content.Index)
	}
	if logEvent.Content.Source != "" {
		hecEvent.SetSource(logEvent.Content.Source)
	}
	if logEvent.Content.Sourcetype != "" {
		hecEvent.SetSourceType(logEvent.Content.Sourcetype)
	}

	// beware: any time before Unix epoch (1970-01-01 00:00:00 UTC) is
	//   negative time_t and will result in an error from Splunk HEC
	if !logEvent.Content.Time.IsZero() {
		hecEvent.SetTime(logEvent.Content.Time)
	}

	hecEvent.SetFields(logEvent.Content.Fields)

	return hecEvent
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

func New(hecUrl, hecToken string) *sess {
	sess := sess{
		hecToken: hecToken,
		hecUrl:   hecUrl,
	}
	return &sess
}
