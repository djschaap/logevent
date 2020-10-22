// +build amqp

package sendamqp

import (
	"github.com/djschaap/logevent"
	"os"
	"testing"
	"time"
)

func TestOpenSvc(t *testing.T) {
	amqpUrl := os.Getenv("AMQP_URL")
	t.Run("with no args",
		func(t *testing.T) {
			obj := New(amqpUrl, "exch-unsed", "rk-unused", "")
			err := obj.OpenSvc()
			if err != nil {
				t.Errorf("OpenSvc() returned err: %s", err)
			}
			err = obj.OpenSvc() // second call
			if err == nil {
				// expect "OpenSvc() called again; that should not be done"
				t.Error("second OpenSvc() expected err but returned nil")
			}
		},
	)
	t.Run("implements MessageSender",
		func(t *testing.T) {
			var _ logevent.MessageSender = New("u", "e", "t", "")
		},
	)
}

func TestRepeatedOpenAndClose(t *testing.T) {
	obj := New("amqp://localhost", "exch", "rk", "")

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

func TestSendMessage_empty(t *testing.T) {
	amqpUrl := os.Getenv("AMQP_URL")
	obj := New(amqpUrl, "amq.headers", "sendamqp_amqp_test_discard", "")
	logEvent := logevent.LogEvent{}

	err := obj.SendMessage(logEvent)
	if err == nil {
		// expect SendMessage() called before OpenSvc()
		t.Error("expected error from SendMessage() but got nil")
	}

	err = obj.OpenSvc()
	if err != nil {
		t.Errorf("OpenSvc() returned unexpected err: %s", err)
	}
	defer obj.CloseSvc()

	obj.SendMessage(logEvent)
	// FUTURE consume message
}

func TestSendMessage_simple(t *testing.T) {
	amqpUrl := os.Getenv("AMQP_URL")
	obj := New(amqpUrl, "amq.headers", "sendamqp_amqp_test_discard", "")
	err := obj.OpenSvc()
	if err != nil {
		t.Errorf("OpenSvc() returned unexpected err: %s", err)
	}
	defer obj.CloseSvc()

	now := time.Now()
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			CustomerCode:      "c1",
			Host:              "h1",
			Source:            "s1",
			SourceEnvironment: "se",
			Sourcetype:        "st1",
		},
		Content: logevent.MessageContent{
			Host:       "h1",
			Index:      "idx1",
			Source:     "s1",
			Sourcetype: "st1",
			Time:       now,
			Event:      `{"msgsrc":"sendamqp_amqp_test TestSendMessage_simple"}`,
			//Fields: {},
		},
	}
	obj.SendMessage(logEvent)
	// FUTURE consume message
}
