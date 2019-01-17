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
			logger Logger
		)

		BeforeEach(func() {
			logger = NewMemStoredLogger("test-logger")
		})

		It("Should store logs correclty", func() {
			logger.WriteString(logger.GetName())
			Expect(string(logger.Peek())).To(BeEquivalentTo(logger.GetName()))
		})

		It("Should store multiple strings", func() {
			logger.WriteString("Hello")
			logger.WriteString(" ")
			logger.WriteString("World")

			Expect(string(logger.Peek())).To(BeEquivalentTo("Hello World"))
		})

		It("Should notify observers", func() {
			var notifiedString string
			logger.PushObserver(func(bytes []byte) {
				notifiedString = string(bytes)
			})

			logger.WriteString(logger.GetName())
			Expect(notifiedString).To(BeEquivalentTo(logger.GetName()))
		})

		It("Should clear its content correctly", func() {
			bytes := []byte(logger.GetName())
			logger.Write(bytes)

			Expect(logger.Clear()).To(BeEquivalentTo(bytes))
			Expect(len(logger.Peek())).To(BeZero())
		})

		It("Should have the correct name", func() {
			Expect(logger.GetName()).To(BeEquivalentTo("test-logger"))
		})
	})
})
