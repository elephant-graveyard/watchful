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

//Container contains a map of registered merkhets as well as manages them
type Container interface {
	Push(merkhet Merkhet)
	Size() uint
	InstallAll()
	PostConnectAll()
	ExecuteAll()
}

//simpleContainer is a basic implementation of the MerkhetContainer interface
type simpleContainer struct {
	RegisteredMerkhets []Merkhet
}

//Push pushes a new merkhet instance into the container
func (s *simpleContainer) Push(m Merkhet) {
	s.RegisteredMerkhets = append(s.RegisteredMerkhets, m)
}

func (s *simpleContainer) Size() uint {
	return uint(len(s.RegisteredMerkhets))
}

//InstallAll installs all Merkhets in the container
//This should be called before the main workflow is executed, as it may take some time
func (s *simpleContainer) InstallAll() {
	for _, merkhet := range s.RegisteredMerkhets {
		merkhet.Install()
	}
}

//PostConnectAll is calling the PostConnect method on every merkhet in the container
func (s *simpleContainer) PostConnectAll() {
	for _, merkhet := range s.RegisteredMerkhets {
		merkhet.PostConnect()
	}
}

//PostConnectAll is calling the Execute method on every merkhet in the container
//This methode can be used during testing phase
func (s *simpleContainer) ExecuteAll() {
	for _, merkhet := range s.RegisteredMerkhets {
		merkhet.Execute()
	}
}

//NewContainer returns a fresh empty instance of the go container
func NewContainer() Container {
	return &simpleContainer{
		RegisteredMerkhets: make([]Merkhet, 0),
	}
}
