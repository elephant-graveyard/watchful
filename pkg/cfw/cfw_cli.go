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
	"strings"
)

// CloudFoundryCLI is a utility object that can create cli commands for the cloud foundry cli
//
// Auth creates a CommandPromise that will try to authenticate against the cloud foundry instance
type CloudFoundryCLI interface {
	Auth(endpoint string, username string, password string) *CommandPromise
}

// BashCloudFoundryCLI is a cli runner that relies on executing bash commands.
// This struct is empty, but will be used as the CloudFoundryCLI interface has an important role
// as it could also be implemented by a CloudFoundryAPI call implementation
type BashCloudFoundryCLI struct {
}

// Auth creates a CommandPromise that will try to authenticate against the cloud foundry instance
func (b *BashCloudFoundryCLI) Auth(endpoint string, username string, password string) *CommandPromise {
	return &CommandPromise{
		command: exec.Command("cf", strings.Split(" ", "login -c ")...), //TODO CONTINUE
	}
}
