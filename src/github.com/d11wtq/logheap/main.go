package main

import (
	"github.com/d11wtq/logheap/io"
	"github.com/d11wtq/logheap/io/docker"
	"github.com/d11wtq/logheap/io/redis"
	"github.com/d11wtq/logheap/io/stdin"
	"github.com/d11wtq/logheap/io/stdout"
	"sync"
)

func GetOutputs(col *io.UnionOutput, urls []string) {
	for _, s := range urls {
		if o, err := io.NewOutput(s); err == nil {
			*col = append(*col, o)
		}
	}
}

func GetInputs(col *io.UnionInput, urls []string) {
	for _, s := range urls {
		if o, err := io.NewInput(s); err == nil {
			*col = append(*col, o)
		}
	}
}

func init() {
	stdin.Register()
	stdout.Register()
	docker.Register()
	redis.Register()
}

func main() {
	var wait sync.WaitGroup

	flags := new(FlagParser)
	flags.Parse()

	output := &io.UnionOutput{}
	GetOutputs(output, flags.Output)

	input := &io.UnionInput{}
	GetInputs(input, flags.Input)

	wait.Add(1)
	go func() {
		output.Listen()
		wait.Done()
	}()

	wait.Add(1)
	go func() {
		input.Listen(output)
		wait.Done()
	}()

	wait.Wait()
}
