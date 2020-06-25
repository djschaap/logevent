package sendhec

import (
	"github.com/djschaap/logevent"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("with no args",
		func(t *testing.T) {
			obj := New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
			if obj.initialized == true {
				t.Errorf("expected initialized=false, got %s", strconv.FormatBool(obj.initialized))
			}
			if obj.trace == true {
				t.Errorf("expected trace=false, got %s", strconv.FormatBool(obj.trace))
			}

		},
	)
	t.Run("implements MessageSender",
		func(t *testing.T) {
			var _ logevent.MessageSender = New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
		},
	)
}

func TestSetHecInsecure(t *testing.T) {
	obj := New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
	if obj.hecInsecure != false {
		t.Errorf("expected initial hecInsecure=false, got %s",
			strconv.FormatBool(obj.hecInsecure))
	}
	obj.SetHecInsecure(true)
	if obj.hecInsecure != true {
		t.Errorf("expected post-change hecInsecure=true, got %s",
			strconv.FormatBool(obj.hecInsecure))
	}
}

func TestSetTrace(t *testing.T) {
	obj := New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
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

func Test_formatLogEvent_empty_LogEvent(t *testing.T) {
	logEvent := logevent.LogEvent{
		Attributes: logevent.Attributes{
			Host: "h1-ignore", // Attributes are IGNORED by splunkhec
		},
		Content: logevent.MessageContent{
			// empty!
		},
	}
	obj := New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
	obj.SetTrace(true)
	hecEvent := obj.formatLogEvent(logEvent)
	t.Run("hec.Event",
		func(t *testing.T) {
			if hecEvent.Event != nil {
				t.Errorf("incorrect Event, expected nil got %#v", hecEvent.Event)
			}
			if hecEvent.Fields != nil {
				t.Errorf("incorrect Fields attribute, expected nil got %#v", hecEvent.Fields)
			}
			if hecEvent.Host != nil {
				t.Errorf("incorrect Host attribute, expected nil got %#v", *hecEvent.Host)
			}
			if hecEvent.Index != nil {
				t.Errorf("incorrect index, expected nil got %#v", *hecEvent.Index)
			}
			if hecEvent.Source != nil {
				t.Errorf("incorect Source attribute, expected nil got %#v", *hecEvent.Source)
			}
			if hecEvent.SourceType != nil {
				t.Errorf("incorrect SourceType attribute, expected nil got %#v", *hecEvent.SourceType)
			}
			if hecEvent.Time != nil {
				t.Errorf("incorrect Time attribute, expected nil got %#v", *hecEvent.Time)
			}
		},
	)
}

func Test_formatLogEvent_simple_LogEvent(t *testing.T) {
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
			Event: "my event",
			Fields: map[string]interface{}{
				"f": "v1",
			},
			Host:       "h1",
			Index:      "idx1",
			Source:     "s1",
			Sourcetype: "st1",
			Time:       now,
		},
	}
	obj := New("https://localhost:8088", "00000000-0000-0000-0000-000000000000")
	hecEvent := obj.formatLogEvent(logEvent)
	t.Run("hec.Event",
		func(t *testing.T) {
			if hecEvent.Event.(string) != "my event" {
				t.Errorf("incorrect Event, expected \"my event\" got %#v", hecEvent.Event)
			}
			if hecEvent.Fields["f"] != "v1" {
				t.Errorf("incorrect Fields[\"f\"] attribute, expected v1 got %#v", hecEvent.Fields["f"])
			}
			if *hecEvent.Host != "h1" {
				t.Errorf("incorrect Host attribute, expected h1 got %#v", *hecEvent.Host)
			}
			if *hecEvent.Index != "idx1" {
				t.Errorf("incorrect index, expected idx1 got %#v", *hecEvent.Index)
			}
			if *hecEvent.Source != "s1" {
				t.Errorf("incorect Source attribute, expected s1 got %#v", *hecEvent.Source)
			}
			if *hecEvent.SourceType != "st1" {
				t.Errorf("incorrect SourceType attribute, expected st1 got %#v", *hecEvent.SourceType)
			}

			// beware: comparison/test ignores fractional seconds
			re := regexp.MustCompile(`(\d+)`)
			nowString := strconv.FormatInt(int64(now.Unix()), 10)
			gotString := re.FindStringSubmatch(*hecEvent.Time)
			if gotString[1] != nowString {
				t.Errorf("incorrect Time attribute, expected %v got %#v", nowString, gotString)
			}
		},
	)
}
