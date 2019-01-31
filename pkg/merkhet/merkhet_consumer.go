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

// Consumer is a simple consumer method that can be called against a merkhet
type Consumer struct {
	consumerMethod ConsumerMethod
	sync           bool
}

// Consume consumes the passed merkhet instance
func (c Consumer) Consume(m Merkhet, relay ControllerChannel) {
	c.consumerMethod(m, relay)
}

// Sync returns if the consumer can run synced with its worker
func (c Consumer) Sync() bool {
	return c.sync
}

// ConsumeSync creates a synced Consumer
func ConsumeSync(m ConsumerMethod) Consumer {
	return Consumer{
		consumerMethod: m,
		sync:           true,
	}
}

// ConsumeAsync creates a synced Consumer
func ConsumeAsync(m ConsumerMethod) Consumer {
	return Consumer{
		consumerMethod: m,
		sync:           false,
	}
}
