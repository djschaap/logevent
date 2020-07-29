package fromenv

import (
	"fmt"
	//"github.com/djschaap/logevent"
	"testing"
)

// os.Getenv mocking concept from alexellis
// https://gist.github.com/alexellis/adc67eb022b7fdca31afc0de6529e5ea
type fakeEnv struct {
	values map[string]string
}

func (env fakeEnv) Getenv(k string) string {
	return env.values[k]
}

func (env fakeEnv) Setenv(k string, v string) {
	env.values[k] = v
}

func (env fakeEnv) Unsetenv(k string) {
	delete(env.values, k)
}

func newFakeEnv() fakeEnv {
	e := fakeEnv{}
	e.values = make(map[string]string)
	return e
}

func TestGetenvBool(t *testing.T) {
	env = newFakeEnv()

	env.Setenv("LOGEVENT_UNIT_TEST_FALSE", "false")
	env.Setenv("LOGEVENT_UNIT_TEST_N", "n")
	env.Setenv("LOGEVENT_UNIT_TEST_TRUE", "true")
	env.Setenv("LOGEVENT_UNIT_TEST_WHITESPACE", " ")
	env.Setenv("LOGEVENT_UNIT_TEST_X", "x")

	t.Run("non-existant env var",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_DOES_NOT_EXIST") != false {
				t.Error("expected false")
			}
		},
	)

	t.Run("with false (BEWARE)",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_FALSE") != true {
				t.Error("expected true")
			}
		},
	)

	t.Run("with n (BEWARE)",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_N") != true {
				t.Error("expected true")
			}
		},
	)

	t.Run("with true",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_TRUE") != true {
				t.Error("expected true")
			}
		},
	)

	t.Run("with whitespace (BEWARE)",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_WHITESPACE") != true {
				t.Error("expected true")
			}
		},
	)

	t.Run("with x",
		func(t *testing.T) {
			if getenvBool("LOGEVENT_UNIT_TEST_X") != true {
				t.Error("expected true")
			}
		},
	)
}

func TestGetMessageSenderFromEnv(t *testing.T) {
	t.Run("no SENDER_PACKAGE",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("SENDER_TRACE", "x")
			s, err := GetMessageSenderFromEnv()
			if err != nil {
				t.Errorf("expected success but got error: %s", err)
			}
			expectedType := "*senddump.Sess"
			senderType := fmt.Sprintf("%T", s)
			if senderType != expectedType {
				t.Errorf("expected %s, got %s", expectedType, senderType)
			}
			// TODO verify s.trace is true
		},
	)

	t.Run("invalid SENDER_PACKAGE",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("SENDER_PACKAGE", "i do not exist")
			expectedError := "FATAL: SENDER_PACKAGE i do not exist is not valid"
			s, err := GetMessageSenderFromEnv()
			errStr := fmt.Sprintf("%s", err)
			if errStr != expectedError {
				t.Errorf("expected: %s but got: %s", expectedError, err)
			}
			if s != nil {
				t.Errorf("expected no MessageSender but got: %#v", s)
			}
		},
	)

	t.Run("minimal sendamqp",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("AMQP_ROUTING_KEY", "x")
			env.Setenv("SENDER_PACKAGE", "sendamqp")
			s, err := GetMessageSenderFromEnv()
			if err != nil {
				t.Errorf("expected success but got error: %s", err)
			}
			expectedType := "*sendamqp.Sess"
			senderType := fmt.Sprintf("%T", s)
			if senderType != expectedType {
				t.Errorf("expected %s, got %s", expectedType, senderType)
			}
			// TODO validate AMQP properties
		},
	)

	t.Run("sendhec w/out token",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("SENDER_PACKAGE", "sendhec")
			expectedError := "FATAL: sendhec requires HEC_TOKEN"
			s, err := GetMessageSenderFromEnv()
			errStr := fmt.Sprintf("%s", err)
			if errStr != expectedError {
				t.Errorf("expected: %s but got: %s", expectedError, err)
			}
			if s != nil {
				t.Errorf("expected no MessageSender but got: %#v", s)
			}
		},
	)

	t.Run("minimal sendhec",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("HEC_TOKEN", "x")
			env.Setenv("SENDER_PACKAGE", "sendhec")
			s, err := GetMessageSenderFromEnv()
			if err != nil {
				t.Errorf("expected success but got error: %s", err)
			}
			expectedType := "*sendhec.Sess"
			senderType := fmt.Sprintf("%T", s)
			if senderType != expectedType {
				t.Errorf("expected %s, got %s", expectedType, senderType)
			}
			// TODO validate HEC properties
		},
	)

	t.Run("minimal sendsns",
		func(t *testing.T) {
			env = newFakeEnv()
			env.Setenv("AWS_SNS_TOPIC", "x")
			env.Setenv("SENDER_PACKAGE", "sendsns")
			s, err := GetMessageSenderFromEnv()
			if err != nil {
				t.Errorf("expected success but got error: %s", err)
			}
			expectedType := "*sendsns.Sess"
			senderType := fmt.Sprintf("%T", s)
			if senderType != expectedType {
				t.Errorf("expected %s, got %s", expectedType, senderType)
			}
			// TODO validate SNS properties
		},
	)
}
