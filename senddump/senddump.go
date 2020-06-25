package senddump

import (
	"errors"
	"github.com/djschaap/logevent"
	"github.com/kr/pretty"
	"log"
	"time"
)

// structs

type sess struct {
	initialized bool
	trace       bool
}

// function(s)

func (self *sess) OpenSvc() error {
	if self.initialized {
		return errors.New("OpenSvc() called again; that should not happen")
	}
	self.initialized = true
	return nil
}

func (self *sess) SendMessage(logEvent logevent.LogEvent) error {
	if !self.initialized {
		return errors.New("SendMessage() called before OpenSvc()")
	}
	timeString := logEvent.Content.Time.UTC().Format(time.RFC3339)
	logEvent.Content.Time = time.Time{}
	self.tracePretty("TRACE_SENDER time =", timeString,
		" logEvent =", logEvent)
	return nil
}

func (self *sess) SetTrace(v bool) {
	self.trace = v
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
