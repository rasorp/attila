// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"fmt"
	"testing"

	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestJobRegisterMethodValidate(t *testing.T) {
	testCases := []struct {
		name           string
		inputMethod    *JobRegisterMethod
		expectedErrors []string
	}{
		{
			name: "one rule and one selector",
			inputMethod: &JobRegisterMethod{
				Name: "test-method",
				Rules: []*JobRegisterMethodRuleLink{
					{Name: "my-rule"},
				},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
						Name:     "my-selector",
						Provider: jobsdk.MethodSelectorProviderExpr,
					}},
				},
			},
			expectedErrors: nil,
		},
		{
			name: "multiple rules and selectors",
			inputMethod: &JobRegisterMethod{
				Name: "test-method",
				Rules: []*JobRegisterMethodRuleLink{
					{Name: "rule-a"},
					{Name: "rule-b"},
					{Name: "rule-c"},
				},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel-1", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel-2", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: nil,
		},
		{
			name: "no rules",
			inputMethod: &JobRegisterMethod{
				Name:  "test-method",
				Rules: nil,
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: []string{"at least one rule required"},
		},
		{
			name: "empty rules",
			inputMethod: &JobRegisterMethod{
				Name:  "test-method",
				Rules: []*JobRegisterMethodRuleLink{},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: []string{"at least one rule required"},
		},
		{
			name: "invalid rule name",
			inputMethod: &JobRegisterMethod{
				Name:  "test-method",
				Rules: []*JobRegisterMethodRuleLink{{Name: ""}},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: []string{"rule 0; method rule \"name\" cannot be empty"},
		},
		{
			name: "valid and invalid rules",
			inputMethod: &JobRegisterMethod{
				Name: "test-method",
				Rules: []*JobRegisterMethodRuleLink{
					{Name: "good"},
					{Name: ""},
					{Name: "also-good"},
					{Name: ""},
				},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "sel", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: []string{"rule 1; method rule \"name\" cannot be empty", "rule 3; method rule \"name\" cannot be empty"},
		},
		{
			name: "no selectors",
			inputMethod: &JobRegisterMethod{
				Name:      "test-method",
				Rules:     []*JobRegisterMethodRuleLink{{Name: "rule"}},
				Selectors: nil,
			},
			expectedErrors: []string{"at least one selector required"},
		},
		{
			name: "empty selectors",
			inputMethod: &JobRegisterMethod{
				Name:      "test-method",
				Rules:     []*JobRegisterMethodRuleLink{{Name: "rule"}},
				Selectors: []*jobsdk.MethodSelectorConfig{},
			},
			expectedErrors: []string{"at least one selector required"},
		},
		{
			name: "invalid selector name",
			inputMethod: &JobRegisterMethod{
				Name:  "test-method",
				Rules: []*JobRegisterMethodRuleLink{{Name: "rule"}},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "", Provider: jobsdk.MethodSelectorProviderExpr,
						},
					},
				},
			},
			expectedErrors: []string{`selector 0; method selector "name" cannot be empty`},
		},
		{
			name: "invalid selector provider",
			inputMethod: &JobRegisterMethod{
				Name:  "test-method",
				Rules: []*JobRegisterMethodRuleLink{{Name: "rule"}},
				Selectors: []*jobsdk.MethodSelectorConfig{{
					MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{Name: "bad-provider", Provider: "bogus"},
				}},
			},
			expectedErrors: []string{`selector 0; unsupported method selector provider "bogus"`},
		},
		{
			name: "invalid rule and selectors",
			inputMethod: &JobRegisterMethod{
				Name:  "broken-method",
				Rules: []*JobRegisterMethodRuleLink{{Name: ""}},
				Selectors: []*jobsdk.MethodSelectorConfig{
					{
						MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
							Name: "", Provider: "bogus",
						},
					},
				},
			},
			expectedErrors: []string{
				`method rule "name" cannot be empty`,
				`method selector "name" cannot be empty`,
				`unsupported method selector provider "bogus"`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			actualError := tc.inputMethod.Validate()

			if tc.expectedErrors == nil {
				must.NoError(t, actualError)
			} else {
				fmt.Println(actualError)
				must.Error(t, actualError)
				for _, wantSubstr := range tc.expectedErrors {
					must.ErrorContains(t, actualError, wantSubstr)
				}
			}
		})
	}
}
