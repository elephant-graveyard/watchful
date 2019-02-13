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

package executable

import (
	"github.com/homeport/watchful/internal/watchful/cfg"
	"github.com/homeport/watchful/internal/watchful/merkhets"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
	"regexp"
	"strconv"
)

var (
	// PercentageThresholdRegex defines the regex that identifies a percentage value
	PercentageThresholdRegex = regexp.MustCompile("([0-9]*\\.)?[0-9]*%")
)

// MerkhetExecutable defines the executable responsible for controlling the merkhets
type MerkhetExecutable struct {
	Configuration cfg.WatchfulConfig
	Pool          merkhet.Pool
	LoggerFactory logger.Factory
	LoggerGroup   logger.Group
}

// NewMerkhetExecutable creates a new merkhet executable service
func NewMerkhetExecutable(Configuration cfg.WatchfulConfig, LoggerFactory logger.Factory, LoggerGroup logger.Group) *MerkhetExecutable {
	return &MerkhetExecutable{
		Configuration: Configuration,
		Pool:          merkhet.NewPool(),
		LoggerGroup:   LoggerGroup,
		LoggerFactory: LoggerFactory,
	}
}

// Execute executes the given service
func (e *MerkhetExecutable) Execute() error {
	for _, c := range e.Configuration.MerkhetConfigurations {
		base, err := e.createMerkhetBase(c)
		if err != nil {
			return err
		}

		switch c.Name {
		case "http-availability":
			e.Pool.StartWorker(merkhets.NewDefaultCurlMerkhet(e.Configuration.CloudFoundryConfig.Domain, base), *c.HeartbeatRate, defaultHeartbeatHandler())
		}
	}

	return nil
}

// createMerkhetBase creates a new merkhet base instance
func (e *MerkhetExecutable) createMerkhetBase(configuration cfg.MerkhetConfiguration) (base merkhet.Base, err error) {
	merkhetLogger := e.LoggerFactory.NewChanneledLogger(configuration.Name)
	e.LoggerGroup.Add(merkhetLogger)

	var merkhetConfig merkhet.Configuration
	if PercentageThresholdRegex.Match([]byte(configuration.Threshold)) {
		pureString := configuration.Threshold[len(configuration.Threshold)-2:]
		if f, err := strconv.ParseFloat(pureString, 64); err == nil {
			merkhetConfig = merkhet.NewPercentageConfiguration(configuration.Name, f)
		} else {
			return nil ,err
		}
	} else {
		if i, err := strconv.ParseInt(configuration.Threshold, 0, 32); err == nil {
			merkhetConfig = merkhet.NewFlatConfiguration(configuration.Name, int(i))
		} else {
			return nil, err
		}
	}

	return merkhet.NewSimpleBase(merkhetLogger, merkhetConfig) , nil
}

// defaultHeartbeatHandler creates a new default consumer
func defaultHeartbeatHandler() merkhet.Consumer {
	return merkhet.ConsumeAsync(func(m merkhet.Merkhet, relay merkhet.ControllerChannel) {
		e := m.Execute()
		relay <- merkhet.ConsumeSync(func(merkhet merkhet.Merkhet, relay merkhet.ControllerChannel) {
			if e == nil {
				merkhet.Base().RecordSuccessfulRun()
			} else {
				merkhet.Base().RecordSuccessfulRun()
			}
		})
	})
}
