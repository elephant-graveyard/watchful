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

// Cluster defines a service that is capable of clustering loggers into on synchronized output
//
// LoggerPipeline() returns the cluster instance the Cluster send the coupled ChannelMessages to
//
// ChannelProvider() returns the channel provider this Cluster uses
//
// StartListening enables the Cluster instance to listen to the channel provided to the cluster
// This will block all code executed, so use it in a go routine
//
// Flush() Flushes the cluster content to the pipeline
//
// WaitGroup() returns the wait-group instance the cluster uses to identify if it is done running
type Cluster interface {
	LoggerPipeline() Pipeline
	ChannelProvider() ChannelProvider
	StartListening()
	Flush()
	WaitGroup() *sync.WaitGroup
}

// SimpleCluster is a default implementation of the Cluster interface
type SimpleCluster struct {
	pipe            Pipeline
	channel         ChannelProvider
	clusterDuration time.Duration
	waitGroup       *sync.WaitGroup
	messageCache    []ChannelMessage
	nextCacheFlush  *time.Time
}

// LoggerPipeline returns the upstream pipeline the Cluster sends all its data to
func (s *SimpleCluster) LoggerPipeline() Pipeline {
	return s.pipe
}

// ChannelProvider returns the channel provider the Cluster uses
func (s *SimpleCluster) ChannelProvider() ChannelProvider {
	return s.channel
}

// StartListening enables the Cluster instance to listen to the channel provided to the cluster
// This will block all code executed, so use it in a go routine
func (s *SimpleCluster) StartListening() {
	s.WaitGroup().Add(1)

	for incomingMessage := range s.ChannelProvider().Channel() {
		if s.timedFlushNeeded() {
			s.Flush()
			s.resetFlushTimer()

			s.appendMessage(incomingMessage)
		} else { // Time wise we don't need to flush
			if s.nextCacheFlush == nil {
				s.resetFlushTimer() // When starting the listener and the first message is received, start the timer as well
			}

			if s.isLoggerMessageCached(incomingMessage.Logger.ID()) { // Overwriting logger message, gotta flush
				s.Flush()
				s.resetFlushTimer()
			}

			s.appendMessage(incomingMessage) // Add message to either this or new release, depending on prev if
		}
	}

	s.Flush()
	s.WaitGroup().Done()
}

// Flush flushes all cached messages to the pipeline
func (s *SimpleCluster) Flush() {
	if len(s.messageCache) < 1 {
		return
	}

	s.LoggerPipeline().Write(s.messageCache)
	s.resetMessageCache()
}

// WaitGroup returns the wait-group used to signal if the listening is running at the moment
func (s *SimpleCluster) WaitGroup() *sync.WaitGroup {
	return s.waitGroup
}

// resetMessageCache simply clears the message cache
func (s *SimpleCluster) resetMessageCache() {
	// signal the garbage collector to remove eventually cached array structure under it: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
	// nil slices will still appendable so we don't need to reallocate a slice instance using make
	s.messageCache = nil
}

// isLoggerMessageCached returns if the logger id already has a cached ChannelMessage instance assigned
func (s *SimpleCluster) isLoggerMessageCached(loggerID int) bool {
	for _, i := range s.messageCache {
		if i.Logger.ID() == loggerID {
			return true
		}
	}
	return false
}

// timedFlushNeeded returns if the cluster needs to flush its values due to the time limit being reached
func (s *SimpleCluster) timedFlushNeeded() bool {
	return s.nextCacheFlush != nil && s.nextCacheFlush.Before(time.Now())
}

// resetFlushTimer resets the flush timer to now
func (s *SimpleCluster) resetFlushTimer() {
	newDuration := time.Now().Add(s.clusterDuration)
	s.nextCacheFlush = &newDuration
}

// appendMessage adds a message to the message cache of the cluster
func (s *SimpleCluster) appendMessage(c ChannelMessage) {
	s.messageCache = append(s.messageCache, c)
}

// NewLoggerCluster creates a new instance of the LoggerCluster interface
func NewLoggerCluster(pipeline Pipeline, channelProvider ChannelProvider, clusterDuration time.Duration) *SimpleCluster {
	return &SimpleCluster{
		channel:         channelProvider,
		pipe:            pipeline,
		clusterDuration: clusterDuration,
		nextCacheFlush:  nil,
		waitGroup:       &sync.WaitGroup{},
	}
}
