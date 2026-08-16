// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
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
					Enable: new(true),
				},
			},
			expectedError: nil,
		},
		{
			name: "no backend enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: new(false),
				},
				File: &FileConfig{
					Enable: new(false),
				},
			},
			expectedError: errors.New("no state backend enabled"),
		},
		{
			name: "all backends enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
				File: &FileConfig{
					Enable: new(true),
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
					Enable: new(true),
				},
			},
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
			},
		},
		{
			name: "merge nil",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
			},
			mergeConfig: nil,
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
			},
		},
		{
			name: "full merge",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enable: new(false),
				},
				File: &FileConfig{
					Enable: new(false),
					Path:   "/my/path",
				},
			},
			mergeConfig: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
				File: &FileConfig{
					Enable: new(true),
					Path:   "/my/new/path",
				},
			},
			expectedOutput: &Config{
				Memory: &MemoryConfig{
					Enable: new(true),
				},
				File: &FileConfig{
					Enable: new(true),
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
			inputMemoryConfig: &MemoryConfig{Enable: new(false)},
			expectedOutput:    false,
		},
		{
			name:              "enabled true",
			inputMemoryConfig: &MemoryConfig{Enable: new(true)},
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
			inputFileConfig: &FileConfig{Enable: new(false)},
			expectedOutput:  false,
		},
		{
			name:            "enabled true",
			inputFileConfig: &FileConfig{Enable: new(true)},
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
				Enable: new(false),
			},
			expectedError: false,
		},
		{
			name: "empty path",
			inputFileConfig: &FileConfig{
				Enable: new(true),
				Path:   "",
			},
			expectedError: true,
		},
		{
			name: "not absolute path",
			inputFileConfig: &FileConfig{
				Enable: new(true),
				Path:   "~/jrasell",
			},
			expectedError: true,
		},
		{
			name: "non-existent path",
			inputFileConfig: &FileConfig{
				Enable: new(true),
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
