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

package services

import (
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
)

// TeardownService tears down the cf test instance
type TeardownService struct {
	WatchfulLogger logger.Logger
	Worker         cfw.CloudFoundryWorker
}

// NewTeardownService creates a new teardown task
func NewTeardownService(watchfulLogger logger.Logger, worker cfw.CloudFoundryWorker) *TeardownService {
	return &TeardownService{WatchfulLogger: watchfulLogger, Worker: worker}
}

// Execute executes a teardown
func (e *TeardownService) Execute() error {
	e.WatchfulLogger.WriteString(logger.Info, "Tearing down test environment")
	if err := e.Worker.TeardownTestEnvironment(); err != nil {
		e.WatchfulLogger.WriteString(logger.Error, "Could not teardown test environment")
		return err
	}

	return nil
}
