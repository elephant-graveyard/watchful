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
	"net/http"
	"net/url"
	"time"

	"github.com/gonvenience/bunt"
	"github.com/homeport/watchful/pkg/logger"
	"github.com/homeport/watchful/pkg/merkhet"
)

// CurlMerkhet is an implementation of the Merkhet interface that curls against a domain
type CurlMerkhet struct {
	CurlDomain        *string
	BaseDomain        string
	BaseReference     merkhet.Base
	HTTPClient        *http.Client
	SingleAppProvider AppProvider
}

// NewDefaultCurlMerkhet creates a new curl merkhet instance
func NewDefaultCurlMerkhet(domain string, baseReference merkhet.Base, appProvider AppProvider) *CurlMerkhet {
	return NewCurlMerkhet(domain, baseReference, &http.Client{}, 30*time.Second, appProvider)
}

// NewCurlMerkhet creates a new curl merkhet instance
func NewCurlMerkhet(baseDomain string, baseReference merkhet.Base, httpClient *http.Client, timeout time.Duration, appProvider AppProvider) *CurlMerkhet {
	httpClient.Timeout = timeout
	return &CurlMerkhet{
		BaseDomain:        baseDomain,
		BaseReference:     baseReference,
		HTTPClient:        httpClient,
		SingleAppProvider: appProvider,
	}
}

// Install installs the merkhet. This does virtually nothing as curl doesn't need setup
func (m *CurlMerkhet) Install() error {
	parsedURL, e := url.Parse(m.BaseDomain)
	if e != nil {
		m.Base().Logger().WriteString(logger.Error, fmt.Sprintf("Could not parse url from domain %s", m.BaseDomain))
		return e
	}

	domain := parsedURL.Scheme + "://" + m.SingleAppProvider.AppName() + "." + parsedURL.Host
	m.CurlDomain = &domain
	return nil
}

// PostConnect is called after watchful was authorized against the cloud foundry cluster
func (m *CurlMerkhet) PostConnect() error {
	infoLog, errorLog, err := m.SingleAppProvider.Push(m.Base().Logger())
	if err != nil {
		m.Base().Logger().WriteString(logger.Error, "Could not post-connect upstream sample-app, printing logs")
		infoLog.Flush()
		errorLog.Flush()
		return err
	}

	m.Base().Logger().WriteString(logger.Info, "Post-Connected curl-merkhet")
	return nil
}

// Execute executes one single test. This will curl against the domain
func (m *CurlMerkhet) Execute() error {
	var curlDomain string
	if m.CurlDomain != nil {
		curlDomain = *m.CurlDomain
	} else {
		curlDomain = m.BaseDomain
	}

	response, err := m.HTTPClient.Get(curlDomain)
	if err != nil {
		m.BaseReference.Logger().WriteString(logger.Error, bunt.Sprintf("Red{Failed to curl: } %s", err.Error()))
		return err
	}

	if response.StatusCode != http.StatusOK {
		m.BaseReference.Logger().WriteString(logger.Error, bunt.Sprintf("Red{Failed to curl: } Response Code: %d", response.StatusCode))
		return fmt.Errorf("the domain %s returned status code %d", m.BaseDomain, response.StatusCode)
	}

	m.BaseReference.Logger().WriteString(logger.Debug, bunt.Sprintf("SpringGreen{Curled successfully}"))
	return nil
}

// Base returns the merkhet base instance of the curl merkhet
func (m *CurlMerkhet) Base() merkhet.Base {
	return m.BaseReference
}
