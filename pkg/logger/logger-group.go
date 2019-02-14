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

// GroupContainer contains all the logger groups
//
// NewGroup creates a new Group but returns the GroupContainer for future interaction
//
// GroupByLogger returns the group the passed logger instance is a part of or nil of not found
//
// Groups returns the list of all groups
//
// GroupCount returns the amount of groups registered
type GroupContainer interface {
	NewGroup(logger ...Logger) GroupContainer
	GroupByLogger(logger Logger) Group
	Groups() []Group
	GroupCount() int
}

// HotAccessGroupContainer is an implementation of the GroupContainer interface that also stores a direct map
// of logger id to group for quicker access
type HotAccessGroupContainer struct {
	IDToGroupMap    []Group
	GroupSlice      []Group
	GroupCountCache int
}

// NewGroup creates a new Group but returns the GroupContainer for future interaction
func (g *HotAccessGroupContainer) NewGroup(logger ...Logger) GroupContainer {
	loggerGroup := NewSlicedGroup(g.GroupCount(), logger)
	g.GroupSlice = append(g.GroupSlice, loggerGroup)
	g.GroupCountCache = len(g.GroupSlice)

	for _, l := range logger {
		currentIDMapLength := len(g.IDToGroupMap)
		if currentIDMapLength <= l.ID() {
			g.IDToGroupMap = append(g.IDToGroupMap, make([]Group, l.ID()-currentIDMapLength+1)...)
		}

		g.IDToGroupMap[l.ID()] = loggerGroup
	}
	return g
}

// GroupByLogger returns the group the passed logger instance is a part of or nil of not found
func (g *HotAccessGroupContainer) GroupByLogger(logger Logger) Group {
	if len(g.IDToGroupMap) > logger.ID() {
		return g.IDToGroupMap[logger.ID()]
	}
	return nil
}

// Groups returns the list of all groups
func (g *HotAccessGroupContainer) Groups() []Group {
	return g.GroupSlice
}

// GroupCount returns the amount of groups
func (g *HotAccessGroupContainer) GroupCount() int {
	return g.GroupCountCache
}

// NewGroupContainer creates a new hot access container instance
func NewGroupContainer() *HotAccessGroupContainer {
	return &HotAccessGroupContainer{}
}

// Group represents a group of loggers all displayed in the same pipe
//
// Loggers will return the loggers in the group
//
// MaxPrefixSize returns the maximum length of the prefix
//
// ID returns the group id
type Group interface {
	Loggers() []Logger
	MaxPrefixSize() int
	ID() int
}

// SlicedGroup is a group implementation based on a slice
type SlicedGroup struct {
	LoggerSlice     []Logger
	MaxPrefixLength int
	IDValue         int
}

// Loggers will return the loggers in the group
func (g *SlicedGroup) Loggers() []Logger {
	return g.LoggerSlice
}

// MaxPrefixSize returns the maximum length of the prefix
func (g *SlicedGroup) MaxPrefixSize() int {
	return g.MaxPrefixLength
}

// ID returns the group id
func (g *SlicedGroup) ID() int {
	return g.IDValue
}

// NewSlicedGroup creates a new group instance containing the loggers
func NewSlicedGroup(id int, loggers []Logger) *SlicedGroup {
	var maxPrefix int
	for _, logger := range loggers {
		prefixLength := len(logger.AsPrefix())
		if prefixLength > maxPrefix {
			maxPrefix = prefixLength
		}
	}

	return &SlicedGroup{
		LoggerSlice:     loggers,
		MaxPrefixLength: maxPrefix,
		IDValue:         id,
	}
}
