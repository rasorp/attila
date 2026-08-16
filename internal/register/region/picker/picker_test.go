// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package picker

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/pkg/api"
	"github.com/rasorp/attila/pkg/job"
)

func Test_New(t *testing.T) {
	testCases := []struct {
		name               string
		inputPickerConfigs []*job.RegionPickerConfig
		expectedProviders  int
		expectedError      bool
	}{
		{
			name: "single valid expr picker",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-expr", job.RegionPickerProviderExpr, map[string]any{"expression": "regions"}),
			},
			expectedProviders: 1,
		},
		{
			name: "single valid filter picker",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-filter", job.RegionPickerProviderFilter, map[string]any{"expression": "region.name == \"eu\""}),
			},
			expectedProviders: 1,
		},
		{
			name: "single valid limit picker",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-limit", job.RegionPickerProviderLimit, map[string]any{"num": 5}),
			},
			expectedProviders: 1,
		},
		{
			name: "single valid random picker",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-random", job.RegionPickerProviderRandom, map[string]any{"seed": int64(42)}),
			},
			expectedProviders: 1,
		},
		{
			name: "multiple valid pickers",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("first", job.RegionPickerProviderFilter, map[string]any{"expression": "region.name == \"eu\""}),
				regionPickerConfig("second", job.RegionPickerProviderLimit, map[string]any{"num": 2}),
			},
			expectedProviders: 2,
		},
		{
			name:               "empty configs",
			inputPickerConfigs: []*job.RegionPickerConfig{},
			expectedProviders:  0,
		},
		{
			name:               "nil configs",
			inputPickerConfigs: nil,
			expectedProviders:  0,
		},
		{
			name: "invalid provider",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("bad-provider", "unknown", map[string]any{}),
			},
			expectedError: true,
		},
		{
			name: "empty name",
			inputPickerConfigs: []*job.RegionPickerConfig{
				{
					RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
						Name:     "",
						Provider: job.RegionPickerProviderExpr,
					},
					ProviderConfig: map[string]any{"expression": "regions"},
				},
			},
			expectedError: true,
		},
		{
			name: "expr missing expression",
			inputPickerConfigs: []*job.RegionPickerConfig{
				{
					RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
						Name:     "test-expr",
						Provider: job.RegionPickerProviderExpr,
					},
					ProviderConfig: map[string]any{},
				},
			},
			expectedError: true,
		},
		{
			name: "filter missing expression",
			inputPickerConfigs: []*job.RegionPickerConfig{
				{
					RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
						Name:     "test-filter",
						Provider: job.RegionPickerProviderFilter,
					},
					ProviderConfig: map[string]any{},
				},
			},
			expectedError: true,
		},
		{
			name: "limit negative num",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-limit", job.RegionPickerProviderLimit, map[string]any{"num": -1}),
			},
			expectedError: true,
		},
		{
			name: "random negative seed",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("test-random", job.RegionPickerProviderRandom, map[string]any{"seed": int64(-1)}),
			},
			expectedError: true,
		},
		{
			name: "mixed valid and invalid",
			inputPickerConfigs: []*job.RegionPickerConfig{
				regionPickerConfig("valid-limit", job.RegionPickerProviderLimit, map[string]any{"num": 3}),
				{
					RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
						Name:     "",
						Provider: "bogus",
					},
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualPicker, actualError := New(tc.inputPickerConfigs)
			if tc.expectedError {
				must.Nil(t, actualPicker)
				must.Error(t, actualError)
				return
			}

			must.NotNil(t, actualPicker)
			must.NoError(t, actualError)
			must.Eq(t, tc.expectedProviders, len(actualPicker.providers))
		})
	}
}

