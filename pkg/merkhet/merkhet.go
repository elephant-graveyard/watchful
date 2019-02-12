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
	"github.com/homeport/disrupt-o-meter/pkg/logger"
)

// Merkhet defines a runnable measurement task that can be executed during the Cloud Foundry maintenance
//
// Install installs the merkhet instance. This method call will be used to setup necessary dependencies of the merkhet
//
// PostConnect is called after DOM successfully authenticated against the cloud foundry instance
//
// PostConnect is called after DOM successfully authenticated against the cloud foundry instance
//
// BuildResult creates a new Result instance containing the
//
// Configuration returns the configuration instances the Merkhet is depending on
//
// Configuration returns the configuration instances the Merkhet is depending on
//
// RecordSuccessfulRun records one successful run
//
// RecordFailedRun records one failed run
type Merkhet interface {
	Install() error
	PostConnect() error
	Execute() error
	BuildResult() Result
	Configuration() Configuration
	Logger() logger.Logger

	RecordSuccessfulRun()
	RecordFailedRun()
}

// Configuration contains the passed configuration values for a Merkhet instance
//
// Name returns the name provided in the configuration.
//
// ValidRun returns whether the provided total relative to the fails is still considered a viable run
// The behaviour of this method is heavily reliant on the implementation
type Configuration interface {
	Name() string
	ValidRun(totalRuns int, failedRuns int) bool
}

// namedConfiguration is a simple structure containing the name of a MerkhetConfiguration
type namedConfiguration struct {
	name string
}

// Name returns the name stored in the named configuration
func (n *namedConfiguration) Name() string {
	return n.name
}

// PercentageConfiguration is a configuration implementation based on a percentage threshold
type PercentageConfiguration struct {
	namedConfiguration  *namedConfiguration
	percentageThreshold float64
}

// Name returns the name stored in the configuration delegate
func (p *PercentageConfiguration) Name() string {
	return p.namedConfiguration.Name()
}

// ValidRun returns if the failed runs compared to the total runs are below the provided percentage threshold
func (p *PercentageConfiguration) ValidRun(totalRuns int, failedRuns int) bool {
	return (float64(failedRuns) / float64(totalRuns)) <= p.percentageThreshold
}

// FlatConfiguration is an implementation of the Configuration interface that is based on a flat amount of failed runs
// to calculate viability
type FlatConfiguration struct {
	namedConfiguration *namedConfiguration
	flatThreshold      int
}

// Name returns the name stored in the configuration delegate
func (f *FlatConfiguration) Name() string {
	return f.namedConfiguration.Name()
}

// ValidRun returns if the failed runs compared to the total runs are below the provided percentage threshold
func (f *FlatConfiguration) ValidRun(totalRuns int, failedRuns int) bool {
	return failedRuns <= f.flatThreshold
}

// NewPercentageConfiguration creates a new configuration instance that uses a percentage threshold
func NewPercentageConfiguration(name string, percentageThreshold float64) *PercentageConfiguration {
	return &PercentageConfiguration{
		namedConfiguration: &namedConfiguration{
			name: name,
		},
		percentageThreshold: percentageThreshold,
	}
}

// NewFlatConfiguration creates a new configuration instance that uses a flat threshold
func NewFlatConfiguration(name string, flatThreshold int) *FlatConfiguration {
	return &FlatConfiguration{
		namedConfiguration: &namedConfiguration{
			name: name,
		},
		flatThreshold: flatThreshold,
	}
}
