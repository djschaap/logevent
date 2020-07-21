package flagarray

import (
	"strings"
)

// StringArray is a simple array of strings, used to store one or more values provided for a flag.
type StringArray []string

// Set appends a value to a StringArray.
func (strArray *StringArray) Set(v string) error {
	*strArray = append(*strArray, v)
	return nil
}

func (strArray *StringArray) String() string {
	return strings.Join(*strArray, " ")
}
