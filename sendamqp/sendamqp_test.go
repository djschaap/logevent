package sendamqp

import (
	"encoding/json"
	"github.com/djschaap/logevent"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("with no args",
		func(t *testing.T) {
			obj := New("amqp://localhost", "exch", "rk")
			if obj.trace == true {
				t.Errorf("expected trace=false, got %s", strconv.FormatBool(obj.trace))
			}
			if obj.amqpExchange != "exch" {
				t.Errorf("expected amqpExchange=t, got %s", obj.amqpExchange)
			}
			if obj.amqpRoutingKey != "rk" {
				t.Errorf("expected amqpRoutingKey=t, got %s", obj.amqpRoutingKey)
			}
		},
	)
	t.Run("implements MessageSender",
		func(t *testing.T) {
			var _ logevent.MessageSender = New("u", "e", "rk")
		},
	)
}

func TestRepeatedOpenAndClose(t *testing.T) {
	obj := New("amqp://localhost", "exch", "rk")

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

func TestSetTrace(t *testing.T) {
	obj := New("u", "e", "rk")
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

func Test_buildAmqpMessage_empty_LogEvent(t *testing.T) {
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			// empty!
		},
		Content: logevent.MessageContent{
			// empty!
		},
	}
	obj := New("u", "e", "rk")
	m := obj.buildAmqpMessage(logEvent)
	t.Run("Headers",
		func(t *testing.T) {
			if m.Headers["customer_code"] != nil {
				val := m.Headers["customer_code"]
				t.Errorf("incorrect customer_code attribute, expected \"\" got %#v", val.(string))
			}
			if m.Headers["host"] != nil {
				val := m.Headers["host"]
				t.Errorf("incorrect host attribute, expected \"\" got %#v", val.(string))
			}
			if m.Headers["source"] != nil {
				val := m.Headers["source"]
				t.Errorf("incorrect source attribute, expected \"\" got %#v", val.(string))
			}
			if m.Headers["source_environment"] != nil {
				val := m.Headers["source_environment"]
				t.Errorf("incorrect source_environment attribute, expected \"\" got %#v", val.(string))
			}
			if m.Headers["sourcetype"] != nil {
				val := m.Headers["sourcetype"]
				t.Errorf("incorrect sourcetype attribute, expected \"\" got %#v", val.(string))
			}
		},
	)
	t.Run("Body",
		func(t *testing.T) {
			var innerMessage map[string]interface{}
			err := json.Unmarshal([]byte(m.Body), &innerMessage)
			if err != nil {
				t.Error("json.Unmarshal error:", err)
			}
			if innerMessage["host"] != nil {
				t.Errorf("incorrect host, expected nil got %#v", innerMessage["host"])
			}
			if innerMessage["index"] != nil {
				t.Errorf("incorrect index, expected nil got %#v", innerMessage["index"])
			}
			if innerMessage["source"] != nil {
				t.Errorf("incorrect source, expected nil got %#v", innerMessage["source"])
			}
			if innerMessage["sourcetype"] != nil {
				t.Errorf("incorrect sourcetype, expected nil got %#v", innerMessage["sourcetype"])
			}
			timeString := "0001-01-01T00:00:00Z"
			if innerMessage["time"].(string) != timeString {
				t.Errorf("incorrect time, expected %#v got %#v", timeString, innerMessage["time"].(string))
			}
		},
	)
}

func Test_buildAmqpMessage_simple_LogEvent(t *testing.T) {
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
			Event:      `{"f1":"v1","f2":"v2"}`,
			//Fields: {},
		},
	}
	obj := New("u", "e", "rk")
	m := obj.buildAmqpMessage(logEvent)
	t.Run("Headers",
		func(t *testing.T) {
			gotCustomerCode := m.Headers["customer_code"]
			if gotCustomerCode != "c1" {
				t.Errorf("incorrect customer_code attribute, expected %#v got %#v", "c1", gotCustomerCode)
			}
			gotHost := m.Headers["host"]
			if gotHost != "h1" {
				t.Errorf("incorrect host attribute, expected %#v got %#v", "h1", gotHost)
			}
			gotSource := m.Headers["source"]
			if gotSource != "s1" {
				t.Errorf("incorrect source attribute, expected %#v got %#v", "s1", gotSource)
			}
			gotSourceEnvironment := m.Headers["source_environment"]
			if gotSourceEnvironment != "se" {
				t.Errorf("incorrect source_environment attribute, expected %#v got %#v", "st1", gotSourceEnvironment)
			}
			gotSourcetype := m.Headers["sourcetype"]
			if gotSourcetype != "st1" {
				t.Errorf("incorrect sourcetype attribute, expected %#v got %#v", "st1", gotSourcetype)
			}
		},
	)
	t.Run("Body",
		func(t *testing.T) {
			var innerMessage map[string]interface{}
			err := json.Unmarshal([]byte(m.Body), &innerMessage)
			if err != nil {
				t.Error("json.Unmarshal error:", err)
			}
			if innerMessage["host"].(string) != "h1" {
				t.Errorf("incorrect host, expected %#v got %#v", "h1", innerMessage["host"].(string))
			}
			if innerMessage["index"].(string) != "idx1" {
				t.Errorf("incorrect index, expected %#v got %#v", "idx1", innerMessage["index"].(string))
			}
			if innerMessage["source"].(string) != "s1" {
				t.Errorf("incorrect source, expected %#v got %#v", "s1", innerMessage["source"].(string))
			}
			if innerMessage["sourcetype"].(string) != "st1" {
				t.Errorf("incorrect sourcetype, expected %#v got %#v", "st1", innerMessage["sourcetype"].(string))
			}
			timeString := now.Format(time.RFC3339Nano)
			if innerMessage["time"].(string) != timeString {
				t.Errorf("incorrect time, expected %#v got %#v", timeString, innerMessage["time"].(string))
			}
			if innerMessage["event"].(string) != `{"f1":"v1","f2":"v2"}` {
				t.Errorf("incorrect event string, got %s", innerMessage["event"].(string))
			}
		},
	)
}
