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

package cfg_test

import (
	"github.com/homeport/watchful/internal/watchful/cfg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration testing", func() {
	Context("Tests the parsing of the config model from string file yaml and json", func() {

		var (
			config cfg.WatchfulConfig
		)

		BeforeEach(func() {
			config = cfg.WatchfulConfig{}
		})

		It("Should read the config from the model yaml file", func() {
			Expect(cfg.ParseFromFile("./test_config.yml", &config)).To(BeNil())

			Expect(config.CloudFoundryConfig.SkipSSLValidation).To(BeTrue())
			Expect(config.CloudFoundryConfig.Domain).To(BeEquivalentTo("file-yaml-test.com"))
		})

		It("Should read the config from a string yaml", func() {
			var body = `---
cf:
  domain: string-yaml-test.com
  api-endpoint: a.test.com
  skip-ssl-validation: true
  custom-cli-parameters:
  - --test
  username: test_user
  password: test_password
logger-config:
  time-location: UTC
  print-logger-name: true`

			Expect(cfg.ParseFromString(body, &config)).To(BeNil())

			Expect(config.CloudFoundryConfig.SkipSSLValidation).To(BeTrue())
			Expect(config.CloudFoundryConfig.Domain).To(BeEquivalentTo("string-yaml-test.com"))
			Expect(config.LoggerConfiguration.PrintLoggerName).To(BeTrue())
		})

		It("Should read the config from the model json file", func() {
			Expect(cfg.ParseFromFile("./test_config.json", &config)).To(BeNil())

			Expect(config.CloudFoundryConfig.SkipSSLValidation).To(BeTrue())
			Expect(config.CloudFoundryConfig.Domain).To(BeEquivalentTo("file-json-test.com"))
		})

		It("Should read the config from a string json", func() {
			var body = `{"cf":{"domain":"string-json-test.com","api-endpoint":"a.test.com","skip-ssl-validation":true,"custom-cli-parameters":["--test"],"username":"testuser","password":"test_password"}}`

			Expect(cfg.ParseFromString(body, &config)).To(BeNil())

			Expect(config.CloudFoundryConfig.SkipSSLValidation).To(BeTrue())
			Expect(config.CloudFoundryConfig.Domain).To(BeEquivalentTo("string-json-test.com"))
		})
	})
})
