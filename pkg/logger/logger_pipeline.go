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
	"fmt"
	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"io"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/homeport/watchful/pkg/utf8chunk"

	"github.com/homeport/gonvenience/pkg/v1/term"
)

// PipelineObserver is an observer that can be registered on a pipeline
type PipelineObserver func(messages string)

var (
	// PipelineSeparator is the char that is used to separate the different pipelines
	PipelineSeparator = ` | `

	// TimeFormat is the format we will adapt the time to
	TimeFormat = "[2/1/2006 15:04:05]"
)

// Pipeline defines an interface that is capable of formatting a LoggerCoupler output
//
// Write formats all passed byte arrays into one final string
//
// Observer adds a new observer to the pipeline that will be triggered on write
type Pipeline interface {
	Write(messages []ChannelMessage)
	Observer(observer PipelineObserver)
}

// SplitPipeline is a logger pipeline that splits the logger output into different channels and produces a split screen view
type SplitPipeline struct {
	config    SplitPipelineConfig
	writer    io.Writer
	observers []PipelineObserver
}

// Write splits the messages
func (s *SplitPipeline) Write(messages []ChannelMessage) {
	if len(messages) < 1 {
		return
	}

	finalOutput := make([][]string, s.config.GroupContainer.GroupCount())

	for _, m := range messages {
		loggerGroup := s.config.GroupContainer.GroupByLogger(m.Logger)
		if loggerGroup == nil {
			fmt.Fprintln(s.writer, "Received a logger with the id ", m.Logger.ID(), "that could not be sorted")
		}

		message := m.MessageAsString()

		var slices []string
		if s.config.ShowLoggerName {
			// Slice it but with prefix
			slices = ChunkSlice(message, s.config.CharacterPerPipe-loggerGroup.MaxPrefixSize())

			loggerPrefix := m.Logger.AsPrefix()
			for index, slice := range slices {
				slices[index] = loggerPrefix + slice
			}
		} else {
			slices = ChunkSlice(message, s.config.CharacterPerPipe) // slice it normally
		}

		loggerGroupOutput := finalOutput[loggerGroup.ID()] // Take output till now

		loggerGroupOutput = append(loggerGroupOutput, slices...) // append new output
		finalOutput[loggerGroup.ID()] = loggerGroupOutput        // set output
	}

	for line := 0; true; line++ {

		output := ""
		empty := true

		for group := 0; group < s.config.GroupContainer.GroupCount(); group++ {
			groupArray := finalOutput[group]
			groupArraySize := len(groupArray)
			if group > 0 { // Print new line if we are not first or last message
				output += PipelineSeparator
			}

			if line < groupArraySize {
				groupOutput := groupArray[line]
				output += groupOutput

				if s.config.GroupContainer.GroupCount() > 1 {
					emptyChars := s.config.CharacterPerPipe - bunt.PlainTextLength(groupOutput) // The prefix is already in the group output
					output += strings.Repeat(" ", emptyChars)
				}
				output += "\u001b[0m"

				empty = false
			} else {
				output += strings.Repeat(" ", s.config.CharacterPerPipe)
			}
		}

		if !empty {
			for _, observer := range s.observers {
				observer(output)
			}

			fmt.Fprintf(s.writer, "%s\n", time.Now().Format(TimeFormat)+output) // Trim leftover empty spaces
			empty = true                                                        // Reset for next run
		} else {
			break
		}
	}
}

// Observer adds a new observer
func (s *SplitPipeline) Observer(o PipelineObserver) {
	s.observers = append(s.observers, o)
}

// ChunkSlice slice a message chunk into smaller chunks
func ChunkSlice(input string, chunkSize int) []string {
	chunkBounds := make([]int, 1)
	result := make([]string, 0)

	s := input
	var (
		printableCount, index int
	)

	for len(s) > 1 {
		m, c := utf8chunk.RemoveStartingColour(s)
		if len(c) != 0 {
			s = m // Simply ignore the color code and jump forward
			index += len(c)
		} else {
			r, _ := utf8.DecodeRuneInString(s)
			index++
			if unicode.IsPrint(r) {
				printableCount++
				if printableCount >= chunkSize {
					chunkBounds = append(chunkBounds, index)
					printableCount = 0
				}
			}

			s = s[1:]
		}
	}

	for index, lowerBound := range chunkBounds {
		var stringPiece string
		if index >= len(chunkBounds)-1 {
			stringPiece = input[lowerBound:]
		} else {
			upperBounds := chunkBounds[index+1]
			stringPiece = input[lowerBound:upperBounds]
		}
		result = append(result, stringPiece)
	}
	return result
}

// NewSplitPipeline returns a new instance of the SplitPipeline object
func NewSplitPipeline(config SplitPipelineConfig, writer io.Writer) *SplitPipeline {
	return &SplitPipeline{
		config: config,
		writer: writer,
	}
}

// SplitPipelineConfig is the config used in a split pipeline
type SplitPipelineConfig struct {
	ShowLoggerName   bool
	Location         *time.Location
	TerminalWidth    int
	GroupContainer   GroupContainer
	CharacterPerPipe int
}

// NewBasicSplitPipelineConfig creates a new split pipeline config using the default terminal length
func NewBasicSplitPipelineConfig(showLoggerName bool, location *time.Location, groups GroupContainer) SplitPipelineConfig {
	w := term.GetTerminalWidth()
	return NewSplitPipelineConfig(showLoggerName, location, w, groups)
}

// NewSplitPipelineConfig creates a new split pipeline config
func NewSplitPipelineConfig(showLoggerName bool, location *time.Location, tWidth int, groups GroupContainer) SplitPipelineConfig {
	tWidth -= len(TimeFormat)
	charactersPerPipe := tWidth/groups.GroupCount() - (len(PipelineSeparator) * (groups.GroupCount() - 1))

	return SplitPipelineConfig{
		ShowLoggerName:   showLoggerName,
		Location:         location,
		TerminalWidth:    tWidth,
		GroupContainer:   groups,
		CharacterPerPipe: charactersPerPipe,
	}
}
