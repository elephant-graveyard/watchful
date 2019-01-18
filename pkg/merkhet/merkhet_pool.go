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

//Pool contains a map of registered merkhets as well as manages them
type Pool interface {
	//StartWorker pushes a new Merkhet instance to the Pool
	StartWorker(merkhet Merkhet)

	//Size returns the current size of the Pool
	Size() uint

	//ForEach executes the provided function for each Merkhet instance currently managed by the Pool
	ForEach(consumer func(m Merkhet))

	//Shutdown shuts the pool and it's workers down
	Shutdown()
}

//simplePool is a basic implementation of the MerkhetPool interface
type simplePool struct {
	Workers []Worker
}

//StartWorker pushes a new merkhet instance into the Pool
func (s *simplePool) StartWorker(m Merkhet) {
	worker := NewMerkhetWorker(m)
	go worker.StartWorker() //Start the worker intance in a different go routine

	s.Workers = append(s.Workers, worker)
}

//Size returns the size of the pool
func (s *simplePool) Size() uint {
	return uint(len(s.Workers))
}

//ForEach executes the provided function for each Merkhet instance currently managed by the pool
func (s *simplePool) ForEach(consumer func(m Merkhet)) {
	for _, worker := range s.Workers {
		worker.GetControllerChannel() <- consumer
	}
}

//Shutdown shuts the pool and it's workers down
func (s *simplePool) Shutdown() {
	for _, worker := range s.Workers {
		close(worker.GetControllerChannel())
	}
}

//NewPool returns a fresh empty instance of the go container
func NewPool() Pool {
	return &simplePool{}
}
