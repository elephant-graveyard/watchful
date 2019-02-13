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

package cfg

import (
	"time"
)

// WatchfulConfig is a structure defining the configuration of the watchful project
type WatchfulConfig struct {
	CloudFoundryConfig    CloudFoundryConfig     `yaml:"cf"`
	TaskConfigurations    []TaskConfiguration    `yaml:"tasks"`
	MerkhetConfigurations []MerkhetConfiguration `yaml:"merkhets"`
	LoggerConfiguration   LoggerConfiguration    `yaml:"logger-config"`
}

// CloudFoundryConfig contains the config to connect and communicate with the cloud foundry instance
type CloudFoundryConfig struct {
	Domain              string   `yaml:"domain"`
	APIEndPoint         string   `yaml:"api-endpoint"`
	SkipSSLValidation   bool     `yaml:"skip-ssl-validation"`
	CustomCLIParameters []string `yaml:"custom-cli-parameters"`
	Username            string   `yaml:"username"`
	Password            string   `yaml:"password"`
}

// TaskConfiguration is the configuration for a simply task that is executed against the cloud foundry instance
type TaskConfiguration struct {
	Executable       string   `yaml:"cmd"`
	Parameters       []string `yaml:"args"`
	MerkhetWhitelist []string `yaml:"merkhet-whitelist"`
	MerkhetBlacklist []string `yaml:"merkhet-blacklist"`
}

// MerkhetConfiguration is the configuration of one merkhet instance running
type MerkhetConfiguration struct {
	Name          string         `yaml:"name"`
	Threshold     string         `yaml:"threshold"`
	HeartbeatRate *time.Duration `yaml:"heartbeat"`
}

// LoggerConfiguration is the config for the logger system watchful uses
type LoggerConfiguration struct {
	TimeLocation    string `yaml:"time-location"`
	PrintLoggerName bool   `yaml:"print-logger-name"`
}

// GetHeartbeatRate returns the rate in which the heart of the merkhet beats
func (m MerkhetConfiguration) GetHeartbeatRate(defaultValue time.Duration) time.Duration {
	if m.HeartbeatRate == nil {
		return defaultValue
	}
	return *m.HeartbeatRate
}
