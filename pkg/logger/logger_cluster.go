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

//--

//Cluster defines a service that is capable of clustering loggers into on syncronized output
type Cluster interface {
	//GetLoggerPipeline returns the cluster instance the Cluster send the coupled ChannelMessages to
	GetLoggerPipeline() Pipeline

	//GetChannelProvider returns the channel provider this Cluster uses
	GetChannelProvider() ChannelProvider

	//StartListening enables the Cluster instance to listen to the channel found in Cluster#getChannelProvider
	StartListening()
}

//simpleCluster is a default implementation of the Cluster interface
type simpleCluster struct {
	pipe           Pipeline
	channel        ChannelProvider
	messageCache   []ChannelMessage
	nextCacheFlush time.Time
}

//GetLoggerPipeline returns the upstream pipeline the Cluster sends all its data to
func (s *simpleCluster) GetLoggerPipeline() Pipeline {
	return s.pipe
}

//GetChannelProvider returns the channel provider the Cluster wi
func (s *simpleCluster) GetChannelProvider() ChannelProvider {
	return s.channel
}

//StartListening enables the Cluster instance to listen to the channel found in Cluster#getChannelProvider
//This will block all code executed, so use it in a go routine
func (s *simpleCluster) StartListening() {
	for incomingMessage := range s.GetChannelProvider().GetChannel() {

		currentTimestamp := time.Now()
		if s.nextCacheFlush.After(currentTimestamp) {
			s.GetLoggerPipeline().Write(s.messageCache)
			s.resetMessageCache()
			s.nextCacheFlush = currentTimestamp.Add(time.Second)
		}

		s.messageCache = append(s.messageCache, incomingMessage)
	}
}

//resetMessageCache simply clears the message cache
func (s *simpleCluster) resetMessageCache() {
	s.messageCache = nil //signal the garbage collector to remove eventually cached array structure under it: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
	//nil slices will still appedabe so we don't need to reallocate a slice instance using make
}
