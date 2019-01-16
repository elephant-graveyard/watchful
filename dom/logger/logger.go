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

package logger

import "time"

//Logger defines an observable logger
type Logger interface {
	//PushObserver adds an observer function to the logger.
	//This function is executed whenever a write call is issued against the Logger intance
	PushObserver(func(bytes []byte))

	//GetName returns the name of the logger
	GetName() string

	//Write simply writes a byte array to the logger. This will trigger all observers
	Write(p []byte) (n int, err error)

	//WriteString writes all bytes of the string to the logger.
	//This will also trigger observers
	WriteString(s string) error

	//Clear removes all logged bytes from the Logger and returns them for future use
	Clear() []byte
}

//--

//Coupler defines a service that is capable of coupling loggers into on syncronized output
type Coupler interface {
	//GetFormatter returns the formatter instance the Coupler uses to combine the loggers handled
	GetFormatter() CoupledFormatter

	//RegisterLogger registers a logger into the coupler, effectivly making the coupler handle the bytes written to the logger
	RegisterLogger(logger Logger)
}

//--

//CoupledFormatter deinfes an interface that is capable of formatting a LoggerCoupler output
type CoupledFormatter interface {
	//FormatLoggerOutput formats all passed byte arrays into one final string
	FormatLoggerOutput(outputs [][]byte) string

	//GetLocation returns the location used to determin the date that is passed into the logs
	GetLocation() time.Location
}
