// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_DefaultConfig(t *testing.T) {
	t.Run("defaults to memory provider", func(t *testing.T) {
		cfg := DefaultConfig()
		must.NotNil(t, cfg)
		must.Eq(t, ProviderMemory, cfg.Provider)
	})
}

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name          string
		inputConfig   *Config
		expectedError error
	}{
		{
			name:          "nil",
			inputConfig:   nil,
			expectedError: errors.New("state config block required"),
		},
		{
			name:          "memory provider",
			inputConfig:   &Config{Provider: ProviderMemory},
			expectedError: nil,
		},
		{
			name:          "empty provider",
			inputConfig:   &Config{},
			expectedError: errors.New("unsupported state provider \"\""),
		},
		{
			name:          "unsupported provider",
			inputConfig:   &Config{Provider: "bad"},
			expectedError: errors.New("unsupported state provider \"bad\""),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputConfig.Validate()
			if tc.expectedError != nil {
				must.Error(t, actualError)
				must.Eq(t, tc.expectedError.Error(), actualError.Error())
			} else {
				must.NoError(t, actualError)
			}
		})
	}
}

func Test_ConfigMerge(t *testing.T) {
	testCases := []struct {
		name                 string
		inputConfigA         *Config
		inputConfigB         *Config
		expectedOutputConfig *Config
	}{
		{
			name:                 "both non-nil, both have provider",
			inputConfigA:         &Config{Provider: ProviderMemory},
			inputConfigB:         &Config{Provider: ProviderMemory},
			expectedOutputConfig: &Config{Provider: ProviderMemory},
		},
		{
			name:                 "a nil, b has provider",
			inputConfigA:         nil,
			inputConfigB:         &Config{Provider: ProviderMemory},
			expectedOutputConfig: &Config{Provider: ProviderMemory},
		},
		{
			name:                 "a has provider, b nil",
			inputConfigA:         &Config{Provider: ProviderMemory},
			inputConfigB:         nil,
			expectedOutputConfig: &Config{Provider: ProviderMemory},
		},
		{
			name:                 "both nil",
			inputConfigA:         nil,
			inputConfigB:         nil,
			expectedOutputConfig: nil,
		},
		{
			name:                 "b provider overrides a",
			inputConfigA:         &Config{Provider: ProviderMemory},
			inputConfigB:         &Config{Provider: ProviderMemory},
			expectedOutputConfig: &Config{Provider: ProviderMemory},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualConfig := tc.inputConfigA.Merge(tc.inputConfigB)
			must.Eq(t, tc.expectedOutputConfig, actualConfig)
		})
	}

	t.Run("merge does not mutate original config", func(t *testing.T) {
		a := &Config{Provider: "a"}
		b := &Config{Provider: "b"}
		result := a.Merge(b)
		must.Eq(t, "b", result.Provider)
		must.Eq(t, "a", a.Provider)
	})
}
