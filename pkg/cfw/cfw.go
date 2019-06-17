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
	"os/exec"

	"github.com/homeport/watchful/pkg/logger"
)

const (
	// WatchfulOrgName is the name of the organization watchful will create
	WatchfulOrgName = "watchful"

	// WatchfulSpaceName is the name of the space watchful will create
	WatchfulSpaceName = "watchful"
)

// CloudFoundryWorker defines a worker object that is capable of executing specific commands
// against a cloud foundry instance. It will do this in the go routine the method call was issued in,
// so make sure to spin up a new routine if you don't want to block your main process
//
// API targets a specific api endpoint on the cli
//
// Authenticate authenticates the worker against the cloud foundry instance provided in the certificate
// It returns any error that may occur, or nil if the operation was successful
//
// CreateTestEnvironment creates a test env in the cloud foundry instance
//
// TeardownTestEnvironment tears down existing test environments
//
// Logger returns the logger instance the worker task is using
//
// Execute executes the provided task against the cloud foundry instance it is working on
// If the worker is not authenticated it will error, else it will try to execute the
// task on the same go routine it was called in
type CloudFoundryWorker interface {
	API(apiEndpoint string, validateSSL bool) error
	Authenticate(cert CloudFoundryCertificate) error
	CreateTestEnvironment() error
	TeardownTestEnvironment() error
	Logger() logger.Logger
	Execute(task CloudFoundryTask) error
}

// SimpleCloudFoundryWorker is a basic implementation of the CloudFoundryWorker interface
type SimpleCloudFoundryWorker struct {
	logger logger.Logger
	cli    CloudFoundryCLI
}

// API targets a specific api endpoint on the cli
func (s *SimpleCloudFoundryWorker) API(apiEndpoint string, validateSSL bool) error {
	return s.Wrap(s.cli.API(apiEndpoint, validateSSL)).Sync()
}

// Authenticate authenticates the worker against the cloud foundry instance provided in the certificate
// It returns any error that may occur, or nil if the operation was successful
func (s *SimpleCloudFoundryWorker) Authenticate(cert CloudFoundryCertificate) error {
	return s.Wrap(s.cli.Auth(cert.Username, cert.Password)).Sync()
}

// CreateTestEnvironment creates a test env in the cloud foundry instance
func (s *SimpleCloudFoundryWorker) CreateTestEnvironment() error {
	if err := s.Wrap(s.cli.CreateOrganization(WatchfulOrgName)).Sync(); err != nil {
		return err
	}

	if err := s.Wrap(s.cli.CreateSpace(WatchfulOrgName, WatchfulSpaceName)).Sync(); err != nil {
		return err
	}

	return s.Wrap(s.cli.Target(WatchfulOrgName, WatchfulSpaceName)).Sync()
}

// TeardownTestEnvironment tears down existing test environments
func (s *SimpleCloudFoundryWorker) TeardownTestEnvironment() error {
	return s.Wrap(s.cli.DeleteOrganization(WatchfulOrgName)).Sync()
}

// Logger returns the logger instance the worker task is using
func (s *SimpleCloudFoundryWorker) Logger() logger.Logger {
	return s.logger
}

// Execute executes the provided task against the cloud foundry instance it is working on
// If the worker is not authenticated it will error, else it will try to execute the
// task on the same go routine it was called in
func (s *SimpleCloudFoundryWorker) Execute(task CloudFoundryTask) error {
	return s.Wrap(NewSimpleCommandPromise(exec.Command(task.Command, task.Parameters...))).Sync()
}

// Wrap will inject the workers logger into the promise
func (s *SimpleCloudFoundryWorker) Wrap(promise CommandPromise) CommandPromise {
	return promise.
		SubscribeOnOut(s.logger.ReportingOn(logger.Info)).
		SubscribeOnErr(s.logger.ReportingOn(logger.Error))
}

// NewCloudFoundryWorker creates a new instance of the SimpleCloudFoundryWorker struct
func NewCloudFoundryWorker(logger logger.Logger, cli CloudFoundryCLI) *SimpleCloudFoundryWorker {
	return &SimpleCloudFoundryWorker{
		logger: logger,
		cli:    cli,
	}
}

// CloudFoundryTask is a simply configuration like struct that contains a specific task
type CloudFoundryTask struct {
	Command    string
	Parameters []string
}
