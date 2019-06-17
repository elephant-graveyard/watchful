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

package merkhets

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gorilla/mux"
	"github.com/homeport/watchful/pkg/logger"
)

func TestMerkhetImplementations(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "watchful internal merkhets")
}

type MockedServer struct {
	Router *mux.Router
	Server *httptest.Server
}

func (m *MockedServer) Route(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	m.Router.HandleFunc(path, handler)
}

func (m *MockedServer) Start() {
	m.Server.Start()
}

func NewMockedServer() *MockedServer {
	router := mux.NewRouter()
	return &MockedServer{
		Server: httptest.NewUnstartedServer(router),
		Router: router,
	}
}

type DevNullLogger struct {
}

func (l *DevNullLogger) AsPrefix() string {
	return "[test]"
}

func (l *DevNullLogger) ChannelProvider() logger.ChannelProvider {
	return nil
}

func (l *DevNullLogger) ID() int {
	return 0
}

func (l *DevNullLogger) Name() string {
	return "console"
}

func (l *DevNullLogger) ReportingOn(level logger.LogLevel) logger.ReportingWriter {
	return &ConsoleLoggerReporter{}
}

func (l *DevNullLogger) Write(p []byte, level logger.LogLevel) (n int, err error) {
	return 0, nil
}

func (l *DevNullLogger) WriteString(level logger.LogLevel, s string) error {
	_, err := l.Write([]byte(s), level)
	return err
}

type ConsoleLoggerReporter struct {
}

func (r *ConsoleLoggerReporter) Refocus(level logger.LogLevel) logger.ReportingWriter {
	return r
}

func (r *ConsoleLoggerReporter) ReviewWith(reviewer logger.ReporterReviewer) logger.ReportingWriter {
	return r
}

func (r *ConsoleLoggerReporter) Write(p []byte) (n int, err error) {
	return 0, nil
}
