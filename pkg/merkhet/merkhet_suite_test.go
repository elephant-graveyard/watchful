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

	"github.com/homeport/disrupt-o-meter/pkg/logger"
	. "github.com/homeport/disrupt-o-meter/pkg/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMerkhet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "disrupt-o-meter pkg merkhet suite")
}

type MerketCallback struct {
	onInstall     func()
	onPostConnect func()
	onExecute     func()
}

type MerkhetMock struct {
	Config      Configuration
	LokggerMock logger.Logger
	FailedRuns  uint
	TotalRuns   uint
	WillExecute bool
	Callback    *MerketCallback
}

func (s *MerkhetMock) Install() {
	if s.Callback.onInstall != nil {
		s.Callback.onInstall()
	}
}

func (s *MerkhetMock) PostConnect() {
	if s.Callback.onPostConnect != nil {
		s.Callback.onPostConnect()
	}
}

func (s *MerkhetMock) Execute() {
	if s.Callback.onExecute != nil {
		s.Callback.onExecute()
	}
}

func (s *MerkhetMock) BuildResult() Result {
	return NewMerkhetResult(s.TotalRuns, s.FailedRuns, s.Configuration().ValidRun(s.TotalRuns, s.FailedRuns))
}

func (s *MerkhetMock) Configuration() Configuration {
	return s.Config
}

func (s *MerkhetMock) Logger() logger.Logger {
	return s.LokggerMock
}

func NewMerkhetMock(config Configuration, totalRuns uint, fails uint, canExecute bool, callback *MerketCallback) *MerkhetMock {
	return &MerkhetMock{
		TotalRuns:   totalRuns,
		FailedRuns:  fails,
		WillExecute: canExecute,
		Config:      config,
		LokggerMock: NewLoggerMock(),
		Callback:    callback,
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

func (l *LoggerMock) WriteString(s string, level logger.LogLevel) error {
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
