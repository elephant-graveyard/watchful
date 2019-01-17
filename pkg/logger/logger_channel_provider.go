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

//MessageChannel defines the type of a channel that can be used to send ChannelMessgae instances
type MessageChannel chan ChannelMessage

//ChannelProvider is a small provider instance that provides a reference to the channel loggers use to communicate with the logger coupler instance
type ChannelProvider interface {

	//Get returns the message channel instance this provider is providing
	Get() MessageChannel
}

//simpleChannelProvider is a basic struct base implementation of the of ChannelProvider interface
type simpleChannelProvider struct {
	Channel MessageChannel
}

//Get returns the instance this ChannelProvider is wrapping
func (s *simpleChannelProvider) Get() MessageChannel {
	return s.Channel
}

//NewChannelProvider creates a new provider providing a fresh channel instance
func NewChannelProvider() ChannelProvider {
	return NewChannelProviderOf(make(chan ChannelMessage))
}

//NewChannelProviderOf returns a channel provider that wraps the passed channel instance
func NewChannelProviderOf(channel MessageChannel) ChannelProvider {
	return &simpleChannelProvider{
		Channel: channel,
	}
}

//--

//ChannelMessage is the content of the Logger channel and contains the logger name as well as the message the logger produced
type ChannelMessage struct {
	Logger  Logger
	Message string
}

//NewChannelMessage is a simple constructor for the ChannelMessage struct
func NewChannelMessage(logger Logger, message string) ChannelMessage {
	return ChannelMessage{
		Logger:  logger,
		Message: message,
	}
}
