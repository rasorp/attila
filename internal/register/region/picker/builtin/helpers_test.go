// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/pkg/job"
)

func TestDecodeParams(t *testing.T) {
	testCases := []struct {
		name          string
		inputRaw      map[string]any
		inputDst      any
		expectedError bool
	}{
		{
			name:          "nil raw",
			inputRaw:      nil,
			inputDst:      new(map[string]string),
			expectedError: false,
		},
		{
			name:          "empty raw",
			inputRaw:      map[string]any{},
			inputDst:      new(map[string]string),
			expectedError: false,
		},
		{
			name:     "simple string map",
			inputRaw: map[string]any{"foo": "bar"},
			inputDst: &struct {
				Foo string `mapstructure:"foo"`
			}{},
			expectedError: false,
		},
		{
			name:     "nested struct",
			inputRaw: map[string]any{"outer": map[string]any{"inner": "value"}},
			inputDst: &struct {
				Outer *struct {
					Inner string `mapstructure:"inner"`
				} `mapstructure:"outer"`
			}{},
			expectedError: false,
		},
		{
			name:     "type mismatch returns error",
			inputRaw: map[string]any{"count": "not a number"},
			inputDst: &struct {
				Count int `mapstructure:"count"`
			}{},
			expectedError: true,
		},
		{
			name:     "missing field ignored by default",
			inputRaw: map[string]any{"foo": "bar"},
			inputDst: &struct {
				Baz string `mapstructure:"baz"`
			}{},
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := decodeParams(tc.inputRaw, tc.inputDst)
			if tc.expectedError {
				must.Error(t, err)
			} else {
				must.NoError(t, err)
			}
		})
	}
}

func TestCopyCandidates(t *testing.T) {
	testCases := []struct {
		name  string
		input []job.RegisterRuleRegionCandidate
	}{
		{
			name:  "nil input",
			input: nil,
		},
		{
			name:  "empty slice",
			input: []job.RegisterRuleRegionCandidate{},
		},
		{
			name: "single candidate",
			input: []job.RegisterRuleRegionCandidate{
				{Name: "eu-west-1", Group: "eu"},
			},
		},
		{
			name: "multiple candidates",
			input: []job.RegisterRuleRegionCandidate{
				{Name: "eu-west-1", Group: "eu"},
				{Name: "us-east-1", Group: "us"},
				{Name: "ap-south-1", Group: "ap"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actual := CopyCandidates(tc.input)

			if len(tc.input) > 0 {
				must.Eq(t, tc.input, actual)
			} else {
				must.SliceEmpty(t, actual)
			}
		})
	}
}

func Test_CopyCandidates(t *testing.T) {
	testCases := []struct {
		name      string
		input     []job.RegisterRuleRegionCandidate
		expectNil bool
	}{
		{
			name:      "nil slice becomes empty slice",
			input:     nil,
			expectNil: false,
		},
		{
			name:      "empty slice",
			input:     []job.RegisterRuleRegionCandidate{},
			expectNil: false,
		},
		{
			name:      "single candidate",
			input:     []job.RegisterRuleRegionCandidate{{Name: "eu-west-1", Group: "eu"}},
			expectNil: false,
		},
		{
			name:      "multiple candidates",
			input:     []job.RegisterRuleRegionCandidate{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			expectNil: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := CopyCandidates(tc.input)
			must.Eq(t, len(tc.input), len(actualOutput))
			if len(tc.input) > 0 {
				must.Eq(t, tc.input, actualOutput)
			}
			if tc.expectNil {
				must.Nil(t, actualOutput)
			} else {
				must.NotNil(t, actualOutput)
			}
		})
	}
}
