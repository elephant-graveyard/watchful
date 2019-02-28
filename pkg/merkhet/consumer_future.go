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

package merkhet

import "sync"

// Future represents the future of a consumer
//
// Complete completes the future with an error or nil
//
// IsCompleted returns if the future has been completed and
// if the first return parameter is true, what error was passed
// not that this error can be nil
type Future interface {
	Complete(err error)
	IsCompleted() (completed bool, err error)
}

// SimpleFuture is a basic future implementation
type SimpleFuture struct {
	completed bool
	e         error
	lock      *sync.Mutex
}

// NewFuture creates a new future instance
func NewFuture() *SimpleFuture {
	return &SimpleFuture{
		completed: false,
		lock:      &sync.Mutex{},
	}
}

// Complete completes the future with an error or nil
func (f *SimpleFuture) Complete(err error) {
	defer f.lock.Unlock()

	f.lock.Lock()
	f.e = err
	f.completed = true
}

// IsCompleted returns if the future has been completed and
// if the first return parameter is true, what error was passed
// not that this error can be nil
func (f *SimpleFuture) IsCompleted() (completed bool, err error) {
	defer f.lock.Unlock()

	f.lock.Lock()
	if f.completed {
		return true, f.e
	}
	return false, nil
}
