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

package cfw

import (
	"github.com/homeport/disrupt-o-meter/pkg/logger"
)

// CloudFoundryWorker defines a worker object that is capable of executing specific commands
// against a cloud foundry instance. It will do this in the go routine the method call was issued in,
// so make sure to spin up a new routine if you don't want to block your main process
//
// Authenticate authenticates the worker against the cloud foundry instance provided in the certificate
// It returns any error that may occure, or nil if the opperation was successful
//
// Logger returns the logger instance the worker task is using
//
// Execute executes the provided task against the cloud foundry instance it is working on
// If the worker is not authenticated it will error, else it will try to execute the
// task on the same go routine it was called in
type CloudFoundryWorker interface {
	Authenticate(cert CloudFoundryCertificate) error
	Logger() logger.Logger
	Execute(task CloudFoundryTask) error
}

// SimpleCloudFoundryWorker is a basic implementation of the CloudFoundryWorker interface
type SimpleCloudFoundryWorker struct {
	logger logger.Logger
	cli    CloudFoundryCLI
}

// Authenticate authenticates the worker against the cloud foundry instance provided in the certificate
// It returns any error that may occure, or nil if the opperation was successful
func (s *SimpleCloudFoundryWorker) Authenticate(cert CloudFoundryCertificate) error {
	return s.cli.Auth(cert.APIEndPoint, cert.Username, cert.Password).Sync()
}

// NewCloudFoundryWorker creates a new instance of the SimpleCloudFoundryWorker struct
func NewCloudFoundryWorker(logger logger.Logger, cli CloudFoundryCLI) SimpleCloudFoundryWorker {
	return SimpleCloudFoundryWorker{
		logger: logger,
		cli:    cli,
	}
}

// CloudFoundryTask is a simply configuration like struct that contains a specific task
type CloudFoundryTask struct {
	Command    string
	Parameters []string
}
