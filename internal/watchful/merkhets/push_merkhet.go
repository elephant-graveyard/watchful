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
	"github.com/pkg/errors"
)

// PushMerkhet is a merkhet implementation that pushes a cf app to the instance
type PushMerkhet struct {
	BaseReference   merkhet.Base
	AssetPath       string
	CloudFoundryCLI cfw.CloudFoundryCLI
}

// NewPushMerkhet creates a new push merkhet
func NewPushMerkhet(baseReference merkhet.Base, assetPath string, cloudFoundryCLI cfw.CloudFoundryCLI) *PushMerkhet {
	return &PushMerkhet{BaseReference: baseReference, AssetPath: assetPath, CloudFoundryCLI: cloudFoundryCLI}
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
	cachedLogger := logger.NewByteBufferCachedLogger(m.Base().Logger().ReportingOn(logger.Debug))

	if err := m.CloudFoundryCLI.Push(m.AssetPath, m.Base().Configuration().Name(), 1).
		SubscribeOnErr(m.Base().Logger().ReportingOn(logger.Error)).
		SubscribeOnOut(cachedLogger).Sync(); err != nil {

		m.Base().Logger().WriteString(logger.Error, bunt.Sprintf("Red{Could not push app to cf instance}"))
		cachedLogger.Flush() // Flush to debug channel if failed
		return errors.Wrap(err, "could not push sample app")
	}
	cachedLogger.Clear()
	m.Base().Logger().WriteString(logger.Info, bunt.Sprintf("SpringGreen{Pushed merkhet app successfully}"))
	return nil
}

// Base returns the base instance of the merkhet
func (m *PushMerkhet) Base() merkhet.Base {
	return m.BaseReference
}
