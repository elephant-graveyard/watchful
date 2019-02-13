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
	"strings"
)

// CloudFoundryCLI is a utility object that can create cli commands for the cloud foundry cli
//
// API targets a specific api endpoint within the cf cli
//
// CreateOrganization creates a new organization on the cloud foundry instance
//
// DeleteOrganization deletes a organization on the cloud foundry instance
//
// CreateSpace creates a new space on the cloud foundry instance
//
// DeleteSpace deletes a space on the cloud foundry instance
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
	API(apiEndpoint string, validateSSL bool) CommandPromise
	CreateOrganization(name string) CommandPromise
	DeleteOrganization(name string) CommandPromise
	CreateSpace(org string, name string) CommandPromise
	DeleteSpace(org string, name string) CommandPromise
	Auth(username string, password string) CommandPromise
	Target(organization string, space string) CommandPromise
	Push(path string, name string, instances int) CommandPromise
	Delete(name string) CommandPromise
	Scale(name string, instances int) CommandPromise
}

// BashCloudFoundryCLI is a cli runner that relies on executing bash commands.
// This struct is empty, but will be used as the CloudFoundryCLI interface has an important role
// as it could also be implemented by a CloudFoundryAPI call implementation
type BashCloudFoundryCLI struct {
}

// API targets a specific api endpoint within the cf cli
func (b *BashCloudFoundryCLI) API(apiEndpoint string, validateSSL bool) CommandPromise {
	if !validateSSL {
		return createCFCommandPromise(fmt.Sprintf("api --skip-ssl-validation %s", apiEndpoint))
	}
	return createCFCommandPromise(fmt.Sprintf("api %s", apiEndpoint))
}

// CreateOrganization creates a new organization on the cloud foundry instance
func (b *BashCloudFoundryCLI) CreateOrganization(name string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("create-org %s", name))
}

// DeleteOrganization deletes a organization on the cloud foundry instance
func (b *BashCloudFoundryCLI) DeleteOrganization(name string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("delete-org -f %s", name))
}

// CreateSpace creates a new space on the cloud foundry instance
func (b *BashCloudFoundryCLI) CreateSpace(org string, name string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("create-space -o %s %s", org, name))
}

// DeleteSpace deletes a space on the cloud foundry instance
func (b *BashCloudFoundryCLI) DeleteSpace(org string, name string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("delete-space -o %s -f %s", org, name))
}

// Auth creates a CommandPromise that will try to authenticate against the cloud foundry instance
func (b *BashCloudFoundryCLI) Auth(username string, password string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("auth %s %s", username, password))
}

// Target targets the given organization instance
func (b *BashCloudFoundryCLI) Target(organization string, space string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("target -o %s -s %s", organization, space))
}

// Push pushes a new instance to the cloud foundry instance. The instance has to be under the provided path
// It will be pushed with the provided name and n instances will be created
func (b *BashCloudFoundryCLI) Push(path string, name string, instances int) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("push %s -i %d -p %s", name, instances, path))
}

// Delete will delete the provided app from the selected space
func (b *BashCloudFoundryCLI) Delete(name string) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("delete %s -r -f", name))
}

// Scale will scale the app instance to the provided amount
func (b *BashCloudFoundryCLI) Scale(name string, instances int) CommandPromise {
	return createCFCommandPromise(fmt.Sprintf("scale %s -i %d", name, instances))
}

// SplitParameterString splits the string into a slice of parameters
func SplitParameterString(parameter string) []string {
	return strings.Split(parameter, " ")
}

// NewBashCloudFoundryCLI creates a new bash based cloud foundry cli
func NewBashCloudFoundryCLI() *BashCloudFoundryCLI {
	return &BashCloudFoundryCLI{}
}

// createCFCommandPromise simply creates a new command promise executing cf plus the parameter list
func createCFCommandPromise(parameter string) CommandPromise {
	return NewSimpleCommandPromise(exec.Command("cf", SplitParameterString(parameter)...))
}
