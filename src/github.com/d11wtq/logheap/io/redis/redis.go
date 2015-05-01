package redis

import (
	"github.com/d11wtq/logheap/io"
	"github.com/fzzy/radix/redis"
	"net/url"
)

// Output handler for writing to redis.
type Output struct {
	client *redis.Client
	ch     chan string
}

// Create a new Output for redis processing.
func NewOutput(u *url.URL) (io.Output, error) {
	if client, err := redis.Dial("tcp", u.Host); err == nil {
		return &Output{client: client}, nil
	} else {
		return nil, err
	}
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
	s.client.Cmd("RPUSH", "logstash", doc)
}

func Register() {
	io.RegisterOutput("redis", NewOutput)
}
