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
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"sync"
)

// AppProvider defines an object that provides one single pushed app in the cf
//
// Push pushes the app if not pushed yet
//
// AppName returns the app name
type AppProvider interface {
	Push(l logger.Logger) (infoLog logger.CachedLogger, errorLog logger.CachedLogger, err error)
	ForcePush(l logger.Logger, appName string) (infoLog logger.CachedLogger, errorLog logger.CachedLogger, err error)
	AppName() string
}

// MutexAppProvider is a mutex.lock based implementation of the app provider
type MutexAppProvider struct {
	AppNameValue string
	Pushed       bool
	Lock         *sync.Mutex
	AssetPath    string
	Cli          cfw.CloudFoundryCLI
}

// AppName returns the app name
func (m *MutexAppProvider) AppName() string {
	return m.AppNameValue
}

// NewMutexSingleAppProvider creates a new single app provider
func NewMutexSingleAppProvider(cli cfw.CloudFoundryCLI, appName string, assetPath string) *MutexAppProvider {
	return &MutexAppProvider{
		AppNameValue: appName,
		Pushed:       false,
		Lock:         &sync.Mutex{},
		AssetPath:    assetPath,
		Cli:          cli,
	}
}

// Push pushes the app if not pushed yet. This will lock any other setup app
func (m *MutexAppProvider) Push(l logger.Logger) (infoLog logger.CachedLogger, errorLog logger.CachedLogger, err error) {
	defer m.Lock.Unlock()

	m.Lock.Lock()
	if m.Pushed {
		return nil, nil, nil
	}

	m.Pushed = true
	return m.ForcePush(l, m.AppNameValue)
}

// ForcePush force pushes the app to the
func (m *MutexAppProvider) ForcePush(l logger.Logger, appName string) (infoLog logger.CachedLogger, errorLog logger.CachedLogger, err error) {
	infoLogger := logger.NewByteBufferCachedLogger(l.ReportingOn(logger.Error))
	errorLogger := logger.NewByteBufferCachedLogger(l.ReportingOn(logger.Error))

	if err := m.Cli.Push(m.AssetPath, appName, 1).SubscribeOnErr(errorLogger).SubscribeOnOut(infoLogger).Sync(); err != nil {
		return infoLogger, errorLogger, err
	}

	infoLogger.Clear()
	errorLogger.Clear()
	return nil, nil, nil
}
