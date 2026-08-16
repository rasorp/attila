// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestCelPicker_SetConfig(t *testing.T) {
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
					Provider: jobsdk.RegionPickerProviderExpr,
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
			name: "empty expression",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: jobsdk.RegionPickerProviderExpr,
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": ""},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "whitespace only expression",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: jobsdk.RegionPickerProviderExpr,
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
					Provider: jobsdk.RegionPickerProviderExpr,
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "regions.filter(r, "},
			},
			outputError: `failed to compile expr expression`,
		},
		{
			name: "valid expression returns no error",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: jobsdk.RegionPickerProviderExpr,
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "regions.filter(r, r.group == \"eu\")"},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: jobsdk.RegionPickerProviderExpr,
					Name:     "test",
				},
				ProviderConfig: nil,
			},
			outputError: `expr config required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &ExprPicker{}
			err := strategy.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}

func TestCelPicker_Name(t *testing.T) {
	t.Run("before config", func(t *testing.T) {
		strategy := &ExprPicker{}
		must.Eq(t, "", strategy.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "my-cel",
			},
			ProviderConfig: map[string]any{"expression": "regions"},
		}
		strategy := &ExprPicker{}
		must.NoError(t, strategy.SetConfig(cfg))
		must.Eq(t, "my-cel", strategy.Name())
	})
}

func TestCelPicker_Run(t *testing.T) {
	t.Run("empty candidates with valid expression", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "regions"},
		}

		result, err := runCEL(cfg, []jobsdk.RegisterRuleRegionCandidate{})
		must.NoError(t, err)
		must.Len(t, 0, result.RegionCandidates)
	})

	t.Run("expression returns explicit single candidate list", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "[regions[0]]"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
		}
		result, err := runCEL(cfg, input)
		must.NoError(t, err)
		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "eu", result.RegionCandidates[0].Group)
	})

	t.Run("filter expression", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "regions.filter(r, r.group == \"eu\")"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
			{Name: "c", Group: "eu"},
		}
		result, err := runCEL(cfg, input)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "c", result.RegionCandidates[1].Name)
	})

	t.Run("expression uses candidates alias", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "candidates.filter(r, r.name == \"b\")"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{{Name: "a"}, {Name: "b"}}
		result, err := runCEL(cfg, input)
		must.NoError(t, err)
		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "b", result.RegionCandidates[0].Name)
	})

	t.Run("expression uses job context", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "regions.filter(r, r.name == job[\"region\"])"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{{Name: "a"}, {Name: "b"}}
		result, err := runCELWithJob(cfg, input, map[string]any{"region": "b"})
		must.NoError(t, err)
		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "b", result.RegionCandidates[0].Name)
	})

	t.Run("preserves metadata and context", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: jobsdk.RegionPickerProviderExpr,
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "regions.filter(r, r.metadata[\"active\"] == true)"},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu", Metadata: map[string]any{"active": true}, Context: map[string]any{"x": 1}},
		}
		result, err := runCEL(cfg, input)
		must.NoError(t, err)
		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "eu", result.RegionCandidates[0].Group)
		must.MapEq(t, input[0].Metadata, result.RegionCandidates[0].Metadata)
		must.MapEq(t, input[0].Context, result.RegionCandidates[0].Context)
	})
}

func runCEL(cfg *jobsdk.RegionPickerConfig, candidates []jobsdk.RegisterRuleRegionCandidate) (*jobsdk.RegionPickerRunResult, error) {
	return runCELWithJob(cfg, candidates, nil)
}

func runCELWithJob(
	cfg *jobsdk.RegionPickerConfig,
	candidates []jobsdk.RegisterRuleRegionCandidate,
	job any,
) (*jobsdk.RegionPickerRunResult, error) {
	strategy := &ExprPicker{}
	if err := strategy.SetConfig(cfg); err != nil {
		return nil, err
	}
	return strategy.Run(&jobsdk.RegionPickerRunRequest{
		Job:              job,
		RegionCandidates: candidates,
	})
}
