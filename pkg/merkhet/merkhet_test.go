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
	. "github.com/homeport/disrupt-o-meter/pkg/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Merkhet code test", func() {
	Context("Testing default behaviour of merkhet instances", func() {
		var (
			container Container
			merkhet   Merkhet
		)

		BeforeEach(func() {
			container = NewContainer()
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 0), 10, 0, true)
		})

		It("Should have 1 merkhats", func() {
			container.Push(merkhet)
			Expect(container.Size()).To(BeNumerically("==", 1))
		})

		It("Should have the correct name", func() {
			Expect(merkhet.GetConfiguration().GetName()).To(BeEquivalentTo("test-config"))
		})

		It("Should have installed", func() {
			container.Push(merkhet)
			container.ForEach(func(m Merkhet) {
				m.Install()
			})

			//TODO check for result
		})

		It("Should pass the merkhet test using a flat config", func() {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 2), 10, 2, true)
			Expect(merkhet.BuildResult().IsValid()).To(BeTrue())
		})

		It("Should not pass the merkhet test using a flat config", func() {
			merkhet = NewMerkhetMock(NewFlatConfiguration("test-config", 1), 10, 2, true)
			Expect(merkhet.BuildResult().IsValid()).To(BeFalse())
		})

		It("Should pass the merkhet test using a percentage config", func() {
			merkhet = NewMerkhetMock(NewPercentageConfiguration("test-config", 0.2), 10, 2, true)
			Expect(merkhet.BuildResult().IsValid()).To(BeTrue())
		})

		It("Should not pass the merkhet test using a percentage config", func() {
			merkhet = NewMerkhetMock(NewPercentageConfiguration("test-config", 0.1), 10, 2, true)
			Expect(merkhet.BuildResult().IsValid()).To(BeFalse())
		})
	})
})
