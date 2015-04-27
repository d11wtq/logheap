package io

import "sync"

// Interface for input handlers.
type Input interface {
	// Listen for incoming documents and process them.
	Listen(output Output)
}

// Interface for output handlers.
type Output interface {
	// Push a document to the handler for procesing.
	//
	// The document will typically be passed to an internal channel.
	Push(doc string)

	// Listen for incoming documents and process them.
	//
	// Documents will typically be pulled from an internal channel.
	Listen()
}

// Presents a collection of Inputs as a single input.
type UnionInput []Input

func (i *UnionInput) Listen(o Output) {
	var wait sync.WaitGroup

	for _, v := range *i {
		wait.Add(1)
		go func() {
			defer wait.Done()
			v.Listen(o)
		}()
	}

	wait.Wait()
}

// Presents a collection of Outputs as a single output.
type UnionOutput []Output

func (o *UnionOutput) Push(doc string) {
	for _, v := range *o {
		v.Push(doc)
	}
}

func (o *UnionOutput) Listen() {
	var wait sync.WaitGroup

	for _, v := range *o {
		wait.Add(1)
		go func() {
			defer wait.Done()
			v.Listen()
		}()
	}

	wait.Wait()
}
