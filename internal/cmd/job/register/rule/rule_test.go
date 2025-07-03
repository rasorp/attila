// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/pkg/api"
)

func Test_formatRegionPicker(t *testing.T) {

	testCases := []struct {
		name            string
		inputRulePicker *api.JobRegisterRulePicker
		expectedOutput  string
	}{
		{
			name:            "nil picker",
			inputRulePicker: nil,
			expectedOutput:  "",
		},
		{
			name:            "nil expression",
			inputRulePicker: &api.JobRegisterRulePicker{},
			expectedOutput:  "",
		},
		{
			name: "populated expression selector",
			inputRulePicker: &api.JobRegisterRulePicker{
				Expression: &api.JobRegisterRuleFilterExpression{
					Selector: "this.expression",
				},
			},
			expectedOutput: "this.expression",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := formatRegionPicker(tc.inputRulePicker)
			must.Eq(t, tc.expectedOutput, actualOutput)
		})
	}
}
