// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestFilterSelector_Provider(t *testing.T) {
	selector := &FilterSelector{}
	must.Eq(t, "filter", selector.Provider())
}

func TestFilterSelector_Name(t *testing.T) {

	t.Run("before config", func(t *testing.T) {
		selector := &FilterSelector{}
		must.Eq(t, "", selector.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "my-filter",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		selector := &FilterSelector{}
		must.NoError(t, selector.SetConfig(cfg))
		must.Eq(t, "my-filter", selector.Name())
	})
}

func TestFilterSelector_Run(t *testing.T) {

	nomadJob := api.Job{
		ID:        new("test-job-1"),
		Name:      new("test-job"),
		Namespace: new("platform"),
		Type:      new("batch"),
	}

	t.Run("nil request returns non-match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "true"},
		}
		selector := &FilterSelector{}
		must.NoError(t, selector.SetConfig(cfg))
		result, err := selector.Run(nil)
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("nil job returns non-match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "true"},
		}
		selector := &FilterSelector{}
		must.NoError(t, selector.SetConfig(cfg))
		result, err := selector.Run(&jobsdk.MethodSelectorRunRequest{Job: nil})
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("expression matches true", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		result, err := runFilterSelector(t, cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("expression matches false", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"other\""},
		}
		result, err := runFilterSelector(t, cfg, &nomadJob)
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("job id match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.ID == \"test-job-1\""},
		}
		result, err := runFilterSelector(t, cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("multiple conditions", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\" && job.Type == \"batch\""},
		}
		result, err := runFilterSelector(t, cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("contains check", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "filter",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Name contains \"test\" "},
		}
		result, err := runFilterSelector(t, cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})
}

func runFilterSelector(
	t *testing.T,
	cfg *jobsdk.MethodSelectorConfig,
	job *api.Job,
) (*jobsdk.MethodSelectorRunResult, error) {

	s := &FilterSelector{}
	must.NoError(t, s.SetConfig(cfg))
	return s.Run(&jobsdk.MethodSelectorRunRequest{Job: job})
}

func TestFilterSelector_SetConfig(t *testing.T) {
	testCases := []struct {
		name        string
		cfg         *jobsdk.MethodSelectorConfig
		outputError string
	}{
		{
			name:        "nil config",
			cfg:         nil,
			outputError: "method selector config is nil",
		},
		{
			name: "empty name",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "",
				},
			},
			outputError: "\"name\" cannot be empty",
		},
		{
			name: "unsupported provider",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "unknown",
					Name:     "test",
				},
			},
			outputError: `unsupported method selector provider "unknown"`,
		},
		{
			name: "empty expression",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": ""},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "whitespace only expression",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "    "},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "invalid expression",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "[invalid"},
			},
			outputError: `failed to compile filter selector expression`,
		},
		{
			name: "valid expression returns no error",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
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
			selector := &FilterSelector{}
			err := selector.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}
