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

package services

import (
	"github.com/homeport/watchful/internal/watchful/cfg"
	"github.com/homeport/watchful/internal/watchful/merkhets"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
	"regexp"
	"strconv"
	"time"
)

var (
	// PercentageThresholdRegex defines the regex that identifies a percentage value
	PercentageThresholdRegex = regexp.MustCompile("([0-9]*\\.)?[0-9]*%")
)

// MerkhetService defines the services responsible for controlling the merkhets
type MerkhetService struct {
	Configuration *cfg.WatchfulConfig
	Pool          merkhet.Pool
	LoggerFactory logger.Factory
	LoggerGroup   logger.Group
	AppProvider   merkhets.AppProvider
	Cli           cfw.CloudFoundryCLI
}

// NewMerkhetService creates a new merkhet services service
func NewMerkhetService(configuration *cfg.WatchfulConfig, loggerFactory logger.Factory, loggerGroup logger.Group, appProvider merkhets.AppProvider, cli cfw.CloudFoundryCLI) *MerkhetService {
	return &MerkhetService{
		Configuration: configuration,
		Pool:          merkhet.NewPool(),
		LoggerGroup:   loggerGroup,
		LoggerFactory: loggerFactory,
		AppProvider:   appProvider,
		Cli:           cli,
	}
}

// Execute executes the given service
func (e *MerkhetService) Execute() error {
	for _, c := range e.Configuration.MerkhetConfigurations {
		base, err := e.createMerkhetBase(c)
		if err != nil {
			return err
		}

		switch c.Name {
		case "http-availability":
			e.Pool.StartWorker(merkhets.NewDefaultCurlMerkhet(e.Configuration.CloudFoundryConfig.Domain, base,
				e.AppProvider),
				c.GetHeartbeatRate(time.Second), e.defaultHeartbeatHandler())
		case "app-pushability":
			e.Pool.StartWorker(merkhets.NewPushMerkhet(base, e.AppProvider,
				cfw.NewBashCloudFoundryCLI()), c.GetHeartbeatRate(time.Minute), e.defaultHeartbeatHandler())
		case "cf-recent-log-functionality":
			e.Pool.StartWorker(merkhets.NewLogRecentMerkhet(e.Cli, e.AppProvider, base),
				c.GetHeartbeatRate(10*time.Second), e.defaultHeartbeatHandler())
		case "cf-log-functionality":
			e.Pool.StartWorker(merkhets.NewLogStreamMerkhet(e.Cli, e.AppProvider, base),
				c.GetHeartbeatRate(30*time.Second), e.defaultHeartbeatHandler())
		}
	}

	return nil
}

// ApplyWhitelist starts every merkhet that is listed in the name list and stops the rest
func (e *MerkhetService) ApplyWhitelist(whitelisted []string) {
	for _, heart := range e.Pool.BeatingHearts() {
		isWhitelisted := contains(heart.Worker().Merkhet().Base().Configuration().Name(), whitelisted)
		if heart.IsBeating() && !isWhitelisted {
			heart.StopBeating()
		}
		if !heart.IsBeating() && isWhitelisted {
			heart.StartBeating()
		}
	}
}

// ApplyBlacklist stops every merkhet that is listed in the name list and starts the rest
func (e *MerkhetService) ApplyBlacklist(blacklisted []string) {
	for _, heart := range e.Pool.BeatingHearts() {
		isBlacklisted := contains(heart.Worker().Merkhet().Base().Configuration().Name(), blacklisted)
		if heart.IsBeating() && isBlacklisted {
			heart.StopBeating()
		}
		if !heart.IsBeating() && !isBlacklisted {
			heart.StartBeating()
		}
	}
}

// createMerkhetBase creates a new merkhet base instance
func (e *MerkhetService) createMerkhetBase(configuration cfg.MerkhetConfiguration) (base merkhet.Base, err error) {
	merkhetLogger := e.LoggerFactory.NewChanneledLogger(configuration.Name)
	e.LoggerGroup.Add(merkhetLogger)

	var merkhetConfig merkhet.Configuration
	if PercentageThresholdRegex.Match([]byte(configuration.Threshold)) {
		pureString := configuration.Threshold[:len(configuration.Threshold)-1]
		if f, err := strconv.ParseFloat(pureString, 64); err == nil {
			merkhetConfig = merkhet.NewPercentageConfiguration(configuration.Name, f/float64(100))
		} else {
			return nil, err
		}
	} else {
		if i, err := strconv.ParseInt(configuration.Threshold, 0, 32); err == nil {
			merkhetConfig = merkhet.NewFlatConfiguration(configuration.Name, int(i))
		} else {
			return nil, err
		}
	}

	return merkhet.NewSimpleBase(merkhetLogger, merkhetConfig), nil
}

// defaultHeartbeatHandler creates a new default consumer
func (e *MerkhetService) defaultHeartbeatHandler() merkhet.Consumer {
	return merkhet.ConsumeAsync(func(m merkhet.Merkhet, future merkhet.Future) {
		future.Complete(m.Execute())
		if _, err := future.IsCompleted(); err != nil {
			m.Base().RecordFailedRun()
		} else {
			m.Base().RecordSuccessfulRun()
		}
	}).Notify(e.Pool.TaskWaitGroup())
}

// contains checks if the string is contained in the string array
func contains(match string, array []string) bool {
	for _, s := range array {
		if s == match {
			return true
		}
	}
	return false
}