func TestPicker_Process(t *testing.T) {
	testCases := []struct {
		name                  string
		configuredPickers     []*job.RegionPickerConfig
		inputPickerRunRequest *job.RegionPickerRunRequest
		useNilReceiver        bool
		expectedErrorContains string
		assertFn              func(t *testing.T, result []job.RegisterRuleRegionCandidate)
	}{
		{
			name: "nil request returns error",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("test", job.RegionPickerProviderLimit, map[string]any{"num": 5}),
			},
			expectedErrorContains: "empy picker request",
		},
		{
			name:           "nil receiver returns error",
			useNilReceiver: true,
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule:             apiRule("rule-1", apiRegionPicker("test", job.RegionPickerProviderLimit, nil)),
				RegionCandidates: []job.RegisterRuleRegionCandidate{{Name: "a"}},
			},
			expectedErrorContains: "picker is not configured",
		},
		{
			name:              "empty rule pickers returns input candidates",
			configuredPickers: []*job.RegionPickerConfig{},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule:             apiRule("rule-1"),
				RegionCandidates: []job.RegisterRuleRegionCandidate{{Name: "eu-west-1"}},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 1, result)
				must.Eq(t, "eu-west-1", result[0].Name)
			},
		},
		{
			name: "limit picks first n candidates",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 2}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 2}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{
					{Name: "a"},
					{Name: "b"},
					{Name: "c"},
					{Name: "d"},
				},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 2, result)
				must.Eq(t, "a", result[0].Name)
				must.Eq(t, "b", result[1].Name)
			},
		},
		{
			name: "limit zero returns empty slice",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 0}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 0}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{{Name: "a"}},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 0, result)
			},
		},
		{
			name: "limit num larger than input returns all",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 100}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("limiter", job.RegionPickerProviderLimit, map[string]any{"num": 100}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{{Name: "a"}, {Name: "b"}},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 2, result)
				must.Eq(t, "a", result[0].Name)
				must.Eq(t, "b", result[1].Name)
			},
		},
		{
			name: "filter picks candidates matching expression",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("f1", job.RegionPickerProviderFilter, map[string]any{"expression": "region.Group == \"eu\""}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("f1", job.RegionPickerProviderFilter, map[string]any{"expression": "region.Group == \"eu\""}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{
					{Name: "a", Group: "eu"},
					{Name: "b", Group: "us"},
					{Name: "c", Group: "eu"},
				},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 2, result)
				must.Eq(t, "a", result[0].Name)
				must.Eq(t, "c", result[1].Name)
			},
		},
		{
			name: "chained pickers apply sequentially",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("f1", job.RegionPickerProviderFilter, map[string]any{"expression": "region.Group == \"eu\""}),
				regionPickerConfig("l1", job.RegionPickerProviderLimit, map[string]any{"num": 1}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("f1", job.RegionPickerProviderFilter, map[string]any{"expression": "region.Group == \"eu\""}),
					apiRegionPicker("l1", job.RegionPickerProviderLimit, map[string]any{"num": 1}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{
					{Name: "a", Group: "eu"},
					{Name: "b", Group: "us"},
					{Name: "c", Group: "eu"},
				},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 1, result)
				must.Eq(t, "a", result[0].Name)
			},
		},
		{
			name: "picker name not found in registry returns error",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("f1", job.RegionPickerProviderFilter, map[string]any{"expression": "true"}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("unknown", job.RegionPickerProviderFilter, map[string]any{"expression": "true"}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{{Name: "a"}},
			},
			expectedErrorContains: `picker provider "filter" not configured`,
		},
		{
			name: "preserves metadata and context",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("pass-through", job.RegionPickerProviderFilter, map[string]any{"expression": "true"}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("pass-through", job.RegionPickerProviderFilter, map[string]any{"expression": "true"}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{
					{Name: "a", Group: "eu", Metadata: map[string]any{"k": "v"}, Context: map[string]any{"x": 1}},
				},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 1, result)
				must.Eq(t, "a", result[0].Name)
				must.Eq(t, "eu", result[0].Group)
				must.MapEq(t, map[string]any{"k": "v"}, result[0].Metadata)
				must.MapEq(t, map[string]any{"x": 1}, result[0].Context)
			},
		},
		{
			name: "rule is available to picker expressions",
			configuredPickers: []*job.RegionPickerConfig{
				regionPickerConfig("rule-aware", job.RegionPickerProviderFilter, map[string]any{"expression": `rule.Name == "rule-1" && region.Group == "eu"`}),
			},
			inputPickerRunRequest: &job.RegionPickerRunRequest{
				Rule: apiRule(
					"rule-1",
					apiRegionPicker("rule-aware", job.RegionPickerProviderFilter, map[string]any{"expression": `rule.Name == "rule-1" && region.Group == "eu"`}),
				),
				RegionCandidates: []job.RegisterRuleRegionCandidate{
					{Name: "a", Group: "eu"},
					{Name: "b", Group: "us"},
				},
			},
			assertFn: func(t *testing.T, result []job.RegisterRuleRegionCandidate) {
				must.Len(t, 1, result)
				must.Eq(t, "a", result[0].Name)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.useNilReceiver {
				var nilPicker *Picker
				actualResponse, err := nilPicker.Process(tc.inputPickerRunRequest)
				must.Nil(t, actualResponse)
				must.ErrorContains(t, err, tc.expectedErrorContains)
				return
			}

			pickerImpl, err := New(tc.configuredPickers)
			must.NoError(t, err)

			actualResponse, err := pickerImpl.Process(tc.inputPickerRunRequest)
			if tc.expectedErrorContains != "" {
				must.Nil(t, actualResponse)
				must.ErrorContains(t, err, tc.expectedErrorContains)
			} else {
				must.NoError(t, err)
				tc.assertFn(t, actualResponse)
			}
		})
	}
}

func regionPickerConfig(name, provider string, cfg map[string]any) *job.RegionPickerConfig {
	return &job.RegionPickerConfig{
		RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
			Name:     name,
			Provider: provider,
		},
		ProviderConfig: cfg,
	}
}

func apiRule(name string, pickers ...*api.JobRegisterRegionPicker) *job.RegionPickerRule {
	regionPickers := make([]*job.RegionPickerConfig, len(pickers))
	for i, p := range pickers {
		regionPickers[i] = &job.RegionPickerConfig{
			RegionPickerBaseConfig: &job.RegionPickerBaseConfig{
				Name:     p.Name,
				Provider: p.Provider,
			},
			ProviderConfig: p.Config,
		}
	}
	return &job.RegionPickerRule{
		Name:          name,
		RegionPickers: regionPickers,
	}
}

func apiRegionPicker(name, provider string, cfg map[string]any) *api.JobRegisterRegionPicker {
	return &api.JobRegisterRegionPicker{
		Name:     name,
		Provider: provider,
		Config:   cfg,
	}
}
