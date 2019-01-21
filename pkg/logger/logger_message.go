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

// ChannelMessage is the content of the Logger channel and contains the logger name as well as the message the logger produced
//
// Logger returns the logger that send this message
//
// Message returns the message in it's raw byte form
//
// MessageAsString returns the message wrapped in a string
type ChannelMessage interface {
	Logger() Logger
	Message() []byte
	MessageAsString() string
}

// SimpleChannelMessage is a small implementation of the ChannelMessage interface
type SimpleChannelMessage struct {
	logger  Logger
	message []byte
}

// Logger returns the logger that send this message
func (c *SimpleChannelMessage) Logger() Logger {
	return c.logger
}

// Message returns the message in it's raw byte form
func (c *SimpleChannelMessage) Message() []byte {
	return c.message
}

// MessageAsString Returns the message as string
func (c *SimpleChannelMessage) MessageAsString() string {
	return string(c.Message())
}

// NewChannelMessage is a simple constructor for the ChannelMessage struct
func NewChannelMessage(logger Logger, message []byte) *SimpleChannelMessage {
	return &SimpleChannelMessage{
		logger:  logger,
		message: message,
	}
}
