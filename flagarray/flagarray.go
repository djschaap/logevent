package flagarray

import (
	"strings"
)

type StringArray []string

func (self *StringArray) Set(v string) error {
	*self = append(*self, v)
	return nil
}

func (self *StringArray) String() string {
	return strings.Join(*self, " ")
}
