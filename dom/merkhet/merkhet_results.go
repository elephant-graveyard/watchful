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

package merkhet

//Result defines a statistic object that contains the data of a merkhet run
type Result interface {
	GetTotalRuns() uint
	GetFailedRuns() uint
	IsValid() bool
}

//simpleResult is a small implementation of the Result interface
type simpleResult struct {
	TotalRuns  uint
	FailedRuns uint
	Valid      bool
}

//GetTotalRuns returns the total amount of runs the merkhet instance ran
//at the time this result instance was created
func (s *simpleResult) GetTotalRuns() uint {
	return s.TotalRuns
}

//GetFailedRuns returns the total amount of faild runs the merkhet instance that build
//this result produced
func (s *simpleResult) GetFailedRuns() uint {
	return s.FailedRuns
}

//IsValid returns if the result was marked valid by the Merkhet instance that build it
func (s *simpleResult) IsValid() bool {
	return s.Valid
}

//NewMerkhetResult creates a new instance of the MerkhetResult interface
func NewMerkhetResult(totalRuns uint, failedRuns uint, valid bool) Result {
	return &simpleResult{
		TotalRuns:  totalRuns,
		FailedRuns: failedRuns,
		Valid:      valid,
	}
}
