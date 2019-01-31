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
	"encoding/base64"
	"fmt"
	"os"
	"path"
)

var (
	// SampleGoApp is the content of the sample app written in go. This string will be encoded in base64
	SampleGoApp string
)

// Extract64 extracts the string content passed as an asset into the provided file.
// This will create the file and all parent directories if needed
// The provided string has to be encoded in base64
func Extract64(asset string, file string) error {
	directoryName := path.Dir(file)
	fileName := path.Base(file)

	if e := os.MkdirAll(directoryName, 0700); e != nil { // 0700 -> leading 0 is stripped -> 700 means write, read and execute for the user creating this file.
		return e
	}

	fileStream, e := os.Create(path.Join(directoryName, fileName))
	if e != nil {
		return e
	}
	defer fileStream.Close()

	decodedBytes, e := base64.StdEncoding.DecodeString(asset)
	if e != nil {
		return e
	}

	fmt.Fprint(fileStream, string(decodedBytes))
	return nil
}
