// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestFilterExprPicker_SetConfig(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         *jobsdk.RegionPickerConfig
		outputError string
	}{
		{
			name:        "nil config",
			cfg:         nil,
			outputError: "region picker config is nil",
		},
		{
			name: "empty name",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "",
				},
			},
			outputError: "\"name\" cannot be empty",
		},
		{
			name: "unsupported provider",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "unknown",
					Name:     "test",
				},
			},
			outputError: `unsupported region picker provider "unknown"`,
		},
		{
			name: "empty selector",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": ""},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "whitespace only selector",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "    "},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "invalid expression",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "[invalid"},
			},
			outputError: `failed to compile filter picker expression`,
		},
		{
			name: "valid expression returns no error",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "region.Name == \"a\""},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: nil,
			},
			outputError: `filter config required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &FilterPicker{}
			err := strategy.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}

func TestFilterExprPicker_Name(t *testing.T) {

	t.Run("before config", func(t *testing.T) {
		strategy := &FilterPicker{}
		must.Eq(t, "", strategy.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "my-filter",
			},
			ProviderConfig: map[string]any{"expression": "region.Name == \"a\""},
		}
		strategy := &FilterPicker{}
		must.NoError(t, strategy.SetConfig(cfg))
		must.Eq(t, "my-filter", strategy.Name())
	})
}

func TestFilterExprPicker_Run(t *testing.T) {

	t.Run("empty candidates returns empty result", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "true"},
		}

		result, err := runFilterExpr(cfg, []jobsdk.RegisterRuleRegionCandidate{})
		must.NoError(t, err)
		must.Len(t, 0, result.RegionCandidates)
	})

	t.Run("no candidates when none match", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "region.Name == \"z\""},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 0, result.RegionCandidates)
	})

	t.Run("matching candidates", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "region.Name == \"b\""},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "b", result.RegionCandidates[0].Name)
	})

	t.Run("preserves group metadata", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "region.Group == \"eu\""},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
			{Name: "c", Group: "eu"},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "eu", result.RegionCandidates[0].Group)
		must.Eq(t, "c", result.RegionCandidates[1].Name)
		must.Eq(t, "eu", result.RegionCandidates[1].Group)
	})

	t.Run("metadata filter", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "candidate.Metadata[\"active\"] == true"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Metadata: map[string]any{"active": true}},
			{Name: "b", Metadata: map[string]any{"active": false}},
			{Name: "c", Metadata: map[string]any{"active": true}},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "c", result.RegionCandidates[1].Name)
	})

	t.Run("context", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "x > 5"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Context: map[string]any{"x": 10}},
			{Name: "b", Context: map[string]any{"x": 3}},
			{Name: "c", Context: map[string]any{"x": 20}},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
	})

	t.Run("all candidates match returns all", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "true"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"},
		}
		result, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
		must.Len(t, 3, result.RegionCandidates)
	})

	t.Run("no candidates match returns nil", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "false"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"},
		}
		_, err := runFilterExpr(cfg, input)
		must.NoError(t, err)
	})
}

func runFilterExpr(cfg *jobsdk.RegionPickerConfig, candidates []jobsdk.RegisterRuleRegionCandidate) (*jobsdk.RegionPickerRunResult, error) {
	strategy := &FilterPicker{}
	if err := strategy.SetConfig(cfg); err != nil {
		panic(err)
	}
	return strategy.Run(&jobsdk.RegionPickerRunRequest{
		RegionCandidates: candidates,
	})
}
