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

//Factory is an interface that is capable of creating logger instances
type Factory interface {

	//NewChanneledLogger returns an new logger with the given name
	NewChanneledLogger(name string) Logger
}

//channeledLoggerFactory is a basic implementation of the Factory interface
type channeledLoggerFactory struct {
	loggerCount     int
	channelProvider ChannelProvider
}

//NewChanneledLogger returns a new logger instance.
//It's id will simply be based on the amount of logger this factory created
func (c *channeledLoggerFactory) NewChanneledLogger(name string) Logger {
	loggerID := c.loggerCount
	c.loggerCount = c.loggerCount + 1

	return &simpleChanneledLogger{
		id:             loggerID,
		name:           name,
		channelProvier: c.channelProvider,
	}
}

//NewChanneledLoggerFactory creates a new instance of the factory, based on the channel provider instance
func NewChanneledLoggerFactory(channelProvider ChannelProvider) Factory {
	return &channeledLoggerFactory{
		channelProvider: channelProvider,
		loggerCount:     0,
	}
}
