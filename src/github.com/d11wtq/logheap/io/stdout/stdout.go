package stdout

import (
	"fmt"
	"github.com/d11wtq/logheap/io"
	"net/url"
)

// Output handler for writing to stdout.
type Output struct {
	ch chan string
}

// Listen for incoming documents and process them.
func (s *Output) Listen() {
	s.ch = make(chan string)
	for doc := range s.ch {
		s.write(doc)
	}
}

// Push a document for processing.
func (s *Output) Push(doc string) {
	s.ch <- doc
}

func (s *Output) write(doc string) {
	fmt.Println(doc)
}

func Register() {
	io.RegisterOutput(
		"stdout",
		func(u *url.URL) (io.Output, error) {
			return &Output{}, nil
		},
	)
}
