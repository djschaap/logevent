package logevent

import (
	"time"
)

// MessageSender is the interface that manages a connection to an external LogEvent destination.
type MessageSender interface {
	CloseSvc() error
	OpenSvc() error
	SendMessage(LogEvent) error
	SetTrace(bool)
}

// Attributes contains the properties of a single LogEvent that may be passed as headers to an intermediate system (such as a message queue).
type Attributes struct {
	CustomerCode      string `json:"customer_code,omitempty"`
	Host              string `json:"host,omitempty"`
	Source            string `json:"source,omitempty"`
	SourceEnvironment string `json:"source_environment,omitempty"`
	Sourcetype        string `json:"sourcetype,omitempty"`
	Type              string
}

// LogEvent contains a single log message, made up of attributes/headers and message content.
type LogEvent struct {
	Attributes Attributes
	Content    MessageContent
}

// MessageContent is the actual message to be sent to the external destination.
type MessageContent struct {
	Host       string                 `json:"host,omitempty"`
	Index      string                 `json:"index,omitempty"`
	Source     string                 `json:"source,omitempty"`
	Sourcetype string                 `json:"sourcetype,omitempty"`
	Time       time.Time              `json:"time,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Event      interface{}            `json:"event,omitempty"`
}
