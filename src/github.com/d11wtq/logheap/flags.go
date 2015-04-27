package main

import (
	"flag"
	"fmt"
)

// Stores a list of string flags.
type FlagList []string

// Sets a new value into the list.
// This is part of the flag.Value interface.
func (f *FlagList) Set(s string) error {
	*f = append(*f, s)
	return nil
}

// Gets the current value of the list.
// This is part of the flag.Value interface.
func (f *FlagList) String() string {
	return fmt.Sprint(*f)
}

// FlagParser parses command line arguments and records what flags were given.
type FlagParser struct {
	Input  FlagList
	Output FlagList
}

// Populate p.Input and p.Output slices.
//
// If no arguments are found, defaults are assigned, where applicable.
func (p *FlagParser) Parse() {
	flag.Var(&p.Input, "in", "Input source URIs")
	flag.Var(&p.Output, "out", "Output destination URIs")

	flag.Parse()

	if len(p.Input) == 0 {
		p.Input = FlagList{"stdin:"}
	}

	if len(p.Output) == 0 {
		p.Output = FlagList{"stdout:"}
	}
}
