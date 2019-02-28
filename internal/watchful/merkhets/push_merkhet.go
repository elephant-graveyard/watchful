// Copyright (c) 2019 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package merkhets

import (
	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
)

// PushMerkhet is a merkhet implementation that pushes a cf app to the instance
type PushMerkhet struct {
	BaseReference   merkhet.Base
	CloudFoundryCLI cfw.CloudFoundryCLI
	AppProvider     AppProvider
}

// NewPushMerkhet creates a new push merkhet
func NewPushMerkhet(baseReference merkhet.Base, appProvider AppProvider, cloudFoundryCLI cfw.CloudFoundryCLI) *PushMerkhet {
	return &PushMerkhet{BaseReference: baseReference, AppProvider: appProvider, CloudFoundryCLI: cloudFoundryCLI}
}

// Install does nothing on the push merkhet
func (m *PushMerkhet) Install() error {
	return nil
}

// PostConnect does nothing
func (m *PushMerkhet) PostConnect() error {
	return nil
}

// Execute tries to push an app to the cloud foundry
func (m *PushMerkhet) Execute() error {
	infoLog, errorLog, err := m.AppProvider.ForcePush(m.Base().Logger(), m.Base().Configuration().Name())
	if err != nil {
		m.Base().Logger().WriteString(logger.Error, "Could not push app to cf instance")

		infoLog.FlushTarget().Refocus(logger.Debug)
		infoLog.FlushTarget()

		errorLog.FlushTarget().Refocus(logger.Debug)
		errorLog.FlushTarget()
		return err
	}

	m.Base().Logger().WriteString(logger.Debug, bunt.Sprintf("SpringGreen{Pushed app to cf instance}"))
	return nil
}

// Base returns the base instance of the merkhet
func (m *PushMerkhet) Base() merkhet.Base {
	return m.BaseReference
}
