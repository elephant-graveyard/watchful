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

package cmd

import (
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/spf13/cobra"
	"os"
	"time"
)

// runCmd is the run command definition using cobra
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run is the sub-command that executes watchful",
	Long:  "Run executes watchful with the provided configuration and merkhet targets.",
	Run:   run,
}

// run is the method called when someone uses the watchful run command
func run(cmd *cobra.Command, args []string) {
	loggerChannelProvider := logger.NewChannelProvider(10)
	loggerFactory := logger.NewChanneledLoggerFactory(loggerChannelProvider)

	cloudFoundryLogger := loggerFactory.NewChanneledLogger("cf-cli-worker task")
	watchfulLogger := loggerFactory.NewChanneledLogger("watchful")

	loggerConfig := logger.NewGroupContainer().NewGroup(cloudFoundryLogger).NewGroup(watchfulLogger)

	loggerClusterConfig := logger.NewBasicSplitPipelineConfig(true, time.Local, loggerConfig)

	loggerCluster := logger.NewLoggerCluster(logger.NewSplitPipeline(loggerClusterConfig, os.Stdout),
		loggerChannelProvider, time.Second)
	go loggerCluster.StartListening()

	worker := cfw.NewCloudFoundryWorker(cloudFoundryLogger, cfw.NewBashCloudFoundryCLI())

	watchfulLogger.WriteString(logger.Info, "Targeting API Endpoint")
	if err := worker.API("https://api.lynx-cluster.eu-de.containers.appdomain.cloud"); err != nil {
		watchfulLogger.WriteString(logger.Error, "Could not target API endpoint "+err.Error())
		goto teardown
	}

	watchfulLogger.WriteString(logger.Info, "Authenticating against API endpoint")
	if err := worker.Authenticate(cfw.CloudFoundryCertificate{
		Username: "apikey",
		Password: "XIBRuNEeMeYMmHkIkjVMmUVv_WpzvDv-UWmSim9frpi6",
	}); err != nil {
		watchfulLogger.WriteString(logger.Error, "Could not authenticate against API endpoint "+err.Error())
		goto teardown
	}

	watchfulLogger.WriteString(logger.Info, "Creating test environment")
	if err := worker.CreateTestEnvironment(); err != nil {
		watchfulLogger.WriteString(logger.Error, "Could not create test environment")
		goto teardown
	}

	watchfulLogger.WriteString(logger.Info, "Tearing down test environment")
	if err := worker.TeardownTestEnvironment(); err != nil {
		watchfulLogger.WriteString(logger.Error, "Could not teardown test environment")
		goto teardown
	}

teardown:
	watchfulLogger.WriteString(logger.Info, "Done ! Shutting down..")
	loggerChannelProvider.Close()
	loggerCluster.WaitGroup().Wait()
}

// init adds the runCmd to the watchful root command
func init() {
	rootCmd.AddCommand(runCmd)
}
