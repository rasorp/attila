// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/pointer"
)

func TestDefaultConfig(t *testing.T) {
	defaultConfig := DefaultConfig()
	must.NotNil(t, defaultConfig)
	must.Nil(t, defaultConfig.Memory)
	must.Nil(t, defaultConfig.File)
}

func TestConfig_Validate(t *testing.T) {

	testCases := []struct {
		name          string
		inputConfig   *Config
		expectedError error
	}{
		{
			name:          "nil config",
			inputConfig:   nil,
			expectedError: errors.New("state config block required"),
		},
		{
			name: "memory enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
			},
			expectedError: nil,
		},
		{
			name: "no backend enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(false),
				},
				File: &FileConfig{
					Enable: pointer.Of(false),
				},
			},
			expectedError: errors.New("no state backend enabled"),
		},
		{
			name: "all backends enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
				File: &FileConfig{
					Enable: pointer.Of(true),
				},
			},
			expectedError: errors.New("only one storage backend can be enabled"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputConfig.Validate()

			if tc.expectedError != nil {
				must.ErrorContains(t, actualError, tc.expectedError.Error())
			} else {
				must.NoError(t, actualError)
			}
		})
	}
}

func TestConfig_Merge(t *testing.T) {

	testCases := []struct {
		name           string
		inputConfig    *Config
		mergeConfig    *Config
		expectedOutput *Config
	}{
		{
			name:           "both nil",
			inputConfig:    nil,
			mergeConfig:    nil,
			expectedOutput: nil,
		},
		{
			name:        "input nil",
			inputConfig: nil,
			mergeConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
			},
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
			},
		},
		{
			name: "merge nil",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
			},
			mergeConfig: nil,
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
			},
		},
		{
			name: "full merge",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(false),
				},
				File: &FileConfig{
					Enable: pointer.Of(false),
					Path:   "/my/path",
				},
			},
			mergeConfig: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
				File: &FileConfig{
					Enable: pointer.Of(true),
					Path:   "/my/new/path",
				},
			},
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: pointer.Of(true),
				},
				File: &FileConfig{
					Enable: pointer.Of(true),
					Path:   "/my/new/path",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			must.Eq(t, tc.expectedOutput, tc.inputConfig.Merge(tc.mergeConfig))
		})
	}
}

func TestMemoryConfig_Enabled(t *testing.T) {

	testCases := []struct {
		name              string
		inputMemoryConfig *MemoryConfig
		expectedOutput    bool
	}{
		{
			name:              "config nil",
			inputMemoryConfig: nil,
			expectedOutput:    false,
		},
		{
			name:              "enabled nil",
			inputMemoryConfig: &MemoryConfig{},
			expectedOutput:    false,
		},
		{
			name:              "enabled false",
			inputMemoryConfig: &MemoryConfig{Enable: pointer.Of(false)},
			expectedOutput:    false,
		},
		{
			name:              "enabled true",
			inputMemoryConfig: &MemoryConfig{Enable: pointer.Of(true)},
			expectedOutput:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			must.Eq(t, tc.expectedOutput, tc.inputMemoryConfig.Enabled())
		})
	}
}

func TestFileConfig_Enabled(t *testing.T) {

	testCases := []struct {
		name            string
		inputFileConfig *FileConfig
		expectedOutput  bool
	}{
		{
			name:            "config nil",
			inputFileConfig: nil,
			expectedOutput:  false,
		},
		{
			name:            "enabled nil",
			inputFileConfig: &FileConfig{},
			expectedOutput:  false,
		},
		{
			name:            "enabled false",
			inputFileConfig: &FileConfig{Enable: pointer.Of(false)},
			expectedOutput:  false,
		},
		{
			name:            "enabled true",
			inputFileConfig: &FileConfig{Enable: pointer.Of(true)},
			expectedOutput:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			must.Eq(t, tc.expectedOutput, tc.inputFileConfig.Enabled())
		})
	}
}

func TestFileConfig_Validate(t *testing.T) {

	testCases := []struct {
		name            string
		inputFileConfig *FileConfig
		expectedError   bool
	}{
		{
			name: "not enabled",
			inputFileConfig: &FileConfig{
				Enable: pointer.Of(false),
			},
			expectedError: false,
		},
		{
			name: "empty path",
			inputFileConfig: &FileConfig{
				Enable: pointer.Of(true),
				Path:   "",
			},
			expectedError: true,
		},
		{
			name: "not absolute path",
			inputFileConfig: &FileConfig{
				Enable: pointer.Of(true),
				Path:   "~/jrasell",
			},
			expectedError: true,
		},
		{
			name: "non-existent path",
			inputFileConfig: &FileConfig{
				Enable: pointer.Of(true),
				Path:   "/jrasell/data",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.inputFileConfig.Validate()
			if tc.expectedError {
				must.Error(t, actualOutput)
			} else {
				must.NoError(t, actualOutput)
			}
		})
	}
}
