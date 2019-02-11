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

package utf8chunk

import "regexp"

// ColorEscapeSeq is the escape sequence for a color
const ColorEscapeSeq = "\x1b"

var (
	// BeginsWithColorPattern is the regex matching.
	// We are compiling it once, as regex is not the fastest
	BeginsWithColorPattern = regexp.MustCompile(ColorEscapeSeq + `\[\d+(;\d+)*m`)
)

// RemoveStartingColour removes the first colour code at index 0 of the string
func RemoveStartingColour(s string) (tail string, colour string) {
	if l := BeginsWithColorPattern.FindStringIndex(s); l != nil && l[0] == 0 {
		return s[l[1]:], s[l[0]:l[1]]
	}
	return s, ""
}
