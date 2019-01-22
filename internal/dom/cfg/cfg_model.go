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

// DOMConfig is a structure defining the configuration of the disrupt-o-meter
type DOMConfig struct {
	CloudFoundryConfig    CloudFoundryConfig     `json:"cf"`
	TaskConfigurations    []TaskConfiguration    `json:"tasks"`
	MerkhetConfigurations []MerkhetConfiguration `json:"merkhets"`
	LoggerConfiguration   LoggerConfiguration    `json:"logger-config"`
}

// CloudFoundryConfig contains the config to connect and communicate with the cloud foundry instance
type CloudFoundryConfig struct {
	Domain              string   `json:"domain"`
	APIEndPoint         string   `json:"api-endpoint"`
	SkipSSLValidation   bool     `json:"skip-ssl-validation"`
	CustomCLIParameters []string `json:"custom-cli-parameters"`
	Username            string   `json:"username"`
	Password            string   `json:"password"`
}

// TaskConfiguration is the configuration for a simply task that is executed against the cloud foundry instance
type TaskConfiguration struct {
	Executable string   `json:"cmd"`
	Parameters []string `json:"args"`
}

// MerkhetConfiguration is the configuration of one merkhet instance running
type MerkhetConfiguration struct {
	Name       string `json:"name"`
	Threshhold string `json:"threshhold"`
}

// LoggerConfiguration is the config for the logger system dom uses
type LoggerConfiguration struct {
	TimeLocation    string `json:"time-location"`
	PrintLoggerName bool   `json:"print-logger-name"`
}
