package logevent

import (
	"testing"
	"time"
)

func Test_AttributesStruct(t *testing.T) {
	expect := map[string]string{
		"CustomerCode":      "c",
		"Host":              "h",
		"Source":            "s",
		"SourceEnvironment": "se",
		"Sourcetype":        "st",
	}
	a := Attributes{
		CustomerCode:      "c",
		Host:              "h",
		Source:            "s",
		SourceEnvironment: "se",
		Sourcetype:        "st",
	}
	if a.CustomerCode != expect["CustomerCode"] {
		t.Errorf("expected CustomerCode %#v, got %#v", expect["CustomerCode"], a.CustomerCode)
	}
	if a.Host != expect["Host"] {
		t.Errorf("expected Host %#v, got %#v", expect["Host"], a.Host)
	}
	if a.Source != expect["Source"] {
		t.Errorf("expected Source %#v, got %#v", expect["Source"], a.Source)
	}
	if a.SourceEnvironment != expect["SourceEnvironment"] {
		t.Errorf("expected SourceEnvironment %#v, got %#v", expect["SourceEnvironment"], a.SourceEnvironment)
	}
	if a.Sourcetype != expect["Sourcetype"] {
		t.Errorf("expected Sourcetype %#v, got %#v", expect["Sourcetype"], a.Sourcetype)
	}
}

func Test_LogEventStruct(t *testing.T) {
	_ = LogEvent{
		Attributes: Attributes{},
		Content:    MessageContent{},
	}
}

func Test_MessageContentStruct(t *testing.T) {
	now := time.Now()
	expect := map[string]string{
		"Host":       "h",
		"Index":      "i",
		"Source":     "s",
		"Sourcetype": "st",
	}
	a := MessageContent{
		Host:       "h",
		Index:      "i",
		Source:     "s",
		Sourcetype: "st",
		Time:       now,
	}
	if a.Host != expect["Host"] {
		t.Errorf("expected Host %#v, got %#v", expect["Host"], a.Host)
	}
	if a.Index != expect["Index"] {
		t.Errorf("expected Index %#v, got %#v", expect["Index"], a.Index)
	}
	if a.Source != expect["Source"] {
		t.Errorf("expected Source %#v, got %#v", expect["Source"], a.Source)
	}
	if a.Sourcetype != expect["Sourcetype"] {
		t.Errorf("expected Sourcetype %#v, got %#v", expect["Sourcetype"], a.Sourcetype)
	}
	if a.Time != now {
		t.Errorf("expected Time %#v, got %#v", now, a.Time)
	}
}
