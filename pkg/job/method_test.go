// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/shoenig/test/must"
)

func TestMethodSelectorConfig_Validate(t *testing.T) {
	testCases := []struct {
		name                      string
		inputMethodSelectorConfig *MethodSelectorConfig
		expectedError             bool
	}{
		{
			name:                      "nil config",
			inputMethodSelectorConfig: nil,
			expectedError:             true,
		},
		{
			name:                      "nil base config",
			inputMethodSelectorConfig: &MethodSelectorConfig{},
			expectedError:             true,
		},
		{
			name:                      "empty name and unknown provider",
			inputMethodSelectorConfig: &MethodSelectorConfig{MethodSelectorBaseConfig: &MethodSelectorBaseConfig{}},
			expectedError:             true,
		},
		{
			name: "valid with expr provider",
			inputMethodSelectorConfig: &MethodSelectorConfig{
				MethodSelectorBaseConfig: &MethodSelectorBaseConfig{
					Name:     "test-selector",
					Provider: MethodSelectorProviderExpr,
				},
			},
			expectedError: false,
		},
		{
			name: "valid with filter provider",
			inputMethodSelectorConfig: &MethodSelectorConfig{
				MethodSelectorBaseConfig: &MethodSelectorBaseConfig{
					Name:     "test-selector",
					Provider: MethodSelectorProviderFilter,
				},
			},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputMethodSelectorConfig.Validate()
			if tc.expectedError {
				must.Error(t, actualError)
			} else {
				must.NoError(t, actualError)
			}
		})
	}
}

func TestMethodSelectorRunRequest_Validate(t *testing.T) {
	testCases := []struct {
		name                          string
		inputMethodSelectorRunRequest *MethodSelectorRunRequest
		expectedError                 bool
	}{
		{
			name:                          "nil request",
			inputMethodSelectorRunRequest: nil,
			expectedError:                 true,
		},
		{
			name:                          "nil job",
			inputMethodSelectorRunRequest: &MethodSelectorRunRequest{},
			expectedError:                 true,
		},
		{
			name:                          "valid",
			inputMethodSelectorRunRequest: &MethodSelectorRunRequest{Job: &api.Job{}},
			expectedError:                 false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputMethodSelectorRunRequest.Validate()
			if tc.expectedError {
				must.Error(t, actualError)
			} else {
				must.NoError(t, actualError)
			}
		})
	}
}
