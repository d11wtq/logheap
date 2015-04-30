package stdin

import (
	"bufio"
	"github.com/d11wtq/logheap/io"
	"net/url"
	"os"
)

// Input handler for reading from stdin.
type Input struct{}

// Listen for incoming documents and process them.
func (s *Input) Listen(o io.Output) {
	tags := map[string]interface{}{"type": "stdin"}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if doc, err := io.Encode(scanner.Text(), tags); err == nil {
			o.Push(doc)
		}
	}
}

func Register() {
	io.RegisterInput(
		"stdin",
		func(u *url.URL) (io.Input, error) {
			return &Input{}, nil
		},
	)
}
