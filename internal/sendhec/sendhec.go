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

// Sess stores sendhec session state.
type Sess struct {
	hecClient   hec.HEC
	hecInsecure bool
	hecToken    string
	hecURL      string
	trace       bool
}

// CloseSvc closes the open session.
// CloseSvc must not be called when no session is open.
func (sender *Sess) CloseSvc() error {
	if sender.hecClient == nil {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	sender.hecClient = nil
	return nil
}

// OpenSvc opens a new session.
// OpenSvc must not be called when a session is already open.
func (sender *Sess) OpenSvc() error {
	if sender.hecClient != nil {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	client := hec.NewCluster(
		[]string{sender.hecURL},
		sender.hecToken,
	)
	if sender.hecInsecure {
		client.SetHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}})
	}
	sender.hecClient = client
	return nil
}

// SendMessage sends a LogEvent to a Splunk HTTP Event Collector.
func (sender *Sess) SendMessage(logEvent logevent.LogEvent) error {
	if sender.hecClient == nil {
		return errors.New("SendMessage() called before OpenSvc()")
	}
	hecEvents := []*hec.Event{
		sender.formatLogEvent(logEvent),
	}
	sender.tracePretty("TRACE_SENDHEC time =",
		logEvent.Content.Time.UTC().Format(time.RFC3339),
		" hecEvents =", hecEvents)
	err := sender.hecClient.WriteBatch(hecEvents)
	return err
}

// SetHecInsecure disables SSL/TLS validation.
// THIS IS INSECURE but may be useful in dev/lab environments.
func (sender *Sess) SetHecInsecure(v bool) {
	sender.hecInsecure = v
}

// SetTrace enables tracing, which dumps all messages to stderr.
func (sender *Sess) SetTrace(v bool) {
	sender.trace = v
}

func (sender *Sess) formatLogEvent(logEvent logevent.LogEvent) *hec.Event {
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
// It requires a Splunk HEC URL and HEC token (typically, GUID).
func New(hecURL, hecToken string) *Sess {
	sess := Sess{
		hecToken: hecToken,
		hecURL:   hecURL,
	}
	return &sess
}
