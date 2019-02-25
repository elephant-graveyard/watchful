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

package merkhets

import (
	"github.com/homeport/watchful/pkg/merkhet"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

var _ = Describe("the implemented merkhets should function correctly" , func() {

	var (
		Server *MockedServer
		MerkhetBase merkhet.Base
	)

	_ = BeforeEach(func() {
		Server = NewMockedServer()
		Server.Start()

		MerkhetBase = merkhet.NewSimpleBase(&ConsoleLogger{} , merkhet.NewFlatConfiguration("config" , 0))
	})

	_ = AfterEach(func() {
		Server.Server.Close()
		Server.Server = nil
	})

	_ = It("should curl correctly" , func() {
		Server.Route("/info" , func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World"))
			w.WriteHeader(http.StatusOK)
		})

		curlMerkhet := NewDefaultCurlMerkhet(Server.Server.URL + "/info", MerkhetBase , "")
		Expect(curlMerkhet.Execute()).To(BeNil())
	})

	_ = It("should fail curl due to timeout" , func(done Done) {
		Server.Route("/info" , func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
		})

		curlMerkhet := NewCurlMerkhet(Server.Server.URL + "/info", MerkhetBase , &http.Client{} , time.Millisecond , "")
		Expect(curlMerkhet.Execute()).To(Not(BeNil()))
		close(done)
	} , 5 * 1000)

	_ = It("should fail curl due to wrong status code" , func(done Done) {
		Server.Route("/info" , func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		})

		curlMerkhet := NewDefaultCurlMerkhet(Server.Server.URL + "/info", MerkhetBase , "")
		Expect(curlMerkhet.Execute()).To(Not(BeNil()))
		close(done)
	})
})
