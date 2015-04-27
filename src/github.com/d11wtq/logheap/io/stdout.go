package io

import (
	"fmt"
	"net/url"
)

// Output handler for writing to stdout.
type Stdout struct {
	ch chan string
}

// Listen for incoming documents and process them.
func (s *Stdout) Listen() {
	s.ch = make(chan string)
	for doc := range s.ch {
		s.write(doc)
	}
}

// Push a document for processing.
func (s *Stdout) Push(doc string) {
	s.ch <- doc
}

func (s *Stdout) write(doc string) {
	fmt.Println(doc)
}

func init() {
	RegisterOutput(
		"stdout",
		func(u *url.URL) (Output, error) {
			return &Stdout{}, nil
		},
	)
}
