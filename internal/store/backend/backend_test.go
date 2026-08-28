// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_New(t *testing.T) {
	testCases := []struct {
		name          string
		inputConfig   *Config
		expectedError error
		expectedName  string
	}{
		{
			name:          "memory backend",
			inputConfig:   &Config{Provider: ProviderMemory},
			expectedError: nil,
			expectedName:  "mem",
		},
		{
			name:          "no backend",
			inputConfig:   &Config{},
			expectedError: errors.New("no state backend configured"),
			expectedName:  "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualBackend, actualError := New(tc.inputConfig)
			if tc.expectedError != nil {
				must.Nil(t, actualBackend)
				must.Error(t, actualError)
			} else {
				must.NotNil(t, actualBackend)
				must.NoError(t, actualError)
				must.Eq(t, tc.expectedName, tc.expectedName)
			}
		})
	}
}
