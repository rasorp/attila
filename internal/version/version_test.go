// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		name            string
		inputVersion    string
		inputPreRelease string
		expectedOutput  string
	}{
		{
			name:            "pre-release",
			inputVersion:    "0.1.0",
			inputPreRelease: "beta.1",
			expectedOutput:  "0.1.0-beta.1",
		},
		{
			name:            "release",
			inputVersion:    "0.1.0",
			inputPreRelease: "",
			expectedOutput:  "0.1.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			version = tc.inputVersion
			versionPrerelease = tc.inputPreRelease

			must.Eq(t, tc.expectedOutput, Get())
		})
	}
}
