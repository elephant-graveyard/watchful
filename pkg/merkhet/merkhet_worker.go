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

//ControllerChannel is the channel type that is used to send consumers to the running Worker instance
type ControllerChannel chan func(merkhet Merkhet)

//Worker is a wrapper for a created go routine that handles an existing merkhet instance running
type Worker interface {

	//GetControllerChannel returns the controller channel this worker listens on
	GetControllerChannel() ControllerChannel

	//GetMerkhet returns the merkhet instance the Worker wraps
	GetMerkhet() Merkhet

	//Start worker starts the workers listening logic.
	//This has to be called in a different go routine, or it will block the calling routine
	StartWorker()
}

//simpleWorkerTask is a basic implementation of the Worker interface
type simpleWorkerTask struct {
	merkhet Merkhet
	master  ControllerChannel
}

//GetControllerChannel returns the controller channe this worker listens on
func (s *simpleWorkerTask) GetControllerChannel() ControllerChannel {
	return s.master
}

//GetMerkhet returns the merkhet instance the Worker wraps
func (s *simpleWorkerTask) GetMerkhet() Merkhet {
	return s.merkhet
}

//Start worker starts the workers listening logic.
//This has to be called in a different go routine, or it will block the calling routine
func (s *simpleWorkerTask) StartWorker() {
	for request := range s.GetControllerChannel() {
		request(s.GetMerkhet())
	}
}

//NewMerkhetWorker returns a new implementation of the Worker instance
func NewMerkhetWorker(merkhet Merkhet) Worker {
	return &simpleWorkerTask{
		merkhet: merkhet,
		master:  make(ControllerChannel),
	}
}
