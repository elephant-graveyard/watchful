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

package services

import (
	"github.com/homeport/gonvenience/pkg/v1/bunt"
	"github.com/homeport/pina-golada/pkg/files"
	"github.com/homeport/watchful/internal/watchful/assets"
	"github.com/homeport/watchful/pkg/logger"
	"os"
	"path/filepath"
)

var (
	// SampleAppSubPath is the sub-path of the export folder in which the sample app is located
	SampleAppSubPath = "sample-app"
)

// AssetService defines the asset services that is responsible for exporting
type AssetService struct {
	SampleAppLanguage string
	Logger            logger.Logger
	ExportPath        string
}

// Cleanup cleans the asset directory
func (e *AssetService) Cleanup() {
	if err := os.RemoveAll(e.ExportPath); err != nil {
		e.Logger.WriteString(logger.Error, bunt.Sprintf("Red{Could not clean asset directory:} %s", err.Error()))
		return
	}
	e.Logger.WriteString(logger.Info, bunt.Sprintf("Cleaned asset directory"))
}

// NewAssetService creates a new asset services
func NewAssetService(sampleAppLanguage string, logger logger.Logger, ExportPath string) *AssetService {
	return &AssetService{
		SampleAppLanguage: sampleAppLanguage,
		Logger:            logger,
		ExportPath:        filepath.Join(".", ExportPath),
	}
}

// Execute exports all the sample apps
func (e *AssetService) Execute() error {
	directory, err := assets.SampleApp(e.SampleAppLanguage)
	if err != nil {
		return err
	}

	if err := files.WriteToDisk(directory, e.ExportPath, true); err != nil {
		e.Logger.WriteString(logger.Error, bunt.Sprintf("Red{could not export sample app}"))
		return err
	}

	originalAssetName := filepath.Join(e.ExportPath, directory.Name().String())
	newAssetName := filepath.Join(e.ExportPath, SampleAppSubPath)

	if err := os.Rename(originalAssetName, newAssetName); err != nil {
		return err
	}

	e.Logger.WriteString(logger.Info, bunt.Sprintf("SpringGreen{Exported sample-app to %s}", newAssetName))
	return nil
}
