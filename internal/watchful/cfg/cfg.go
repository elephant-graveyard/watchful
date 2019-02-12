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

package cfg

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ParseFromFile parses the contents of the file located under the given string
// into the config instance passed as the second parameter.
// The function returns an error if one occurred or nil if everything worked fine
func ParseFromFile(file string, config interface{}) error {
	fileContent, e := ioutil.ReadFile(file)
	if e != nil {
		return e
	}

	return yaml.UnmarshalStrict(fileContent, config)
}

// ParseFromString parses contents of the string provided as a parameter
// into the config instance passed as the second parameter.
// The function returns an error if one occurred or nil if everything worked fine
func ParseFromString(content string, config interface{}) error {
	return yaml.UnmarshalStrict([]byte(content), config)
}
