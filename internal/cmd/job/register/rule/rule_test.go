// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/pkg/api"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

func Test_formatRegionPicker(t *testing.T) {

	testCases := []struct {
		name            string
		inputRulePicker []*api.JobRegisterRegionPicker
		expectedOutput  string
	}{
		{
			name:            "nil picker",
			inputRulePicker: nil,
			expectedOutput:  "",
		},
		{
			name:            "empty strategies",
			inputRulePicker: []*api.JobRegisterRegionPicker{},
			expectedOutput:  "",
		},
		{
			name: "populated strategies",
			inputRulePicker: []*api.JobRegisterRegionPicker{
				{Name: "foo", Provider: jobsdk.RegionPickerProviderFilter},
				{Name: "bar", Provider: jobsdk.RegionPickerProviderFilter},
			},
			expectedOutput: "filter::foo, filter::bar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := formatRegionPicker(tc.inputRulePicker)
			must.Eq(t, tc.expectedOutput, actualOutput)
		})
	}
}
