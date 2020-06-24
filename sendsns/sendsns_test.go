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

func Test_buildSnsMessage_empty_LogEvent(t *testing.T) {
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			// empty!
		},
		Content: logevent.MessageContent{
			// empty!
		},
	}
	obj := New()
	m := obj.buildSnsMessage(logEvent)
	t.Run("SnsMessage.MessageAttributes",
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
		},
	)
	t.Run("SnsMessage.Message",
		func(t *testing.T) {
			var inner_message map[string]interface{}
			err := json.Unmarshal([]byte(m.Message), &inner_message)
			if err != nil {
				t.Error("json.Unmarshal error:", err)
			}
			if inner_message["host"] != nil {
				t.Errorf("incorrect host, expected nil got %#v", inner_message["host"])
			}
			if inner_message["index"] != nil {
				t.Errorf("incorrect index, expected nil got %#v", inner_message["index"])
			}
			if inner_message["source"] != nil {
				t.Errorf("incorrect source, expected nil got %#v", inner_message["source"])
			}
			if inner_message["sourcetype"] != nil {
				t.Errorf("incorrect sourcetype, expected nil got %#v", inner_message["sourcetype"])
			}
			timeString := "0001-01-01T00:00:00Z"
			if inner_message["time"].(string) != timeString {
				//if ! time.IsZero(inner_message["time"]) {
				//if ! inner_message["time"].(time).IsZero() {
				t.Errorf("incorrect time, expected %#v got %#v", timeString, inner_message["time"].(string))
				//t.Errorf("incorrect time, expected zero got %#v", inner_message["time"])
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
	obj := New()
	m := obj.buildSnsMessage(logEvent)
	t.Run("SnsMessage.MessageAttributes",
		func(t *testing.T) {
			got_customer_code := m.MessageAttributes["customer_code"].StringValue
			if *got_customer_code != "c1" {
				t.Errorf("incorrect customer_code attribute, expected %#v got %#v", "c1", *got_customer_code)
			}
			got_host := m.MessageAttributes["host"].StringValue
			if *got_host != "h1" {
				t.Errorf("incorrect host attribute, expected %#v got %#v", "h1", *got_host)
			}
			got_source := m.MessageAttributes["source"].StringValue
			if *got_source != "s1" {
				t.Errorf("incorrect source attribute, expected %#v got %#v", "s1", *got_source)
			}
			got_source_environment := m.MessageAttributes["source_environment"].StringValue
			if *got_source_environment != "se" {
				t.Errorf("incorrect source_environment attribute, expected %#v got %#v", "st1", *got_source_environment)
			}
			got_sourcetype := m.MessageAttributes["sourcetype"].StringValue
			if *got_sourcetype != "st1" {
				t.Errorf("incorrect sourcetype attribute, expected %#v got %#v", "st1", *got_sourcetype)
			}
		},
	)
	t.Run("SnsMessage.Message",
		func(t *testing.T) {
			var inner_message map[string]interface{}
			err := json.Unmarshal([]byte(m.Message), &inner_message)
			if err != nil {
				t.Error("json.Unmarshal error:", err)
			}
			if inner_message["host"].(string) != "h1" {
				t.Errorf("incorrect host, expected %#v got %#v", "h1", inner_message["host"].(string))
			}
			if inner_message["index"].(string) != "idx1" {
				t.Errorf("incorrect index, expected %#v got %#v", "idx1", inner_message["index"].(string))
			}
			if inner_message["source"].(string) != "s1" {
				t.Errorf("incorrect source, expected %#v got %#v", "s1", inner_message["source"].(string))
			}
			if inner_message["sourcetype"].(string) != "st1" {
				t.Errorf("incorrect sourcetype, expected %#v got %#v", "st1", inner_message["sourcetype"].(string))
			}
			timeString := now.Format(time.RFC3339Nano)
			if inner_message["time"].(string) != timeString {
				t.Errorf("incorrect time, expected %#v got %#v", timeString, inner_message["time"].(string))
			}
		},
	)
}
