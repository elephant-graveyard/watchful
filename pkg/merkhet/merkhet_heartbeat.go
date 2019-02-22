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

import (
	"os"
	"time"
)

// Heartbeat defines the heartbeat of a merkhet worker that beats every given duration
//
// StartBeating starts the heartbeat
//
// IsBeating returns if the heartbeat is currently beating
//
// StopBeating stops the heartbeat
//
// Worker returns the worker instance this heartbeat is working on
type Heartbeat interface {
	StartBeating()
	StopBeating()
	IsBeating() bool
	Worker() Worker
}

// TickedHeartbeat is a Heartbeat implementation based on a time.Ticker instance
type TickedHeartbeat struct {
	Ticker        *time.Ticker
	WrappedWorker Worker
	Consumer      Consumer
	Interval      time.Duration
	closure       chan os.Signal
}

// StartBeating starts the heartbeats go routine
func (t *TickedHeartbeat) StartBeating() {
	t.closure = make(chan os.Signal)

	if t.Consumer == nil {
		return
	}

	t.Ticker = time.NewTicker(t.Interval)
	go func() {
		t.Consumer.Consume(t.Worker().Merkhet(), NewFuture())

		for {
			select {
			case <-t.Ticker.C:
				t.Consumer.Consume(t.Worker().Merkhet(), NewFuture())
			case <-t.closure:
				return
			}
		}
	}()
}

// IsBeating returns if the heartbeat is currently beating
func (t *TickedHeartbeat) IsBeating() bool {
	return t.closure != nil
}

// StopBeating stops the heartbeats go routine
func (t *TickedHeartbeat) StopBeating() {
	if t.Ticker == nil {
		return
	}

	t.Ticker.Stop()

	t.closure <- os.Kill
	close(t.closure)
	t.closure = nil
}

// Worker returns the wrapped worker
func (t *TickedHeartbeat) Worker() Worker {
	return t.WrappedWorker
}

// NewTickedHeartbeat creates a new instance of the TickedHeartbeat struct
func NewTickedHeartbeat(worker Worker, duration time.Duration, consumer Consumer) *TickedHeartbeat {
	return &TickedHeartbeat{
		WrappedWorker: worker,
		Consumer:      consumer,
		Interval:      duration,
	}
}
