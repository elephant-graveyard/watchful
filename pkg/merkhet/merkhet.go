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

//Merkhet deinfes a runnable mesaurement task that can be executed during the Cloud Foundery maintanance
type Merkhet interface {

	//Install installs the merkhet instance. This method call will be used to setup necessary dependencies of the merkhet
	Install()

	//PostConnect is called after DOM successfully authenticated against the cloud foundery instance
	PostConnect()

	//Execute is called each time the merkhet should test against the cloud foundery service
	Execute()

	//BuildResult creates a new Result instance containing the
	BuildResult() Result

	//GetConfiguration returns the configuration instances the Merkhet is depending on
	GetConfiguration() Configuration

	//GetLogger returns the logger the merkhet instance is redirecting it's output to
	GetLogger() logger.Logger
}

//--

//Configuration contains the passed configuration values for a Merkhet instance
type Configuration interface {
	GetName() string
	IsValidRun(totalRuns uint, failedRuns uint) bool
}

//namedConfiguration is a simple structure containing the name of a MerhetConfiguration
type namedConfiguration struct {
	Name string
}

//GetName returns the name stored in the named configuration
func (n *namedConfiguration) GetName() string {
	return n.Name
}

//--

//percentageConfiguration is a configuration implementation based on a percentage threshhold
type percentageConfiguration struct {
	NamedConfiguration   *namedConfiguration
	PercentageThreshhold float64
}

//GetName returns the name stored in the configuration delegate
func (p *percentageConfiguration) GetName() string {
	return p.NamedConfiguration.GetName()
}

//IsValidRun returns if the failed runs comapred to the total runs are below the provided percentage threshhold
func (p *percentageConfiguration) IsValidRun(totalRuns uint, failedRuns uint) bool {
	return (float64(failedRuns) / float64(totalRuns)) <= p.PercentageThreshhold
}

//--

//flatConfiguration is an implementation of the Configuration interface that is based on a flat amount of failed runs
//to calucalte viability
type flatConfiguration struct {
	NamedConfiguration *namedConfiguration
	FlatThreshhold     uint
}

//GetName returns the name stored in the configuration delegate
func (f *flatConfiguration) GetName() string {
	return f.NamedConfiguration.GetName()
}

//IsValidRun returns if the failed runs comapred to the total runs are below the provided percentage threshhold
func (f *flatConfiguration) IsValidRun(totalRuns uint, failedRuns uint) bool {
	return failedRuns <= f.FlatThreshhold
}

//--

//NewPercentageConfiguration creates a new configuration intaces that uses a percentage threshold
func NewPercentageConfiguration(name string, percentageTreshhold float64) Configuration {
	return &percentageConfiguration{
		NamedConfiguration: &namedConfiguration{
			Name: name,
		},
		PercentageThreshhold: percentageTreshhold,
	}
}

//NewFlatConfiguration creates a new configuration intaces that uses a flat threshold
func NewFlatConfiguration(name string, flatThreshhold uint) Configuration {
	return &flatConfiguration{
		NamedConfiguration: &namedConfiguration{
			Name: name,
		},
		FlatThreshhold: flatThreshhold,
	}
}
