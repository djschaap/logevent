package senddump

import (
	"github.com/djschaap/logevent"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("with no args",
		func(t *testing.T) {
			obj := New()
			if obj.trace == true {
				t.Errorf("expected trace=false, got %s", strconv.FormatBool(obj.trace))
			}
		},
	)
	t.Run("implements MessageSender",
		func(t *testing.T) {
			var _ logevent.MessageSender = New()
		},
	)
}

func TestRepeatedOpenAndClose(t *testing.T) {
	obj := New()

	err := obj.OpenSvc()
	if err != nil {
		t.Errorf("OpenSvc() returned unexpected error %v", err)
	}

	err = obj.OpenSvc()
	if err == nil {
		t.Error("expected error from OpenSvc() but got nil")
	}

	err = obj.CloseSvc()
	if err != nil {
		t.Errorf("CloseSvc() returned unexpected error %v", err)
	}

	err = obj.CloseSvc()
	if err == nil {
		t.Error("expected error from CloseSvc() but got nil")
	}
}

func TestSendMessage(t *testing.T) {
	obj := New()
	logEvent := logevent.LogEvent{}

	err := obj.SendMessage(logEvent)
	if err == nil {
		t.Error("expected error from SendMessage() but got nil")
	}

	err = obj.OpenSvc()
	if err != nil {
		t.Errorf("OpenSvc() returned unexpected error %v", err)
	}
	defer obj.CloseSvc()
	obj.SendMessage(logEvent)
}

func TestSetTrace(t *testing.T) {
	obj := New()
	if obj.trace != false {
		t.Errorf("expected initial trace=false, got %s",
			strconv.FormatBool(obj.trace))
	}
	obj.SetTrace(true)
	if obj.trace != true {
		t.Errorf("expected post-change trace=true, got %s",
			strconv.FormatBool(obj.trace))
	}
	obj.tracePretty("test tracePretty output")
	obj.tracePrintln("test tracePrintln output")
}
