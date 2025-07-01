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
	must.NotNil(t, defaultConfig.Memory)
	must.True(t, *defaultConfig.Memory.Enabled)
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
					Enabled: pointer.Of(true),
				},
			},
			expectedError: nil,
		},
		{
			name: "memory not enabled",
			inputConfig: &Config{
				Memory: &MemoryConfig{
					Enabled: pointer.Of(false),
				},
			},
			expectedError: errors.New("memory state must be enabled"),
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
