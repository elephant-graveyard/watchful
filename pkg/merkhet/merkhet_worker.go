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

// ControllerChannel is the channel type that is used to send consumers to the running Worker instance
type ControllerChannel chan func(merkhet Merkhet)

// Worker is a wrapper for a created go routine that handles an existing merkhet instance running
//
// ControllerChannel returns the controller channel this worker listens on
//
// Merkhet returns the merkhet instance the Worker wraps
//
// StartWorker starts the workers listening logic.
// This has to be called in a different go routine, or it will block the calling routine
type Worker interface {
	ControllerChannel() ControllerChannel
	Merkhet() Merkhet
	StartWorker()
}

// SimpleWorkerTask is a basic implementation of the Worker interface
type SimpleWorkerTask struct {
	merkhet Merkhet
	master  ControllerChannel
}

//ControllerChannel returns the controller channe this worker listens on
func (s *SimpleWorkerTask) ControllerChannel() ControllerChannel {
	return s.master
}

// Merkhet returns the merkhet instance the Worker wraps
func (s *SimpleWorkerTask) Merkhet() Merkhet {
	return s.merkhet
}

// StartWorker starts the workers listening logic.
// This has to be called in a different go routine, or it will block the calling routine
func (s *SimpleWorkerTask) StartWorker() {
	for request := range s.ControllerChannel() {
		request(s.Merkhet())
	}
}

// NewMerkhetWorker returns a new implementation of the Worker instance
func NewMerkhetWorker(merkhet Merkhet) *SimpleWorkerTask {
	return &SimpleWorkerTask{
		merkhet: merkhet,
		master:  make(ControllerChannel),
	}
}
