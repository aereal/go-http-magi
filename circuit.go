package main

import "fmt"

type Operation int

func (f *Operation) String() string {
	for name, op := range name2op {
		if *f == op {
			return name
		}
	}
	return "UNKNOWN"
}

func (f *Operation) Set(value string) error {
	if op, ok := name2op[value]; ok {
		*f = op
		return nil
	}
	return fmt.Errorf("Unknown value: %s", value)
}

const (
	opAnd Operation = iota
)

var name2op = map[string]Operation{
	"and": opAnd,
}
