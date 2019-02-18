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
	"fmt"
	"github.com/homeport/watchful/pkg/merkhet"
	"net/http"
	"time"
)

// CurlMerkhet is an implementation of the Merkhet interface that curls against a domain
type CurlMerkhet struct {
	Domain        string
	BaseReference merkhet.Base
	HTTPClient    *http.Client
}

// NewDefaultCurlMerkhet creates a new curl merkhet instance
func NewDefaultCurlMerkhet(domain string, baseReference merkhet.Base) *CurlMerkhet {
	return NewCurlMerkhet(domain, baseReference, &http.Client{} , 115 * time.Second)
}

// NewCurlMerkhet creates a new curl merkhet instance
func NewCurlMerkhet(domain string, baseReference merkhet.Base, httpClient *http.Client, timeout time.Duration) *CurlMerkhet {
	httpClient.Timeout = timeout
	return &CurlMerkhet{Domain: domain, BaseReference: baseReference, HTTPClient: httpClient}
}

// Install installs the merkhet. This does virtually nothing as curl doesn't need setup
func (m *CurlMerkhet) Install() error {
	return nil
}

// PostConnect is called after watchful was authorized against the cloud foundry cluster
func (m *CurlMerkhet) PostConnect() error {
	return nil
}

// Execute executes one single test. This will curl against the domain
func (m *CurlMerkhet) Execute() error {
	response, err := m.HTTPClient.Get(m.Domain)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("the domain %s returned status code %d", m.Domain, response.StatusCode)
	}
	return nil
}

// Base returns the merkhet base instance of the curl merkhet
func (m *CurlMerkhet) Base() merkhet.Base {
	return m.BaseReference
}
