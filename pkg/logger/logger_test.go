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

package logger_test

import (
	. "github.com/homeport/disrupt-o-meter/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger code test", func() {
	Context("Testing default behaviour of logger instances", func() {
		var (
			logger          Logger
			channelProvider ChannelProvider
		)

		BeforeEach(func() {
			channelProvider = NewChannelProvider(1)
			logger = NewChanneledLogger("test-logger", channelProvider)
		})

		It("Should have the correct name", func() {
			Expect(logger.GetName()).To(BeEquivalentTo("test-logger"))
		})

		It("Should forward messages correctly", func() {
			logger.WriteString("test")
			logger.WriteString("more-tests")
			Expect(channelProvider.Read().GetMessageAsString()).To(BeEquivalentTo("test"))
			Expect(channelProvider.Read().GetMessageAsString()).To(BeEquivalentTo("more-tests"))
		})
	})
})
