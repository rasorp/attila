// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestJobRegisterRule_Validate(t *testing.T) {

	testCases := []struct {
		name        string
		inputRule   *JobRegisterRule
		expectedErr bool
	}{
		{
			name:        "nil rule",
			inputRule:   nil,
			expectedErr: true,
		},
		{
			name: "valid rule",
			inputRule: &JobRegisterRule{
				RegionPicker: &JobRegisterRulePicker{
					Expression: &JobRegisterRuleFilterExpression{
						Selector: "filter(regions, .Group == \"europe\")",
					},
				},
			},
			expectedErr: false,
		},
		{
			name: "invalid region picker expression",
			inputRule: &JobRegisterRule{
				RegionPicker: &JobRegisterRulePicker{
					Expression: &JobRegisterRuleFilterExpression{
						Selector: "any(region_namespace, {.Name == \"platform\"})",
					},
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.inputRule.Validate()
			if tc.expectedErr {
				must.Error(t, actualOutput)
			} else {
				must.NoError(t, actualOutput)
			}
		})
	}
}

func TestJobRegisterRulePicker_Validate(t *testing.T) {

	testCases := []struct {
		name            string
		inputRulePicker *JobRegisterRulePicker
		expectedErr     bool
	}{
		{
			name:            "nil picker",
			inputRulePicker: nil,
			expectedErr:     false,
		},
		{
			name:            "nil expression",
			inputRulePicker: &JobRegisterRulePicker{},
			expectedErr:     false,
		},
		{
			name: "invalid expression selector",
			inputRulePicker: &JobRegisterRulePicker{
				Expression: &JobRegisterRuleFilterExpression{
					Selector: "",
				},
			},
			expectedErr: true,
		},
		{
			name: "valid expression selector",
			inputRulePicker: &JobRegisterRulePicker{
				Expression: &JobRegisterRuleFilterExpression{
					Selector: "filter(regions, .Group == \"europe\")",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.inputRulePicker.Validate()
			if tc.expectedErr {
				must.Error(t, actualOutput)
			} else {
				must.NoError(t, actualOutput)
			}
		})
	}
}
