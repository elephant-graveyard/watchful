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

//Logger defines an observable logger
type Logger interface {
	//GetChannelProvider returns the channel provider this Cluster uses
	GetChannelProvider() ChannelProvider

	//GetName returns the name of the logger
	GetName() string

	//Write simply writes a byte array to the logger. This will trigger all observers
	Write(p []byte) (n int, err error)

	//WriteString writes all bytes of the string to the logger.
	//This will also trigger observers
	WriteString(s string) error
}

//simpleChanneledLogger is an implementation of the Logger interface that the loggers will use in order to store their logged values
//It does not flush the written messages, instead it will store them in the buffer, ready to be read from the Coupler
type simpleChanneledLogger struct {
	channelProvier ChannelProvider
	name           string
}

//GetChannelProvider returns the channel provider the Cluster wi
func (m *simpleChanneledLogger) GetChannelProvider() ChannelProvider {
	return m.channelProvier
}

//GetName simply returns the name the logger instance was assigned on creation
func (m *simpleChanneledLogger) GetName() string {
	return m.name
}

//Write simply stores the written string inside the buffer
func (m *simpleChanneledLogger) Write(b []byte) (int, error) {
	m.GetChannelProvider().Push(NewChannelMessage(m, b))
	return len(b), nil
}

//WriteString writes an entire string to the logger instance
func (m *simpleChanneledLogger) WriteString(s string) error {
	_, err := m.Write([]byte(s))
	return err
}

//NewChanneledLogger creates a new instance of the memStoredLogger
func NewChanneledLogger(name string, channelProvider ChannelProvider) Logger {
	return &simpleChanneledLogger{
		name:           name,
		channelProvier: channelProvider,
	}
}
