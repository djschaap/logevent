package senddump

import (
	"errors"
	"github.com/djschaap/logevent"
	"github.com/kr/pretty"
	"log"
	"time"
)

// Sess stores senddump session state.
type Sess struct {
	initialized bool
	trace       bool
}

// CloseSvc closes the open session.
// CloseSvc must not be called when no session is open.
func (sender *Sess) CloseSvc() error {
	if !sender.initialized {
		return errors.New("CloseSvc() called again or before OpenSvc(); that should not be done")
	}
	sender.initialized = false
	return nil
}

// OpenSvc opens a new session.
// OpenSvc must not be called when a session is already open.
func (sender *Sess) OpenSvc() error {
	if sender.initialized {
		return errors.New("OpenSvc() called again; that should not be done")
	}
	sender.initialized = true
	return nil
}

// SendMessage dumps a LogEvent to stderr when tracing is enabled.
// No output is generated unless tracing is enabled.
func (sender *Sess) SendMessage(logEvent logevent.LogEvent) error {
	if !sender.initialized {
		return errors.New("SendMessage() called before OpenSvc()")
	}
	timeString := logEvent.Content.Time.UTC().Format(time.RFC3339)
	logEvent.Content.Time = time.Time{}
	sender.tracePretty("TRACE_SENDER time =", timeString,
		" logEvent =", logEvent)
	return nil
}

// SetTrace enables tracing, which dumps all messages to stderr.
func (sender *Sess) SetTrace(v bool) {
	sender.trace = v
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

// New creates a new senddump object/session.
func New() *Sess {
	sess := Sess{}
	return &sess
}
