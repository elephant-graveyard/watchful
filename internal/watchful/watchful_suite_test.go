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

package watchful_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWatchful(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "watchful internal suite")
}

type CloudFoundry struct {
	Server *httptest.Server
}

func NewFakeCloudFoundry() CloudFoundry {
	r := mux.NewRouter()

	r.HandleFunc("/v2/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, asJson(map[string]interface{}{
			"name": "Fake Cloud Foundry",
		}))
	})

	r.HandleFunc("/v2/apps", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, asJson(map[string]interface{}{
			"total_results": 0,
			"total_pages":   0,
			"prev_url":      nil,
			"next_url":      nil,
			"resources":     []map[string]interface{}{},
		}))
	})

	// catch all unhandled paths
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Unexpected call: %v", spew.Sdump(r))

		w.WriteHeader(500)
		fmt.Fprintf(w, "path not implemented")
	})

	return CloudFoundry{
		Server: httptest.NewServer(r),
	}
}

func (cf *CloudFoundry) Teardown() {
	if cf.Server != nil {
		cf.Server.Close()
		cf.Server = nil
	}
}

func asJson(input map[string]interface{}) string {
	bytes, err := json.Marshal(input)
	Expect(err).To(BeNil())

	return string(bytes)
}
