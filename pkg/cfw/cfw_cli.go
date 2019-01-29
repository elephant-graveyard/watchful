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
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// CloudFoundryCLI is a utility object that can create cli commands for the cloud foundry cli
//
// Auth creates a CommandPromise that will try to authenticate against the cloud foundry instance
//
// Target targets the given organization instance. This will fail if the cli is not authenticated
//
// Push pushes a new instance to the cloud foundry instance. The instance has to be under the provided path
// It will be pushed with the provided name and n instances will be created
//
// Delete will delete the provided app from the selected space
//
// Scale will scale the app instance to the provided amount
type CloudFoundryCLI interface {
	Auth(endpoint string, username string, password string) CommandPromise
	Target(organization string, space string) CommandPromise
	Push(path string, name string, instances int) CommandPromise
	Delete(name string) CommandPromise
	Scale(name string, instances int) CommandPromise
	StreamLogs(name string)
}

// BashCloudFoundryCLI is a cli runner that relies on executing bash commands.
// This struct is empty, but will be used as the CloudFoundryCLI interface has an important role
// as it could also be implemented by a CloudFoundryAPI call implementation
type BashCloudFoundryCLI struct {
}

// Auth creates a CommandPromise that will try to authenticate against the cloud foundry instance
func (b *BashCloudFoundryCLI) Auth(endpoint string, username string, password string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("login -a %s -u %s -p %s", endpoint, username, password))
}

// Target targets the given organization instance
func (b *BashCloudFoundryCLI) Target(organization string, space string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("target -o %s -s %s", organization, space))
}

// Push pushes a new instance to the cloud foundry instance. The instance has to be under the provided path
// It will be pushed with the provided name and n instances will be created
func (b *BashCloudFoundryCLI) Push(path string, name string, instances int) CommandPromise {
	return &SimpleCommandPromise{
		command: exec.Command("cf",
			"push", name,
			"-i", strconv.Itoa(instances),
			"-p", path),
	}
}

// Delete will delete the provided app from the selected space
func (b *BashCloudFoundryCLI) Delete(name string) CommandPromise {
	return &SimpleCommandPromise{
		command: exec.Command("cf",
			"delete", name,
			"-r",
			"-f"),
	}
}

// Scale will scale the app instance to the provided amount
func (b *BashCloudFoundryCLI) Scale(name string, instances int) CommandPromise {
	return &SimpleCommandPromise{
		command: exec.Command("cf",
			"scale", name,
			"-i", strconv.Itoa(instances)),
	}
}

// SplitParameterString splits the string into a slice of parameters
func SplitParameterString(parameter string) []string {
	return strings.Split(parameter, " ")
}

// createCFCommandPromise simply creates a new command promise executing cf plus the parameter list
func createCFCommandPromise(parameter string) CommandPromise {
	return &SimpleCommandPromise{
		command: exec.Command("cf", SplitParameterString(parameter)...),
	}
}
