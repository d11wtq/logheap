package main

import (
	"github.com/d11wtq/logheap/io"
	"github.com/d11wtq/logheap/io/docker"
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
		defer wait.Done()
		output.Listen()
	}()

	wait.Add(1)
	go func() {
		defer wait.Done()
		input.Listen(output)
	}()

	wait.Wait()
}
