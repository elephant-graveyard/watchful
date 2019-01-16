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

	"github.com/homeport/disrupt-o-meter/dom/logger"

	. "github.com/homeport/disrupt-o-meter/dom/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMerkhet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "disrupt-o-meter dom merkhet suite")
}

type MerkhetMock struct {
	Base        Base
	FailedRuns  uint
	TotalRuns   uint
	WillExecute bool
}

func (s *MerkhetMock) Install() {
	s.GetBase().GetLogger().WriteString("Install")
}

func (s *MerkhetMock) PostConnect() {
	s.GetBase().GetLogger().WriteString("PostConnect")
}

func (s *MerkhetMock) Execute() {
	s.GetBase().GetLogger().WriteString("Execute")
}

func (s *MerkhetMock) GetBase() Base {
	return s.Base
}

func (s *MerkhetMock) BuildResult() Result {
	return NewMerkhetResult(s.TotalRuns, s.FailedRuns, s.GetBase().GetConfiguration().IsValidRun(s.TotalRuns, s.FailedRuns))
}

func NewMerkhetMock(config Configuration, totalRuns uint, fails uint, canExecute bool) *MerkhetMock {
	return &MerkhetMock{
		Base:        NewMerkhetBase(NewLoggerMock(), config),
		TotalRuns:   totalRuns,
		FailedRuns:  fails,
		WillExecute: canExecute,
	}
}

//--

type LoggerMock struct {
	Buffer    *bytes.Buffer
	Observers []func(b []byte)
}

func (l *LoggerMock) PushObserver(observer func(bytes []byte)) {
	l.Observers = append(l.Observers, observer)
}

func (l *LoggerMock) GetName() string {
	return "LoggerMock"
}

func (l *LoggerMock) Write(p []byte) (n int, err error) {
	for _, observer := range l.Observers {
		observer(p)
	}

	return l.Buffer.Write(p)
}

func (l *LoggerMock) WriteString(s string) error {
	_, err := l.Write([]byte(s))
	return err
}

func (l *LoggerMock) Clear() []byte {
	b := l.Buffer.Bytes()
	l.Buffer.Reset()
	return b
}

func NewLoggerMock() logger.Logger {
	return &LoggerMock{
		Buffer:    &bytes.Buffer{},
		Observers: make([]func(bytes []byte), 0),
	}
}
