package flagarray

import (
	"testing"
)

func TestSet(t *testing.T) {
	t.Run("with string",
		func(t *testing.T) {
			obj := StringArray{}
			if len(obj) > 0 {
				t.Errorf("expected len(obj) == 0 but got %d", len(obj))
			}

			err := obj.Set("one")
			if err != nil {
				t.Error("first Set failed, got error:", err)
			}
			if len(obj) != 1 {
				t.Errorf("expected len(obj) == 1 but got %d; obj=%v", len(obj), obj)
			}
			if obj[0] != "one" {
				t.Errorf("first Set stored wrong value; obj=%v", obj)
			}

			err = obj.Set("two")
			if err != nil {
				t.Error("second Set failed, got error:", err)
			}
			if len(obj) != 2 {
				t.Errorf("expected len(obj) == 2 but got %d; obj=%v", len(obj), obj)
			}
			if obj[1] != "two" {
				t.Errorf("second Set stored wrong value; obj=%v", obj)
			}
		},
	)
}

func TestString(t *testing.T) {
	t.Run("with string",
		func(t *testing.T) {
			obj := StringArray{"one", "two"}
			obj.Set("three")
			s := obj.String()
			if s != "one two three" {
				t.Errorf("incorrect String() response, got %v", s)
			}
		},
	)
}
