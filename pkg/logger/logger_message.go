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

//--

//ChannelMessage is the content of the Logger channel and contains the logger name as well as the message the logger produced
type ChannelMessage interface {
	//GetLogger returns the logger that send this message
	GetLogger() Logger

	//GetMessage returns the message in it's raw byte form
	GetMessage() []byte

	//GetMessageAsString returns the message wrapped in a string
	GetMessageAsString() string
}

//simpleChannelMessage is a small implementation of the ChannelMessage interface
type simpleChannelMessage struct {
	logger  Logger
	message []byte
}

//GetLogger returns the logger that send this message
func (c *simpleChannelMessage) GetLogger() Logger {
	return c.logger
}

//GetMessage returns the message in it's raw byte form
func (c *simpleChannelMessage) GetMessage() []byte {
	return c.message
}

//Returns the message as string
func (c *simpleChannelMessage) GetMessageAsString() string {
	return string(c.GetMessage())
}

//NewChannelMessage is a simple constructor for the ChannelMessage struct
func NewChannelMessage(logger Logger, message []byte) ChannelMessage {
	return &simpleChannelMessage{
		logger:  logger,
		message: message,
	}
}
