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

package cmd

import (
	"fmt"
	"os"

	"github.com/gonvenience/term"
	"github.com/homeport/watchful/internal/watchful/services"
	"github.com/spf13/cobra"
)

var (
	// TerminalWidth is an argument that one can pass to override the found terminal width
	TerminalWidth int

	// ConfigContent represents the string content of the config file
	ConfigContent string

	// PushedAppSampleLanguage is the language in which the pushed sample app should be written
	PushedAppSampleLanguage string

	// Verbose defines if the program will output debug information
	Verbose bool
)

// runCmd is the run command definition using cobra
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run is the sub-command that executes watchful",
	Long:  "Run executes watchful with the provided configuration and merkhet targets.",
	Run:   run,
}

// run is the method called when someone uses the watchful run command
func run(cmd *cobra.Command, args []string) {
	e := services.MainService{
		TerminalWidth:           TerminalWidth,
		ConfigContent:           ConfigContent,
		PushedAppSampleLanguage: PushedAppSampleLanguage,
		Verbose:                 Verbose,
	}

	if err := e.Execute(); err != nil {
		fmt.Println(fmt.Sprintf("watchful failed: %s", err.Error()))
		os.Exit(1)
	}

	os.Exit(0) // At this point nothing is allowed to run and we force exit
}

// init adds the runCmd to the watchful root command
func init() {
	runCmd.PersistentFlags().IntVarP(&TerminalWidth, "terminalWidth", "w", term.GetTerminalWidth(), "Provides the terminal width")
	runCmd.PersistentFlags().StringVarP(&ConfigContent, "config", "c", "", "Provides the configuration for watchful")
	runCmd.PersistentFlags().StringVarP(&PushedAppSampleLanguage, "language", "l", "go", "Defines in which language "+
		"the push sample app should be written")
	runCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Toggles whether the app is run in verbose mode")
	rootCmd.AddCommand(runCmd)
}
