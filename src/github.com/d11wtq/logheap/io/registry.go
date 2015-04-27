package io

import "net/url"

// Func that can create an Input from a URL.
type InputFactory func(u *url.URL) (Input, error)

// Func that can create an Output from a URL.
type OutputFactory func(u *url.URL) (Output, error)

// Register a new Input handler.
// Error message for unsupported schemes.
type SchemeError string

func (e SchemeError) Error() string {
	return string(e)
}

// Registry in which to store Input factories
var inputs = make(map[string]InputFactory)

// Registry in which to store Output factories
var outputs = make(map[string]OutputFactory)

//
// The scheme is the URI scheme that triggers the use of the handler.
// The func is provided with the actual URI and returns the Input.
func RegisterInput(scheme string, f InputFactory) {
	inputs[scheme] = f
}

// Register a new Output handler.
//
// The scheme is the URI scheme that triggers the use of the handler.
// The func is provided with the actual URI and returns the Output.
func RegisterOutput(scheme string, f OutputFactory) {
	outputs[scheme] = f
}

// Create the Input for the given URL s.
func NewInput(s string) (Input, error) {
	if u, err := url.Parse(s); err == nil {
		if f, ok := inputs[u.Scheme]; ok {
			return f(u)
		} else {
			return nil, SchemeError(u.Scheme)
		}
	} else {
		return nil, err
	}
}

// Create the Output for the given URL s.
func NewOutput(s string) (Output, error) {
	if u, err := url.Parse(s); err == nil {
		if f, ok := outputs[u.Scheme]; ok {
			return f(u)
		} else {
			return nil, SchemeError(u.Scheme)
		}
	} else {
		return nil, err
	}
}
