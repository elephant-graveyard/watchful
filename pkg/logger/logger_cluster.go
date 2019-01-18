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

import (
	"sync"
	"time"
)

//--

//Cluster defines a service that is capable of clustering loggers into on syncronized output
type Cluster interface {
	//GetLoggerPipeline returns the cluster instance the Cluster send the coupled ChannelMessages to
	GetLoggerPipeline() Pipeline

	//GetChannelProvider returns the channel provider this Cluster uses
	GetChannelProvider() ChannelProvider

	//StartListening enables the Cluster instance to listen to the channel found in Cluster#getChannelProvider
	StartListening()

	//Flushes the cluster content to the pipeline
	Flush()

	//GetWaitGroup returns the waitgroup instance the cluster uses to identify if it is done running
	GetWaitGroup() *sync.WaitGroup
}

//simpleCluster is a default implementation of the Cluster interface
type simpleCluster struct {
	pipe            Pipeline
	channel         ChannelProvider
	clusterDuration time.Duration
	waitGroup       *sync.WaitGroup
	messageCache    []ChannelMessage
	nextCacheFlush  time.Time
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
	s.GetWaitGroup().Add(1)

	for incomingMessage := range s.GetChannelProvider().GetChannel() {
		if s.timedFlushNeeded() {
			if s.isLoggerMessageCached(incomingMessage.GetLogger().GetID()) { //The message is not flushed this release
				s.Flush()
				s.resetFlushTimer()

				s.appendMessage(incomingMessage)
			} else { //The message is flushed in this release
				s.appendMessage(incomingMessage)

				s.Flush()
				s.resetFlushTimer()
			}
		} else { //Time wise we don't need to flush
			if s.isLoggerMessageCached(incomingMessage.GetLogger().GetID()) { //Overwriting logger message, gotta flush
				s.Flush()
				s.resetFlushTimer()
			}

			s.appendMessage(incomingMessage) //Add message to either this or new release, depending on prev if
		}
	}

	s.Flush()
	s.GetWaitGroup().Done()
}

//Flush flushes all cached messages to the pipeline
func (s *simpleCluster) Flush() {
	if len(s.messageCache) < 1 {
		return
	}

	s.GetLoggerPipeline().Write(s.messageCache)
	s.resetMessageCache()
}

//GetWaitGroup returns the waitgroup used to signal if the listening is running at the momement
func (s *simpleCluster) GetWaitGroup() *sync.WaitGroup {
	return s.waitGroup
}

//resetMessageCache simply clears the message cache
func (s *simpleCluster) resetMessageCache() {
	s.messageCache = nil
	//signal the garbage collector to remove eventually cached array structure under it: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
	//nil slices will still appedabe so we don't need to reallocate a slice instance using make
}

//isLoggerMessageCached returns if the logger id already has a cached ChannelMessage instance assigned
func (s *simpleCluster) isLoggerMessageCached(loggerID int) bool {
	for _, i := range s.messageCache {
		if i.GetLogger().GetID() == loggerID {
			return true
		}
	}
	return false
}

//timedFlushNeeded returns if the cluster needs to flush its values due to the time limit being reached
func (s *simpleCluster) timedFlushNeeded() bool {
	return s.nextCacheFlush.Before(time.Now())
}

//resetFlushTimer resets the flush timer to now
func (s *simpleCluster) resetFlushTimer() {
	s.nextCacheFlush = time.Now().Add(s.clusterDuration)
}

//appendMessage adds a message to the message cache of the cluster
func (s *simpleCluster) appendMessage(c ChannelMessage) {
	s.messageCache = append(s.messageCache, c)
}

//NewLoggerCluster creates a new instance of the LoggerCluster inferface
func NewLoggerCluster(pipeline Pipeline, channelProvider ChannelProvider, clusterDuration time.Duration) Cluster {
	return &simpleCluster{
		channel:         channelProvider,
		pipe:            pipeline,
		clusterDuration: clusterDuration,
		nextCacheFlush:  time.Now().Add(clusterDuration),
		waitGroup:       &sync.WaitGroup{},
	}
}
