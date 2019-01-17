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

package logger

//ChannelProvider is a small provider instance that provides a reference to the channel loggers use to communicate with the logger coupler instance
type ChannelProvider interface {

	//Push pushes the channel message onto the channel
	Push(message ChannelMessage)

	//Read reads the channel message from the wrapped channel
	Read() ChannelMessage

	//GetChannel returns the wrapped channel instance
	GetChannel() chan ChannelMessage
}

//simpleChannelProvider is a basic struct base implementation of the of ChannelProvider interface
type simpleChannelProvider struct {
	Channel     chan ChannelMessage
	ChannelSize int
}

//Push pushes the channel message onto the channel
func (c *simpleChannelProvider) Push(message ChannelMessage) {
	if len(c.Channel) < c.ChannelSize {
		c.GetChannel() <- message
	} else {
		go func() {
			c.GetChannel() <- message
		}()
	}
}

//Read reads the channel message from the wrapped channel
func (c *simpleChannelProvider) Read() ChannelMessage {
	return <-c.GetChannel()
}

func (c *simpleChannelProvider) GetChannel() chan ChannelMessage {
	return c.Channel
}

//NewChannelProvider returns a channel provider that wraps the passed channel instance
func NewChannelProvider(size int) ChannelProvider {
	return &simpleChannelProvider{
		Channel:     make(chan ChannelMessage, size),
		ChannelSize: size,
	}
}
