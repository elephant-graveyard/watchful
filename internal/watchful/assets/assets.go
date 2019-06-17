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

package assets

import (
	"fmt"
	"strings"

	"github.com/homeport/pina-golada/pkg/files"
)

var (
	// Provider is the provider instance to access pina-goladas asset framework
	Provider ProviderInterface

	// LanguageMap represents a map of language to asset provider
	LanguageMap map[string]func() (directory files.Directory, e error)
)

// ProviderInterface is the interface that provides the asset file
// @pgl(package=assets&injector=Provider)
type ProviderInterface interface {
	// GetGoSampleApp returns the directory containing the go sample app
	// @pgl(asset=/assets/go-cf-sample/&compressor=tar)
	GetGoSampleApp() (directory files.Directory, e error)
}

// SampleApp returns the sample app stored for the given language
func SampleApp(language string) (directory files.Directory, e error) {
	provider := LanguageMap[strings.ToLower(language)]
	if provider != nil {
		return provider()
	}
	return nil, fmt.Errorf("could not find provider method for language %s", language)
}

// init initializes the languageMap
func init() {
	LanguageMap = make(map[string]func() (directory files.Directory, e error))

	LanguageMap["go"] = func() (directory files.Directory, e error) {
		return Provider.GetGoSampleApp()
	}
}
