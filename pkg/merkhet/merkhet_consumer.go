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

// ConsumerMethod describes the type of the methods allowed to consumer a merkhet
type ConsumerMethod func(merkhet Merkhet, relay ControllerChannel)

// Consumer defines an object that can consume a merkhet instnace
type Consumer interface {
	Consume(merkhet Merkhet, relay ControllerChannel)
}

// SyncedConsumer is a consumer that executes the consuming method in sync with the current go routine
type SyncedConsumer struct {
	consumer ConsumerMethod
}

// Consume calls the passed consumer method in the current go routine
func (s *SyncedConsumer) Consume(merkhet Merkhet, relay ControllerChannel) {
	s.consumer(merkhet, relay)
}

// ConsumeSync creates a synced Consumer
func ConsumeSync(m ConsumerMethod) Consumer {
	return &SyncedConsumer{
		consumer: m,
	}
}

// AsyncedConsumer is a consumer that executes the consuming method on a new go routine
type AsyncedConsumer struct {
	consumer ConsumerMethod
}

// Consume calls the passed consumer method in a new go routine
func (a *AsyncedConsumer) Consume(merkhet Merkhet, relay ControllerChannel) {
	go a.consumer(merkhet, relay)
}

// ConsumeAsync creates a asynced Consumer
func ConsumeAsync(m ConsumerMethod) Consumer {
	return &AsyncedConsumer{
		consumer: m,
	}
}
