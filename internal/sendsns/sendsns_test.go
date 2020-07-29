package sendsns

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
			obj := New("t")
			if obj.trace == true {
				t.Errorf("expected trace=false, got %s", strconv.FormatBool(obj.trace))
			}
			if obj.snsTopicArn != "t" {
				t.Errorf("expected snsTopicArn=t, got %s", obj.snsTopicArn)
			}
			if obj.svc != nil {
				t.Errorf("expected svc=nil, got %#v", obj.svc)
			}
		},
	)
	t.Run("implements MessageSender",
		func(t *testing.T) {
			var _ logevent.MessageSender = New("t")
		},
	)
}

func TestRepeatedOpenAndClose(t *testing.T) {
	obj := New("t")

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
	obj := New("t")
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

func Test_buildSnsMessage_empty_LogEvent(t *testing.T) {
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			// empty!
		},
		Content: logevent.MessageContent{
			// empty!
		},
	}
	obj := New("t")
	m := obj.buildSnsMessage(logEvent)
	t.Run("snsMessage.MessageAttributes",
		func(t *testing.T) {
			if m.MessageAttributes["customer_code"] != nil {
				val := m.MessageAttributes["customer_code"].StringValue
				t.Errorf("incorrect customer_code attribute, expected \"\" got %#v", *val)
			}
			if m.MessageAttributes["host"] != nil {
				val := m.MessageAttributes["host"].StringValue
				t.Errorf("incorrect host attribute, expected \"\" got %#v", *val)
			}
			if m.MessageAttributes["source"] != nil {
				val := m.MessageAttributes["source"].StringValue
				t.Errorf("incorrect source attribute, expected \"\" got %#v", *val)
			}
			if m.MessageAttributes["source_environment"] != nil {
				val := m.MessageAttributes["source_environment"].StringValue
				t.Errorf("incorrect source_environment attribute, expected \"\" got %#v", *val)
			}
			if m.MessageAttributes["sourcetype"] != nil {
				val := m.MessageAttributes["sourcetype"].StringValue
				t.Errorf("incorrect sourcetype attribute, expected \"\" got %#v", *val)
			}
			if m.MessageAttributes["type"] != nil {
				val := m.MessageAttributes["type"].StringValue
				t.Errorf("incorrect type attribute, expected \"\" got %#v", *val)
			}
		},
	)
	t.Run("snsMessage.Message",
		func(t *testing.T) {
			var innerMessage map[string]interface{}
			err := json.Unmarshal([]byte(m.Message), &innerMessage)
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
				//if ! time.IsZero(innerMessage["time"]) {
				//if ! innerMessage["time"].(time).IsZero() {
				t.Errorf("incorrect time, expected %#v got %#v", timeString, innerMessage["time"].(string))
				//t.Errorf("incorrect time, expected zero got %#v", innerMessage["time"])
			}
		},
	)
}

func Test_buildSnsMessage_simple_LogEvent(t *testing.T) {
	now := time.Now()
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			CustomerCode:      "c1",
			Host:              "h1",
			Source:            "s1",
			SourceEnvironment: "se",
			Sourcetype:        "st1",
			Type:              "t1",
		},
		Content: logevent.MessageContent{
			Host:       "h1",
			Index:      "idx1",
			Source:     "s1",
			Sourcetype: "st1",
			Time:       now,
			//Fields: now,
			//Event: now,
		},
	}
	obj := New("t")
	m := obj.buildSnsMessage(logEvent)
	t.Run("snsMessage.MessageAttributes",
		func(t *testing.T) {
			gotCustomerCode := m.MessageAttributes["customer_code"].StringValue
			if *gotCustomerCode != "c1" {
				t.Errorf("incorrect customer_code attribute, expected %#v got %#v", "c1", *gotCustomerCode)
			}
			gotHost := m.MessageAttributes["host"].StringValue
			if *gotHost != "h1" {
				t.Errorf("incorrect host attribute, expected %#v got %#v", "h1", *gotHost)
			}
			gotSource := m.MessageAttributes["source"].StringValue
			if *gotSource != "s1" {
				t.Errorf("incorrect source attribute, expected %#v got %#v", "s1", *gotSource)
			}
			gotSourceEnvironment := m.MessageAttributes["source_environment"].StringValue
			if *gotSourceEnvironment != "se" {
				t.Errorf("incorrect source_environment attribute, expected %#v got %#v", "st1", *gotSourceEnvironment)
			}
			gotSourcetype := m.MessageAttributes["sourcetype"].StringValue
			if *gotSourcetype != "st1" {
				t.Errorf("incorrect sourcetype attribute, expected %#v got %#v", "st1", *gotSourcetype)
			}
			gotType := m.MessageAttributes["type"].StringValue
			if *gotType != "t1" {
				t.Errorf("incorrect type attribute, expected %#v got %#v", "st1", *gotType)
			}
		},
	)
	t.Run("snsMessage.Message",
		func(t *testing.T) {
			var innerMessage map[string]interface{}
			err := json.Unmarshal([]byte(m.Message), &innerMessage)
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
		},
	)
}
