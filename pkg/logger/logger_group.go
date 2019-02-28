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
// UpdateGroup updates the group of the logger to fit the given group
//
// GroupByLogger returns the group the passed logger instance is a part of or nil of not found
//
// Groups returns the list of all groups
//
// GroupCount returns the amount of groups registered
type GroupContainer interface {
	NewGroup(logger ...Logger) GroupContainer
	UpdateGroup(logger Logger, group Group)
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
	loggerGroup := NewSlicedGroup(g, g.GroupCount(), logger)
	g.GroupSlice = append(g.GroupSlice, loggerGroup)
	g.GroupCountCache = len(g.GroupSlice)

	for _, l := range logger {
		g.UpdateGroup(l, loggerGroup)
	}
	return g
}

// UpdateGroup updates the group of the logger to fit the given group
func (g *HotAccessGroupContainer) UpdateGroup(logger Logger, group Group) {
	currentIDMapLength := len(g.IDToGroupMap)
	if currentIDMapLength <= logger.ID() {
		g.IDToGroupMap = append(g.IDToGroupMap, make([]Group, logger.ID()-currentIDMapLength+1)...)
	} else {
		priorGroup := g.IDToGroupMap[logger.ID()]
		if priorGroup != nil {
			priorGroup.Remove(logger)
		}
	}

	g.IDToGroupMap[logger.ID()] = group
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
//
// Add adds a logger to the group
//
// Remove removes a logger from a group
type Group interface {
	Loggers() []Logger
	MaxPrefixSize() int
	ID() int
	Add(logger Logger)
	Remove(logger Logger)
}

// LinkedSlicedGroup is a group implementation based on a slice
type LinkedSlicedGroup struct {
	LoggerSlice     []Logger
	MaxPrefixLength int
	IDValue         int
	Parent          GroupContainer
}

// Add adds a logger to the group
func (g *LinkedSlicedGroup) Add(logger Logger) {
	g.LoggerSlice = append(g.LoggerSlice, logger)
	g.updateMaxPrefix()
	g.Parent.UpdateGroup(logger, g)
}

// Remove removes a logger from a group
func (g *LinkedSlicedGroup) Remove(logger Logger) {
	for i, log := range g.LoggerSlice {
		if log.ID() == logger.ID() {
			g.LoggerSlice[i] = g.LoggerSlice[0]
			g.LoggerSlice = g.LoggerSlice[1:]
			return
		}
	}
}

// Loggers will return the loggers in the group
func (g *LinkedSlicedGroup) Loggers() []Logger {
	return g.LoggerSlice
}

// MaxPrefixSize returns the maximum length of the prefix
func (g *LinkedSlicedGroup) MaxPrefixSize() int {
	return g.MaxPrefixLength
}

// ID returns the group id
func (g *LinkedSlicedGroup) ID() int {
	return g.IDValue
}

// updateMaxPrefix updates the maximum prefix
func (g *LinkedSlicedGroup) updateMaxPrefix() {
	var maxPrefix int
	for _, logger := range g.LoggerSlice {
		prefixLength := len(logger.AsPrefix())
		if prefixLength > maxPrefix {
			maxPrefix = prefixLength
		}
	}
	g.MaxPrefixLength = maxPrefix
}

// NewSlicedGroup creates a new group instance containing the loggers
func NewSlicedGroup(parent GroupContainer, id int, loggers []Logger) *LinkedSlicedGroup {
	group := &LinkedSlicedGroup{
		Parent:      parent,
		LoggerSlice: loggers,
		IDValue:     id,
	}
	group.updateMaxPrefix()
	return group
}
