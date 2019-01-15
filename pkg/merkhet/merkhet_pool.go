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

// Pool contains a map of registered merkhets as well as manages them
//
// StartWorker pushes a new Merkhet instance to the Pool.
// This method will also instantly start the worker sub routine
//
// Size returns the current size of the Pool
//
// ForEach executes the provided function for each Merkhet instance currently managed by the Pool
//
// Shutdown shuts the pool and it's workers down
type Pool interface {
	StartWorker(merkhet Merkhet)
	Size() uint
	ForEach(consumer func(m Merkhet))
	Shutdown()
}

// SimplePool is a basic implementation of the MerkhetPool interface
type SimplePool struct {
	workers []Worker
}

// StartWorker pushes a new merkhet instance into the Pool
func (s *SimplePool) StartWorker(m Merkhet) {
	worker := NewMerkhetWorker(m)
	go worker.StartWorker() //Start the worker intance in a different go routine

	s.workers = append(s.workers, worker)
}

// Size returns the size of the pool
func (s *SimplePool) Size() uint {
	return uint(len(s.workers))
}

// ForEach executes the provided function for each Merkhet instance currently managed by the pool
func (s *SimplePool) ForEach(consumer func(m Merkhet)) {
	for _, worker := range s.workers {
		worker.ControllerChannel() <- consumer
	}
}

// Shutdown shuts the pool and it's workers down
func (s *SimplePool) Shutdown() {
	for _, worker := range s.workers {
		close(worker.ControllerChannel())
	}
}

// NewPool returns a fresh empty instance of the go container
func NewPool() *SimplePool {
	return &SimplePool{}
}
