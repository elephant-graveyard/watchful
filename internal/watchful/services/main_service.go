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
	"fmt"
	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"github.com/homeport/watchful/internal/watchful/cfg"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"time"
)

var (
	// ExportPath is the path assets are exported to
	ExportPath = "temp"
)

// Service is an interface that defines an object that is capable of being executed
//
// Execute executes the services
type Service interface {
	Execute() error
}

// Cleanable defines objects that need cleanup
type Cleanable interface {
	Cleanup()
}

// MainService defines the main services service of watchful
type MainService struct {
	TerminalWidth           int
	ConfigContent           string
	PushedAppSampleLanguage string
}

// Execute executes watchful with all outside parameters
func (e *MainService) Execute() error {
	config := &cfg.WatchfulConfig{}
	if len(e.ConfigContent) > 0 { // Load config
		if err := cfg.ParseFromString(e.ConfigContent, config); err != nil {
			return errors.Wrap(err, "could not parse cli provided config")
		}
	} else {
		if err := cfg.ParseFromFile("config.yml", config); err != nil {
			return errors.Wrap(err, "could not parse file based config")
		}
	}

	loggerChannelProvider := logger.NewChannelProvider(10)                   // Create logger channel provider
	loggerFactory := logger.NewChanneledLoggerFactory(loggerChannelProvider) // Create logger factory

	cloudFoundryLogger := loggerFactory.NewChanneledLogger("cf-cli-worker") // Create task logger
	taskLogger := loggerFactory.NewChanneledLogger("task-worker")           // Create task logger
	watchfulLogger := loggerFactory.NewChanneledLogger("watchful")          // Create watchful logger
	assetLogger := loggerFactory.NewChanneledLogger("asset-service")        // Create asset logger

	loggerConfig := logger.NewGroupContainer().
		NewGroup(cloudFoundryLogger, taskLogger).
		NewGroup(watchfulLogger, assetLogger)

	location, er := time.LoadLocation(config.LoggerConfiguration.TimeLocation) // Load the configured timezone
	if er != nil {
		return er
	}

	loggerClusterConfig := logger.NewSplitPipelineConfig(config.LoggerConfiguration.PrintLoggerName, location, e.TerminalWidth, loggerConfig) // Create cluster
	loggerCluster := logger.NewLoggerCluster(logger.NewSplitPipeline(loggerClusterConfig, os.Stdout), // Create pipeline
		loggerChannelProvider, time.Second)
	go loggerCluster.StartListening() // Start cluster

	assetService := NewAssetService(e.PushedAppSampleLanguage, assetLogger, ExportPath)
	if err := assetService.Execute(); err != nil {
		return err
	}

	watchfulLogger.WriteString(logger.Info, fmt.Sprintf("Using time location %s", location.String()))

	shutdownNotifier := make(chan os.Signal, 1) // We want to be able to kill it in the same routine
	signal.Notify(shutdownNotifier, os.Interrupt)
	signal.Notify(shutdownNotifier, os.Kill)

	worker := cfw.NewCloudFoundryWorker(cloudFoundryLogger, cfw.NewBashCloudFoundryCLI())

	merkhetCore := NewMerkhetService(config, loggerFactory, loggerConfig.GroupByLogger(watchfulLogger), ExportPath)
	if err := merkhetCore.Execute(); err != nil {
		return err
	}

	go func() {
		watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{Installing merkhets}"))
		watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{  ->}"))
		if err := merkhetCore.Pool.ForEach(merkhet.ConsumeAsync(func(m merkhet.Merkhet, future merkhet.Future) {
			future.Complete(m.Install())
			watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Installed Aqua{%s}", m.Base().Configuration().Name()))
		})).Wait().FirstError(); err != nil {
			shutdownNotifier <- &ErrorSignal{InnerError: errors.Wrap(err, "could not install merkhets")}
			return
		}

		NewSetupService(config.CloudFoundryConfig, watchfulLogger, taskLogger, worker).Execute() // Run setup logic
		taskWorker := NewCloudFoundryService(config.TaskConfigurations, cloudFoundryLogger)

		watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{Post-Connecting merkhets}"))
		watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{  ->}"))
		if err := merkhetCore.Pool.ForEach(merkhet.ConsumeAsync(func(m merkhet.Merkhet, future merkhet.Future) {
			future.Complete(m.PostConnect())
			watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Post-Connected Aqua{%s}", m.Base().Configuration().Name()))
		})).Wait().FirstError(); err != nil {
			shutdownNotifier <- &ErrorSignal{InnerError: errors.Wrap(err, "could not post-connect merkhets")}
			return
		}

		taskIndex := 0
		for currentTaskConfig := taskWorker.Next(); currentTaskConfig != nil; currentTaskConfig = taskWorker.Next() { // Loop over all tasks
			taskIndex++

			if len(currentTaskConfig.MerkhetWhitelist) > 0 {
				merkhetCore.ApplyWhitelist(currentTaskConfig.MerkhetWhitelist)
			} else if len(currentTaskConfig.MerkhetBlacklist) > 0 {
				merkhetCore.ApplyBlacklist(currentTaskConfig.MerkhetBlacklist)
			} else {
				merkhetCore.Pool.StartHeartbeats()
			}

			watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{Executing task} #%d", taskIndex))
			watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Aqua{Using merkhets: }")) // Print currently running merkhets
			for _, beat := range merkhetCore.Pool.BeatingHearts() {
				watchfulLogger.WriteString(logger.Info, bunt.Sprintf("Gray{ - } Aqua{%s}", beat.Worker().Merkhet().Base().Configuration().Name()))
			}

			if err := taskWorker.Execute(); err != nil {
				taskLogger.WriteString(logger.Error, err.Error())
				watchfulLogger.WriteString(logger.Error, bunt.Sprintf("Red{Task #%d failed}", taskIndex))
				shutdownNotifier <- &ErrorSignal{InnerError: errors.Wrap(err, fmt.Sprintf("Faliure in task #%d", taskIndex))} // Shutdown with the given error
				return
			}

			watchfulLogger.WriteString(logger.Info, bunt.Sprintf("DarkGreen{Finished task #%d}", taskIndex))
			_ = taskWorker.Pop()
		}

		shutdownNotifier <- &ErrorSignal{} // Kill silently
	}()

	output := <-shutdownNotifier
	merkhetCore.Pool.Shutdown()
	watchfulLogger.WriteString(logger.Info, bunt.Sprintf("DarkGreen{Shutdown merkhets}")) // stop merkhets

	NewTeardownService(watchfulLogger, worker).Execute() // Teardown cf env
	assetService.Cleanup()                               // Cleans the asset service

	watchfulLogger.WriteString(logger.Info, "Done ! Shutting down..")
	loggerChannelProvider.Close()
	loggerCluster.WaitGroup().Wait()

	switch signalType := output.(type) {
	case *ErrorSignal:
		return signalType.Error()
	default:
		return fmt.Errorf("received external system signal: %s", signalType.String())
	}
}

// ErrorSignal is a an implementation of the Signal interface that contains an error
type ErrorSignal struct {
	InnerError error
}

// String returns the error included in the signal
func (e *ErrorSignal) String() string {
	return fmt.Sprintf("error-signal: %s", e.InnerError.Error())
}

// Signal does literally nothing
func (e *ErrorSignal) Signal() {}

// Error returns the internal error
func (e *ErrorSignal) Error() error {
	return e.InnerError
}
