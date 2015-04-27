package main

import (
	"github.com/d11wtq/logheap/docker"
	"github.com/d11wtq/logheap/io"
	"sync"
)

func main() {
	var wait sync.WaitGroup

	flags := new(FlagParser)
	flags.Parse()

	output := &io.UnionOutput{new(io.Stdout)}
	input := &io.UnionInput{new(docker.Input)}

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
