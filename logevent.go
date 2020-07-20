package logevent

import (
	"time"
)

type MessageSender interface {
	CloseSvc() error
	OpenSvc() error
	SendMessage(LogEvent) error
	SetTrace(bool)
}

type Attributes struct {
	CustomerCode      string `json:"customer_code,omitempty"`
	Host              string `json:"host,omitempty"`
	Source            string `json:"source,omitempty"`
	SourceEnvironment string `json:"source_environment,omitempty"`
	Sourcetype        string `json:"sourcetype,omitempty"`
}

type LogEvent struct {
	Attributes Attributes
	Content    MessageContent
}

type MessageContent struct {
	Host       string                 `json:"host,omitempty"`
	Index      string                 `json:"index,omitempty"`
	Source     string                 `json:"source,omitempty"`
	Sourcetype string                 `json:"sourcetype,omitempty"`
	Time       time.Time              `json:"time,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Event      interface{}            `json:"event,omitempty"`
}
