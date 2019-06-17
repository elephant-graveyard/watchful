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
	"os/exec"

	"github.com/homeport/watchful/internal/watchful/cfg"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
)

// CloudFoundryService is an services that executes a cloud foundry task
type CloudFoundryService struct {
	Tasks              []cfg.TaskConfiguration
	CloudFoundryLogger logger.Logger
}

// NewCloudFoundryService creates a new cloud foundry executor service
func NewCloudFoundryService(tasks []cfg.TaskConfiguration, cloudFoundryLogger logger.Logger) *CloudFoundryService {
	return &CloudFoundryService{Tasks: tasks, CloudFoundryLogger: cloudFoundryLogger}
}

// Next returns if the cloud foundry services has a next task to run
// It will also return the next task it would execute
func (e *CloudFoundryService) Next() (configuration *cfg.TaskConfiguration) {
	if len(e.Tasks) > 0 {
		return &e.Tasks[0]
	}
	return nil
}

// Execute executes the current task
func (e *CloudFoundryService) Execute() error {
	config := e.Next()

	commandPromise := cfw.NewSimpleCommandPromise(exec.Command(config.Executable, config.Parameters...))
	commandPromise.SubscribeOnOut(e.CloudFoundryLogger.ReportingOn(logger.Info))
	commandPromise.SubscribeOnErr(e.CloudFoundryLogger.ReportingOn(logger.Error))

	return commandPromise.Sync()
}

// Pop the first cloud foundry services
func (e *CloudFoundryService) Pop() *cfg.TaskConfiguration {
	if len(e.Tasks) > 0 {
		top := e.Tasks[0]
		e.Tasks = e.Tasks[1:]
		return &top
	}
	return nil
}
