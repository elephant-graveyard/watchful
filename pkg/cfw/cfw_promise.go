// Copyright © 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cfw

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// CommandPromise is a basic promise interface that promises a commands execution future
//
// SubscribeOnOut subscribes the provided writer to the output stream of the command promise
// This will overwrite previous subscriber
//
// SubscribeOnErr subscribes the provided writer to the error stream of the command promise
// This will overwrite previous subscriber
//
// Sync executes the command in sync to the go routine it was called in, returning the result
//
// Async executes the command in a new go routine, calls the passed callback after execution and returns a wait-group
// that is tracking the state of the execution
type CommandPromise interface {
	SubscribeOnOut(writer io.Writer) CommandPromise
	SubscribeOnErr(writer io.Writer) CommandPromise
	Environment(key string, value string) CommandPromise
	Sync() error
	Async(subscriber func(e error)) *sync.WaitGroup
}

// NewSimpleSlicedCommandPromise returns a simple command promise implementation containing multiple commands
func NewSimpleSlicedCommandPromise(commands []*exec.Cmd) *SimpleCommandPromise {
	return &SimpleCommandPromise{commands: commands}
}

// NewSimpleCommandPromise returns a simple command promise implementation containing one command
func NewSimpleCommandPromise(command *exec.Cmd) *SimpleCommandPromise {
	return &SimpleCommandPromise{commands: []*exec.Cmd{command}}
}

// SimpleCommandPromise is an implementation of the CommandPromise interface that is able to run multiple commands
type SimpleCommandPromise struct {
	commands []*exec.Cmd
}

// SubscribeOnOut will subscribe the writer instance to the command promise
func (c *SimpleCommandPromise) SubscribeOnOut(writer io.Writer) CommandPromise {
	for _, c := range c.commands {
		c.Stdout = writer
	}
	return c
}

// SubscribeOnErr will subscribe the writer instance to the command promise
func (c *SimpleCommandPromise) SubscribeOnErr(writer io.Writer) CommandPromise {
	for _, c := range c.commands {
		c.Stderr = writer
	}
	return c
}

// Environment adds a new environment variable
func (c *SimpleCommandPromise) Environment(key string, value string) CommandPromise {
	for _, c := range c.commands {
		if len(c.Env) < 1 {
			c.Env = os.Environ() // Store the default env anyway !
		}
		c.Env = append(c.Env, key+"="+value)
	}
	return c
}

// Sync executes the command promise and returns the result
// The command will be executed on the same go routine
func (c *SimpleCommandPromise) Sync() error {
	for _, c := range c.commands {
		e := c.Run()
		if e != nil {
			return e
		}
	}
	return nil
}

// Async executes the command promise and returns the result to the passed subscriber
// The command will be executed on a new go routine
// It returns a WaitGroup instance that will finish once the task did
func (c *SimpleCommandPromise) Async(subscriber func(e error)) *sync.WaitGroup {
	w := &sync.WaitGroup{}
	w.Add(1)

	go func(w *sync.WaitGroup) {
		result := c.Sync()

		if subscriber != nil {
			subscriber(result)
		}
		w.Done()
	}(w)

	return w
}
