// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestLimitPicker_SetConfig(t *testing.T) {
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
					Provider: "limit",
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
			name: "negative num",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "limit",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"num": -1},
			},
			outputError: `limit config option "num" cannot be negative`,
		},
		{
			name: "num zero",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "limit",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"num": 0},
			},
		},
		{
			name: "num positive",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "limit",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"num": 3},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.RegionPickerConfig{
				RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
					Provider: "limit",
					Name:     "test",
				},
				ProviderConfig: nil,
			},
			outputError: `limit config required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			strategy := &LimitPicker{}
			err := strategy.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}

func TestLimitPicker_Name(t *testing.T) {

	t.Run("before config", func(t *testing.T) {
		strategy := &LimitPicker{}
		must.Eq(t, "", strategy.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "limit",
				Name:     "my-limit",
			},
			ProviderConfig: map[string]any{"num": 2},
		}
		strategy := &LimitPicker{}
		must.NoError(t, strategy.SetConfig(cfg))
		must.Eq(t, "my-limit", strategy.Name())
	})
}

func TestLimitPicker_Run(t *testing.T) {

	t.Run("num zero returns empty slice not nil", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "limit",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"num": 0},
		}
		strategy := &LimitPicker{}
		must.NoError(t, strategy.SetConfig(cfg))

		req := &jobsdk.RegionPickerRunRequest{
			RegionCandidates: []jobsdk.RegisterRuleRegionCandidate{{Name: "a"}},
		}

		result, err := strategy.Run(req)
		must.NoError(t, err)
		must.Eq(t, 0, len(result.RegionCandidates))
	})

	t.Run("num equal returns all candidates", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "limit",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"num": 2},
		}
		strategy := &LimitPicker{}
		must.NoError(t, strategy.SetConfig(cfg))

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
		}
		req := &jobsdk.RegionPickerRunRequest{
			RegionCandidates: input,
		}

		result, err := strategy.Run(req)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "eu", result.RegionCandidates[0].Group)
		must.Eq(t, "b", result.RegionCandidates[1].Name)
		must.Eq(t, "us", result.RegionCandidates[1].Group)
	})

	t.Run("num less returns first n candidates in order", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "limit",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"num": 2},
		}
		strategy := &LimitPicker{}
		must.NoError(t, strategy.SetConfig(cfg))

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
			{Name: "c", Group: "ap"},
		}
		req := &jobsdk.RegionPickerRunRequest{
			RegionCandidates: input,
		}

		result, err := strategy.Run(req)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "b", result.RegionCandidates[1].Name)
	})

	t.Run("num greater returns all candidates", func(t *testing.T) {
		cfg := &jobsdk.RegionPickerConfig{
			RegionPickerBaseConfig: &jobsdk.RegionPickerBaseConfig{
				Provider: "limit",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"num": 10},
		}
		strategy := &LimitPicker{}
		must.NoError(t, strategy.SetConfig(cfg))

		input := []jobsdk.RegisterRuleRegionCandidate{
			{Name: "a", Group: "eu"},
			{Name: "b", Group: "us"},
		}
		req := &jobsdk.RegionPickerRunRequest{
			RegionCandidates: input,
		}

		result, err := strategy.Run(req)
		must.NoError(t, err)
		must.Len(t, 2, result.RegionCandidates)
		must.Eq(t, "a", result.RegionCandidates[0].Name)
		must.Eq(t, "b", result.RegionCandidates[1].Name)
	})
}
