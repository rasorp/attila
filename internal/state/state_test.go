// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/pointer"
)

func TestNewBackend(t *testing.T) {
	testCases := []struct {
		name          string
		inputConfig   *Config
		expectedError error
		expectedName  string
	}{
		{
			name: "memory backend",
			inputConfig: &Config{
				Memory: &MemoryConfig{Enable: pointer.Of(true)},
			},
			expectedError: nil,
			expectedName:  "mem",
		},
		{
			name: "file backend",
			inputConfig: &Config{
				File: &FileConfig{
					Enable: pointer.Of(true),
					Path:   "/tmp/attila-" + strconv.FormatInt(time.Now().UnixNano(), 10),
				},
			},
			expectedError: nil,
			expectedName:  "file",
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

			stateBackend, actualErr := NewBackend(tc.inputConfig)

			if tc.expectedError != nil {
				must.ErrorContains(t, actualErr, tc.expectedError.Error())
				must.Nil(t, stateBackend)
			} else {
				must.NoError(t, actualErr)
				must.Eq(t, tc.expectedName, stateBackend.Name())
			}
		})
	}
}
