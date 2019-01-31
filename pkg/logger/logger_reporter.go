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

import "errors"

// ReporterReviewer defines a method that reviews a byte slice and returns if it should be passed to the logger
type ReporterReviewer func(p []byte) error

// ReportingWriter is a writer that is able to review the written bytes
//
// Write writes the bytes to the writer
//
// ReviewWith sets the Reporters reviewer. A ReportingWriter can have only on reviewer
type ReportingWriter interface {
	Write(p []byte) (n int, err error)
	ReviewWith(reviewer ReporterReviewer) ReportingWriter
}

// ReportingLoggerWriter is a simple internal implementation of the writer interface.
// it forwards every byte slice to its master logger
type ReportingLoggerWriter struct {
	logger   Logger
	level    LogLevel
	reviewer ReporterReviewer
}

// NewSimpleLoggerReporter creates a new reporter that forwards every written byte to the logger
func NewSimpleLoggerReporter(logger Logger, level LogLevel) *ReportingLoggerWriter {
	return &ReportingLoggerWriter{logger: logger, level: level, reviewer: func(p []byte) error {
		return nil
	}}
}

// ReviewWith sets the Reporters reviewer. A ReportingWriter can have only on reviewer
func (l *ReportingLoggerWriter) ReviewWith(reviewer ReporterReviewer) ReportingWriter {
	l.reviewer = reviewer
	return l
}

// Write simply forwards every call to the master
func (l *ReportingLoggerWriter) Write(p []byte) (n int, err error) {
	if l.reviewer == nil {
		panic(errors.New("the logger reporter was initialized without a basic reviewer. Always use the constructor"))
	}

	if e := l.reviewer(p); e != nil {
		return 0, e // Return the error here and don't write it
	}

	return l.logger.Write(p, l.level) // Forward it
}
