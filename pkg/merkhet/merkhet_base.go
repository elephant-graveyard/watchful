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

package merkhet

import (
	"github.com/homeport/watchful/pkg/logger"
	"sync"
)

// Base represents a merkhet base which contains the values every merkhet needs
//
// Logger returns the logger instance the merkhet uses
//
// Configuration returns the configuration instance the merkhet uses
//
// RecordSuccessfulRun records one new successful run
//
// RecordFailedRuns records one new failed run
//
// NewResultSet builds a new result set
type Base interface {
	Logger() logger.Logger
	Configuration() Configuration
	RecordSuccessfulRun()
	RecordFailedRun()
	NewResultSet() Result
}

// SimpleBase is a basic variable driven implementation of the Base interface
type SimpleBase struct {
	LoggerReference        logger.Logger
	ConfigurationReference Configuration
	SuccessfulRuns         int
	FailedRun              int
	Lock                   *sync.Mutex
}

// NewSimpleBase creates a new basic simple base
func NewSimpleBase(loggerReference logger.Logger, configurationReference Configuration) *SimpleBase {
	return NewSetSimpleBase(loggerReference, configurationReference, 0, 0)
}

// NewSetSimpleBase creates a new simple base with setting failed and successful runs
func NewSetSimpleBase(loggerReference logger.Logger, configurationReference Configuration, successfulRuns int, failedRun int) *SimpleBase {
	return &SimpleBase{
		LoggerReference:        loggerReference,
		ConfigurationReference: configurationReference,
		SuccessfulRuns:         successfulRuns,
		FailedRun:              failedRun,
		Lock:                   &sync.Mutex{},
	}
}

// Logger returns the logger reference
func (b *SimpleBase) Logger() logger.Logger {
	return b.LoggerReference
}

// Configuration returns the configuration reference of the base
func (b *SimpleBase) Configuration() Configuration {
	return b.ConfigurationReference
}

// RecordSuccessfulRun records a successful run
func (b *SimpleBase) RecordSuccessfulRun() {
	defer b.Lock.Unlock()

	b.Lock.Lock()
	b.SuccessfulRuns++
}

// RecordFailedRun records a failed run
func (b *SimpleBase) RecordFailedRun() {
	defer b.Lock.Unlock()

	b.Lock.Lock()
	b.FailedRun++
}

// NewResultSet builds a new result set instance
func (b *SimpleBase) NewResultSet() Result {
	totalRuns := b.FailedRun + b.SuccessfulRuns
	return NewMerkhetResult(totalRuns, b.FailedRun, b.Configuration().ValidRun(totalRuns, b.FailedRun))
}
