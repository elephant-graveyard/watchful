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
	"time"

	. "github.com/homeport/disrupt-o-meter/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logger code test", func() {
	Context("Testing default behaviour of logger instances", func() {
		var (
			channelProvider ChannelProvider
			pipeline        *PipelineMock
			cluster         Cluster
			loggerFactory   Factory
			logger          Logger
		)

		BeforeEach(func() {
			channelProvider = NewChannelProvider(10)
			pipeline = &PipelineMock{}
			cluster = NewLoggerCluster(pipeline, channelProvider, time.Second*1)
			loggerFactory = NewChanneledLoggerFactory(channelProvider)
			logger = loggerFactory.NewChanneledLogger("test-logger")
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

		It("Should cluster same logs differently", func(done Done) {
			pipeline.callback = func(i int, c []ChannelMessage) {
				if c[0].GetMessageAsString() == "first" {
					Expect(i).To(BeZero())
				}
				if c[0].GetMessageAsString() == "second" {
					Expect(i).To(BeEquivalentTo(1))
					close(done)
				}
			}

			go cluster.StartListening()

			logger.WriteString("first")
			logger.WriteString("second")
			close(channelProvider.GetChannel()) //Ends all communication as we wanna unit test on main thread
		})

		It("Should cluster different logs together", func(done Done) {
			pipeline.callback = func(i int, c []ChannelMessage) {
				Expect(len(c)).To(BeEquivalentTo(2))
				Expect(c[0].GetMessageAsString()).To(BeEquivalentTo("first"))
				Expect(c[1].GetMessageAsString()).To(BeEquivalentTo("second"))
				close(done)
			}

			go cluster.StartListening()

			logger.WriteString("first")
			loggerFactory.NewChanneledLogger("other-logger").WriteString("second")
			close(channelProvider.GetChannel()) //Ends all communication as we wanna unit test on main thread
		})

		It("Should cluster logs differently due to time", func(done Done) {
			cluster = NewLoggerCluster(pipeline, channelProvider, time.Nanosecond)

			pipeline.callback = func(i int, c []ChannelMessage) {
				if c[0].GetMessageAsString() == "first" {
					Expect(i).To(BeZero())
				}
				if c[0].GetMessageAsString() == "second" {
					Expect(i).To(BeEquivalentTo(1))
					close(done)
				}
			}

			go cluster.StartListening()
			go func() {
				logger.WriteString("first")
				time.Sleep(time.Second)
				loggerFactory.NewChanneledLogger("other-logger").WriteString("second")

				close(channelProvider.GetChannel()) //Ends all communication as we wanna unit test on main thread
			}()
		}, 5*1000) //We give the test a 5 second timeout, as we expect the loggers to take some time to call
	})
})
