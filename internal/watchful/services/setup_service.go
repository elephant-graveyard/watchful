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
	"github.com/homeport/watchful/internal/watchful/cfg"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
)

// SetupService is a service to setup a new cloud foundry
type SetupService struct {
	Config             cfg.CloudFoundryConfig
	WatchfulLogger     logger.Logger
	CloudFoundryLogger logger.Logger
	Worker             cfw.CloudFoundryWorker
}

// NewSetupService creates a new instance
func NewSetupService(config cfg.CloudFoundryConfig, watchfulLogger logger.Logger, cloudFoundryLogger logger.Logger, worker cfw.CloudFoundryWorker) *SetupService {
	return &SetupService{Config: config, WatchfulLogger: watchfulLogger, CloudFoundryLogger: cloudFoundryLogger, Worker: worker}
}

// Execute sets up the cloud foundry instance
func (e *SetupService) Execute() error {
	e.WatchfulLogger.WriteString(logger.Info, "Targeting API Endpoint")
	if err := e.Worker.API(e.Config.APIEndPoint, !e.Config.SkipSSLValidation); err != nil {
		e.WatchfulLogger.WriteString(logger.Error, "Could not target API endpoint "+err.Error())
		return err
	}

	e.WatchfulLogger.WriteString(logger.Info, "Authenticating against API endpoint")
	e.WatchfulLogger.WriteString(logger.Info, "Username: "+e.Config.Username)
	e.WatchfulLogger.WriteString(logger.Info, "Passkey: "+e.Config.Password)
	if err := e.Worker.Authenticate(cfw.CloudFoundryCertificate{
		Username: e.Config.Username,
		Password: e.Config.Password,
	}); err != nil {
		e.WatchfulLogger.WriteString(logger.Error, "Could not authenticate against API endpoint "+err.Error())
		return err
	}

	e.WatchfulLogger.WriteString(logger.Info, "Creating test environment")
	if err := e.Worker.CreateTestEnvironment(); err != nil {
		e.WatchfulLogger.WriteString(logger.Error, "Could not create test environment")
		return err
	}
	return nil
}
