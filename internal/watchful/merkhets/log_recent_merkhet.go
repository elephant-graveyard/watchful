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
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gonvenience/bunt"
	"github.com/homeport/watchful/pkg/cfw"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
)

var (
	// TimestampRegex defines the regex that retrieves a timestamp
	TimestampRegex = regexp.MustCompile(".*Timestamp{([0-9]*)}.*")
)

// LogRecentMerkhet is an implementation of the Merkhet interface that tests of the recent log fetch is successful
type LogRecentMerkhet struct {
	Cli           cfw.CloudFoundryCLI
	AppProvider   AppProvider
	BaseReference merkhet.Base
	lastTimeStamp int64
}

// NewLogRecentMerkhet creates a new instance of the merkhet implementation to check log recent
func NewLogRecentMerkhet(cli cfw.CloudFoundryCLI, appProvider AppProvider, baseReference merkhet.Base) *LogRecentMerkhet {
	return &LogRecentMerkhet{Cli: cli, AppProvider: appProvider, BaseReference: baseReference}
}

// Install installs the merkhet, in this case does nothing
func (m *LogRecentMerkhet) Install() error {
	return nil
}

// PostConnect post connects the merkhet, in this case pushes the sample app if not done
func (m *LogRecentMerkhet) PostConnect() error {
	infoLog, errorLog, err := m.AppProvider.Push(m.Base().Logger())
	if err != nil {
		m.Base().Logger().WriteString(logger.Error, "Could not post-connect upstream sample-app, printing logs")
		infoLog.Flush()
		errorLog.Flush()
		return err
	}

	m.Base().Logger().WriteString(logger.Info, "Post-Connected log-recent-merkhet")
	return nil
}

// Execute tests the recent logs
func (m *LogRecentMerkhet) Execute() error {
	errorLog := logger.NewByteBufferCachedLogger(m.Base().Logger().ReportingOn(logger.Debug))
	infoLog := &bytes.Buffer{}

	if err := m.Cli.RecentLogs(m.AppProvider.AppName()).SubscribeOnOut(infoLog).SubscribeOnErr(errorLog).Sync(); err != nil {
		m.Base().Logger().WriteString(logger.Error, "Could not fetch recent logs")
		errorLog.Flush()
		return err
	}

	foundTimestamps := TimestampRegex.FindAllStringSubmatch(infoLog.String(), -1)
	if len(foundTimestamps) < 1 {
		m.Base().Logger().WriteString(logger.Error, "Could not find timestamp in fetched logs")
		return fmt.Errorf("log did not contain timestamps")
	}

	timeStamp, e := strconv.ParseInt(foundTimestamps[len(foundTimestamps)-1][1], 0, 64)
	if e != nil {
		m.Base().Logger().WriteString(logger.Error, "Could not parse timestamp in fetched logs")
		return e
	}

	defer func() {
		m.lastTimeStamp = timeStamp
	}()

	if timeStamp <= m.lastTimeStamp {
		m.Base().Logger().WriteString(logger.Error, "Found timestamp is <= to previous one, no new logs")
		return fmt.Errorf("found timestamp is <= to previous one, no new logs")
	}

	m.Base().Logger().WriteString(logger.Debug, bunt.Sprintf("SpringGreen{Fetched recent logs successfully}"))
	return nil
}

// Base returns the base reference of the merkhet
func (m *LogRecentMerkhet) Base() merkhet.Base {
	return m.BaseReference
}
