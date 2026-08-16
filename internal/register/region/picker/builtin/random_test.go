// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestRandomPicker_SetConfig(t *testing.T) {
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
					Provider: "random",
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
			name: "negative seed",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "random",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"seed": -1},
			},
			outputError: `random config option "seed" cannot be negative`,
		},
		{
			name: "seed zero",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "random",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"seed": 0},
			},
		},
		{
			name: "seed positive",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "random",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"seed": 42},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "random",
					Name:     "test",
				},
				ProviderConfig: nil,
			},
			outputError: `random config required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &RandomPicker{}
			err := strategy.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}

func TestRandomPicker_Name(t *testing.T) {

	t.Run("before config", func(t *testing.T) {
		strategy := &RandomPicker{}
		must.Eq(t, "", strategy.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "my-random",
			},
			ProviderConfig: map[string]any{"seed": 42},
		}
		strategy := &RandomPicker{}
		must.NoError(t, strategy.SetConfig(cfg))
		must.Eq(t, "my-random", strategy.Name())
	})
}

func TestRandomPicker_Run(t *testing.T) {

	t.Run("with seed returns deterministic result", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"seed": 12345},
		}

		first := runRandom(t, cfg, []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"},
		})
		second := runRandom(t, cfg, []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"},
		})

		must.Eq(t, first.RegionCandidates, second.RegionCandidates)
	})

	t.Run("with seed returns shuffled candidates", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"seed": 12345},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a"}, {Name: "b"}, {Name: "c"}, {Name: "d"},
		}
		result := runRandom(t, cfg, input)

		must.Len(t, 4, result.RegionCandidates)

		// With seed 12345 the order should not be a,b,c,d  as that would be
		// unshuffled.
		notAllInOrder := false
		for i, c := range result.RegionCandidates {
			if c.Name != input[i].Name {
				notAllInOrder = true
				break
			}
		}
		must.True(t, notAllInOrder)
	})

	t.Run("with seed zero still works", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"seed": 0},
		}

		input := []jobsdk.RegisterRuleRegionCandidate{{Name: "a"}, {Name: "b"}, {Name: "c"}}
		result := runRandom(t, cfg, input)

		must.Len(t, 3, result.RegionCandidates)
	})

	t.Run("empty candidates", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"seed": 12345},
		}

		result := runRandom(t, cfg, []jobsdk.RegisterRuleRegionCandidate{})
		must.Len(t, 0, result.RegionCandidates)
	})

	t.Run("single candidate returns single", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "random",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"seed": 12345},
		}

		result := runRandom(t, cfg, []jobsdk.RegisterRuleRegionCandidate{{Name: "a"}})

		must.Len(t, 1, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
	})
}

// runRandom is a helper that creates a new strategy, applies config, and runs.
func runRandom(t *testing.T, cfg *jobsdk.RegionPickerConfig, candidates []jobsdk.RegisterRuleRegionCandidate) *jobsdk.RegionPickerRunResult {
	t.Helper()
	strategy := &RandomPicker{}
	must.NoError(t, strategy.SetConfig(cfg))

	result, err := strategy.Run(&jobsdk.RegionPickerRunRequest{
		RegionCandidates: candidates,
	})
	must.NoError(t, err)
	return result
}
