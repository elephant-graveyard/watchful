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

package merkhet_test

import (
	"time"

	. "github.com/homeport/watchful/pkg/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Merkhet code test", func() {
	Context("Testing default behaviour of merkhet instances", func() {
		var (
			pool     Pool
			callback *MerkhetCallback
			merkhet  Merkhet
		)

		BeforeEach(func() {
			pool = NewPool()
			callback = &MerkhetCallback{}
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 0), 0, 0, true, callback)
		})

		It("should have 1 merkhet", func() {
			pool.StartWorker(merkhet, time.Second, nil)
			Expect(pool.Size()).To(BeNumerically("==", 1))
		})

		It("should have the correct name", func() {
			Expect(merkhet.Base().Configuration().Name()).To(BeEquivalentTo("test-config"))
		})

		It("should have installed", func() {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 2), 10, 2, true, &MerkhetCallback{
				onInstall: func() error {
					Succeed()
					return nil
				},
			})

			pool.StartWorker(merkhet, time.Second, nil)

			err := pool.ForEach(ConsumeSync(func(m Merkhet, future Future) {
				future.Complete(m.Install())
			})).Wait().FirstError()

			Expect(err).To(BeNil())
			pool.Shutdown()
		})

		It("should have executed correctly", func(done Done) {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 2), 0, 0, true, &MerkhetCallback{
				onExecute: func() error {
					time.Sleep(time.Second)
					return nil
				},
			})

			pool.StartWorker(merkhet, time.Second, nil)

			pool.ForEach(ConsumeAsync(func(m Merkhet, future Future) {
				future.Complete(merkhet.Execute())
				Succeed()
				close(done)
			}))
		}, 5*1000)

		It("should not produce a data race and record a successful run", func(done Done) {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 2), 2, 2, true, &MerkhetCallback{
				onExecute: func() error {
					time.Sleep(time.Second)
					return nil
				},
			})

			pool.StartWorker(merkhet, time.Second, nil)
			er := pool.ForEach(ConsumeAsync(func(m Merkhet, future Future) {
				future.Complete(m.Execute())
				m.Base().RecordSuccessfulRun()
			})).Wait().FirstError()

			Expect(er).To(BeNil())
			Expect(merkhet.Base().NewResultSet().SuccessfulRuns()).To(BeEquivalentTo(1))
			close(done)
		}, 5*1000)

		It("should beat correctly", func(done Done) {
			pool.StartWorker(merkhet, 10*time.Millisecond, ConsumeSync(func(m Merkhet, future Future) {
				m.Base().RecordSuccessfulRun()
				future.Complete(nil)
			}))
			pool.StartHeartbeats()
			time.Sleep(time.Second)
			Expect(len(pool.BeatingHearts())).To(BeEquivalentTo(1))
			pool.Shutdown()
			Expect(len(pool.BeatingHearts())).To(BeEquivalentTo(0))

			Expect(merkhet.Base().NewResultSet().SuccessfulRuns()).To(BeNumerically("~", 100, 2))
			close(done)
		}, 10*1000)

		It("should pass the merkhet test using a flat config", func() {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 2), 10, 2, true, callback)
			Expect(merkhet.Base().NewResultSet().Valid()).To(BeTrue())
		})

		It("should not pass the merkhet test using a flat config", func() {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 1), 10, 2, true, callback)
			Expect(merkhet.Base().NewResultSet().Valid()).To(BeFalse())
		})

		It("should pass the merkhet test using a percentage config", func() {
			merkhet = NewMerkhetMock(NewPercentageConfiguration("test-config", 0.2), 10, 2, true, callback)
			Expect(merkhet.Base().NewResultSet().Valid()).To(BeTrue())
		})

		It("should not pass the merkhet test using a percentage config", func() {
			merkhet = NewMerkhetMock(NewPercentageConfiguration("test-config", 0.1), 10, 2, true, callback)
			Expect(merkhet.Base().NewResultSet().Valid()).To(BeFalse())
		})
	})
})
