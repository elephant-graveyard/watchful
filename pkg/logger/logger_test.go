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
	"os"
	"strings"
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

		It("should have the correct name", func() {
			Expect(logger.Name()).To(BeEquivalentTo("test-logger"))
		})

		It("should forward messages correctly", func() {
			logger.WriteString("test", Info)
			logger.WriteString("more-tests", Info)
			Expect(channelProvider.Read().MessageAsString()).To(BeEquivalentTo("test"))
			Expect(channelProvider.Read().MessageAsString()).To(BeEquivalentTo("more-tests"))
		})

		It("should cluster same logs differently", func(done Done) {
			pipeline.callback = func(i int, c []ChannelMessage) {
				if c[0].MessageAsString() == "first" {
					Expect(i).To(BeZero())
				}
				if c[0].MessageAsString() == "second" {
					Expect(i).To(BeEquivalentTo(1))
					close(done)
				}
			}

			go cluster.StartListening()

			logger.WriteString("first", Info)
			logger.WriteString("second", Info)
			close(channelProvider.Channel()) // Ends all communication as we wanna unit test on main thread
		})

		It("should cluster different logs together", func(done Done) {
			pipeline.callback = func(i int, c []ChannelMessage) {
				Expect(len(c)).To(BeEquivalentTo(2))
				Expect(c[0].MessageAsString()).To(BeEquivalentTo("first"))
				Expect(c[1].MessageAsString()).To(BeEquivalentTo("second"))
				close(done)
			}

			go cluster.StartListening()

			logger.WriteString("first", Info)
			loggerFactory.NewChanneledLogger("other-logger").WriteString("second", Info)
			close(channelProvider.Channel()) // Ends all communication as we wanna unit test on main thread
		})

		It("should cluster logs differently due to time", func(done Done) {
			cluster = NewLoggerCluster(pipeline, channelProvider, time.Nanosecond)

			pipeline.callback = func(i int, c []ChannelMessage) {
				if c[0].MessageAsString() == "first" {
					Expect(i).To(BeZero())
				}
				if c[0].MessageAsString() == "second" {
					Expect(i).To(BeEquivalentTo(1))
					close(done)
				}
			}

			go cluster.StartListening()
			go func() {
				logger.WriteString("first", Info)
				time.Sleep(time.Second)
				loggerFactory.NewChanneledLogger("other-logger").WriteString("second", Info)

				close(channelProvider.Channel()) // Ends all communication as we wanna unit test on main thread
			}()
		}, 5*1000) // We give the test a 5 second timeout, as we expect the loggers to take some time to call

		It("should chunk slices correctly", func() {
			first, second, third, fourth := "11111", "2222\n2", "\033[31m"+"33333", "44"+"\033[31m"+"444"

			result := ChunkSlice(first+second+third+fourth, 4)

			Expect(len(result)).To(BeEquivalentTo(5))
		})

		It("should log different levels", func() {
			logger.ReportingOn(Info).Write([]byte("info"))
			logger.ReportingOn(Error).Write([]byte("error"))

			infoMessage := channelProvider.Read()
			Expect(infoMessage.Level).To(BeEquivalentTo(Info))
			Expect(infoMessage.Message).To(BeEquivalentTo("info"))

			errorMessage := channelProvider.Read()
			Expect(errorMessage.Level).To(BeEquivalentTo(Error))
			Expect(errorMessage.Message).To(BeEquivalentTo("error"))
		})

		It("should notify the reviewers", func() {
			logger.ReportingOn(Info).
				ReviewWith(func(p []byte) error {
					Expect(string(p)).To(BeEquivalentTo("info"))
					return nil
				}).
				Write([]byte("info"))
		})

		It("Should send messages to the terminal correctly", func(done Done) {
			loggerGroupConfig := []int{0, 1}
			p := NewSplitPipeline(NewSplitPipelineConfig(true, *time.FixedZone("UTC", 0), 80, loggerGroupConfig), os.Stdout) // 200 is the fixed size of a terminal this thing will use, as ginkgo overwrites it
			c := NewLoggerCluster(p, channelProvider, time.Second)
			other := loggerFactory.NewChanneledLogger("other")
			p.Observer(func(s string) {
				if strings.Contains(s, "done") {
					Succeed()
					close(done)
				} else {
					Expect(s).To(ContainSubstring("[%s]", logger.Name()))
				}
			})

			go c.StartListening()

			go func() {
				logger.WriteString("\033[31m1", Info)
				logger.WriteString("\033[31m2", Info)
				other.WriteString("\033[32mdone", Info)
				close(channelProvider.Channel())
			}()
		}, 5*1000)
	})
})
