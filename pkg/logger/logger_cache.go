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

// CachedLogger defines a logger wrapper that caches the messages send and will only flushes them when called fr
//
// Write writes the bytes to the cached logger
//
// Clear clears the content of the logger
//
// Flush flushes the logger content
type CachedLogger interface {
	Write(b []byte) (n int, err error)
	Clear()
	Flush()
}

// ByteBufferCachedLogger is an implementation of the CachedLogger interface relying on byte buffers
type ByteBufferCachedLogger struct {
	cache   []string
	wrapped ReportingWriter
}

// Write writes the bytes to the logger
func (l *ByteBufferCachedLogger) Write(b []byte) (n int, err error) {
	l.cache = append(l.cache, string(b))
	return len(b), nil
}

// Flush flushes the content to the wrapped ReportingWriter
func (l *ByteBufferCachedLogger) Flush() {
	for _, message := range l.cache {
		l.wrapped.Write([]byte(message))
	}
}

// Clear clears the content of the logger
func (l *ByteBufferCachedLogger) Clear() {
	l.cache = nil
}

// NewByteBufferCachedLogger creates a new cached logger that wraps the reporting writer
func NewByteBufferCachedLogger(wrapped ReportingWriter) *ByteBufferCachedLogger {
	return &ByteBufferCachedLogger{wrapped: wrapped}
}
