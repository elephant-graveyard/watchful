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

// LogLevel describes the log level type a logger can write with
type LogLevel byte

const (
	// Info represents the logging level info.
	Info LogLevel = iota

	// Error represents the logging level info
	Error
)

// Logger defines an observable logger
//
// ChannelProvider returns the channel provider this Cluster uses
//
// Name returns the name of the logger
//
// ID returns the id of the logger
//
// Write simply writes a byte array to the logger
//
// WriteString writes all bytes of the string to the logger
//
// ReportingOn creates a new writer that reports every written byte slice to the logger on the given log level
type Logger interface {
	ChannelProvider() ChannelProvider
	Name() string
	ID() int
	Write(p []byte, level LogLevel) (n int, err error)
	WriteString(level LogLevel, s string) error
	ReportingOn(level LogLevel) ReportingWriter
}

// SimpleChanneledLogger is an implementation of the Logger interface that the loggers will use in order to store their logged values
// It does not flush the written messages, instead it will store them in the buffer, ready to be read from the Coupler
type SimpleChanneledLogger struct {
	channelProvider ChannelProvider
	name            string
	id              int
}

// ChannelProvider returns the channel provider the Cluster wi
func (l *SimpleChanneledLogger) ChannelProvider() ChannelProvider {
	return l.channelProvider
}

// Name simply returns the name the logger instance was assigned on creation
func (l *SimpleChanneledLogger) Name() string {
	return l.name
}

// ID returns the id the logger instance was assigned on creation
func (l *SimpleChanneledLogger) ID() int {
	return l.id
}

// Write simply stores the written string inside the buffer
func (l *SimpleChanneledLogger) Write(b []byte, level LogLevel) (int, error) {
	l.ChannelProvider().Push(NewChannelMessage(l, b, level))
	return len(b), nil
}

// WriteString writes an entire string to the logger instance
func (l *SimpleChanneledLogger) WriteString(level LogLevel, s string) error {
	_, err := l.Write([]byte(s), level)
	return err
}

// ReportingOn creates a new writer that reports every written byte slice to the logger on the given log level
func (l *SimpleChanneledLogger) ReportingOn(level LogLevel) ReportingWriter {
	return NewSimpleLoggerReporter(l , level)
}
