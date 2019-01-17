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

import "bytes"

//Logger defines an observable logger
type Logger interface {
	//PushObserver adds an observer function to the logger.
	//This function is executed whenever a write call is issued against the Logger intance
	PushObserver(func(bytes []byte))

	//GetName returns the name of the logger
	GetName() string

	//Write simply writes a byte array to the logger. This will trigger all observers
	Write(p []byte) (n int, err error)

	//WriteString writes all bytes of the string to the logger.
	//This will also trigger observers
	WriteString(s string) error

	//Peek returns the currently stored logged bytes in the logger
	Peek() []byte

	//Clear removes all logged bytes from the Logger and returns them for future use
	Clear() []byte
}

//memStoredLogger is an implementation of the Logger interface that the loggers will use in order to store their logged values
//It does not flush the written messages, instead it will store them in the buffer, ready to be read from the Coupler
type memStoredLogger struct {
	buffer    bytes.Buffer
	observers []func(bytes []byte)
	name      string
}

//PushObserver simply adds an observer to slice of registered observers
func (m *memStoredLogger) PushObserver(observer func(bytes []byte)) {
	m.observers = append(m.observers, observer)
}

//GetName simply returns the name the logger instance was assigned on creation
func (m *memStoredLogger) GetName() string {
	return m.name
}

//Write simply stores the written string inside the buffer
func (m *memStoredLogger) Write(b []byte) (int, error) {
	for _, ob := range m.observers {
		ob(b)
	}

	return m.buffer.Write(b)
}

//WriteString writes an entire string to the logger instance
func (m *memStoredLogger) WriteString(s string) error {
	_, err := m.Write([]byte(s))
	return err
}

//Peek returns the bytes currently stored in the logger
func (m *memStoredLogger) Peek() []byte {
	return m.buffer.Bytes()
}

//Clear simply clears and resets the buffer.
//This returns the stored values in the logger
func (m *memStoredLogger) Clear() []byte {
	bytes := m.Peek()
	m.buffer.Reset()
	return bytes
}

//NewMemStoredLogger creates a new instance of the memStoredLogger
func NewMemStoredLogger(name string) Logger {
	return &memStoredLogger{
		observers: make([]func(bytes []byte), 0),
		name:      name,
		buffer:    bytes.Buffer{},
	}
}
