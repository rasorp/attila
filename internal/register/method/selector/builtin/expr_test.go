// Copyright James Rasell 2026, 2025
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/shoenig/test/must"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestExprSelector_Provider(t *testing.T) {
	selector := &ExprSelector{}
	must.Eq(t, "expr", selector.Provider())
}

func TestExprSelector_Name(t *testing.T) {

	t.Run("before config", func(t *testing.T) {
		selector := &ExprSelector{}
		must.Eq(t, "", selector.Name())
	})

	t.Run("after config", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "my-expr",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		selector := &ExprSelector{}
		must.NoError(t, selector.SetConfig(cfg))
		must.Eq(t, "my-expr", selector.Name())
	})
}

func TestExprSelector_Run(t *testing.T) {

	nomadJob := api.Job{
		ID:        new("test-job-1"),
		Namespace: new("platform"),
		Type:      new("batch"),
	}

	t.Run("nil request returns non-match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		s := &ExprSelector{}
		must.NoError(t, s.SetConfig(cfg))
		result, err := s.Run(nil)
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("nil job returns non-match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		s := &ExprSelector{}
		must.NoError(t, s.SetConfig(cfg))
		result, err := s.Run(&jobsdk.MethodSelectorRunRequest{Job: nil})
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("expression matches", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
		}
		result, err := runExprSelector(cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("expression does not match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Type == \"batch\""},
		}
		result, err := runExprSelector(cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("complex expression matches", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\" && job.Type == \"batch\""},
		}
		result, err := runExprSelector(cfg, &nomadJob)
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("complex expression does not match", func(t *testing.T) {
		cfg := &jobsdk.MethodSelectorConfig{
			MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
				Provider: "expr",
				Name:     "test",
			},
			ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\" && job.Type == \"service\""},
		}
		result, err := runExprSelector(cfg, &nomadJob)
		must.NoError(t, err)
		must.False(t, result.Match)
	})
}

func runExprSelector(cfg *jobsdk.MethodSelectorConfig, job *api.Job) (*jobsdk.MethodSelectorRunResult, error) {
	s := &ExprSelector{}
	if err := s.SetConfig(cfg); err != nil {
		panic(err)
	}
	return s.Run(&jobsdk.MethodSelectorRunRequest{Job: job})
}

func TestExprSelector_SetConfig(t *testing.T) {
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
					Provider: "expr",
					Name:     "",
				},
			},
			outputError: "\"name\" cannot be empty",
		},
		{
			name: "empty expression",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": ""},
			},
			outputError: `param "expression" cannot be empty`,
		},
		{
			name: "invalid expression",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "[invalid"},
			},
			outputError: `failed to compile expr expression`,
		},
		{
			name: "valid expression returns no error",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "test",
				},
				ProviderConfig: map[string]any{"expression": "job.Namespace == \"platform\""},
			},
		},
		{
			name: "no params",
			cfg: &jobsdk.MethodSelectorConfig{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "test",
				},
				ProviderConfig: nil,
			},
			outputError: `expr config required`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := &ExprSelector{}
			err := s.SetConfig(tc.cfg)

			if tc.outputError != "" {
				must.ErrorContains(t, err, tc.outputError)
			} else {
				must.NoError(t, err)
			}
		})
	}
}
