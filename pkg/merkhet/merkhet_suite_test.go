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

package merkhet_test

import (
	"bytes"
	"testing"

	"github.com/homeport/watchful/pkg/logger"
	. "github.com/homeport/watchful/pkg/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMerkhet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "watchful pkg merkhet suite")
}

type MerkhetCallback struct {
	onInstall     func() error
	onPostConnect func() error
	onExecute     func() error
}

type MerkhetMock struct {
	Config         Configuration
	LoggerMock     logger.Logger
	FailedRuns     int
	SuccessfulRuns int
	WillExecute    bool
	Callback       *MerkhetCallback
}

func (s *MerkhetMock) RecordSuccessfulRun() {
	s.SuccessfulRuns = s.SuccessfulRuns + 1
}

func (s *MerkhetMock) RecordFailedRun() {
	s.FailedRuns = s.FailedRuns + 1
}

func (s *MerkhetMock) Install() error {
	if s.Callback.onInstall != nil {
		return s.Callback.onInstall()
	}
	return nil
}

func (s *MerkhetMock) PostConnect() error {
	if s.Callback.onPostConnect != nil {
		return s.Callback.onPostConnect()
	}
	return nil
}

func (s *MerkhetMock) Execute() error {
	if s.Callback.onExecute != nil {
		return s.Callback.onExecute()
	}
	return nil
}

func (s *MerkhetMock) BuildResult() Result {
	return NewMerkhetResult(s.SuccessfulRuns, s.FailedRuns, s.Configuration().ValidRun(s.SuccessfulRuns, s.FailedRuns))
}

func (s *MerkhetMock) Configuration() Configuration {
	return s.Config
}

func (s *MerkhetMock) Logger() logger.Logger {
	return s.LoggerMock
}

func NewMerkhetMock(config Configuration, totalRuns int, fails int, canExecute bool, callback *MerkhetCallback) *MerkhetMock {
	return &MerkhetMock{
		SuccessfulRuns: totalRuns,
		FailedRuns:     fails,
		WillExecute:    canExecute,
		Config:         config,
		LoggerMock:     NewLoggerMock(),
		Callback:       callback,
	}
}

type LoggerMock struct {
	Buffer *bytes.Buffer
}

func (l *LoggerMock) ReportingOn(level logger.LogLevel) logger.ReportingWriter {
	return logger.NewSimpleLoggerReporter(l, level)
}

func (l *LoggerMock) Name() string {
	return "LoggerMock"
}

func (l *LoggerMock) ID() int {
	return 0
}

func (l *LoggerMock) Write(p []byte, level logger.LogLevel) (n int, err error) {
	return l.Buffer.Write(p)
}

func (l *LoggerMock) WriteString(level logger.LogLevel, s string) error {
	_, err := l.Write([]byte(s), level)
	return err
}

func (l *LoggerMock) ChannelProvider() logger.ChannelProvider {
	return nil
}

func NewLoggerMock() logger.Logger {
	return &LoggerMock{
		Buffer: &bytes.Buffer{},
	}
}
