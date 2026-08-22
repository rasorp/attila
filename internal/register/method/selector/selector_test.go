// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package selector

import (
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/shoenig/test/must"
	"go.uber.org/zap/zaptest"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func TestNew(t *testing.T) {
	t.Run("empty configs returns nil errors", func(t *testing.T) {
		s, err := New(zaptest.NewLogger(t), nil)
		must.NoError(t, err)
		must.NotNil(t, s)
	})

	t.Run("valid configs are registered", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "type-checker",
				},
				ProviderConfig: map[string]any{"expression": "job.type == \"batch\""},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.NoError(t, err)
		must.NotNil(t, s)
	})

	t.Run("invalid config produces error", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "filter",
					Name:     "",
				},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.Error(t, err)
		must.Nil(t, s)
	})

	t.Run("invalid provider produces error", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "unknown-provider",
					Name:     "test",
				},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.Error(t, err)
		must.Nil(t, s)
	})
}

func TestSelector_Run(t *testing.T) {

	matchAllJob := &api.Job{ID: new("match-all"), Namespace: new("platform")}
	nomatchJob := &api.Job{ID: new("nomatch-job"), Namespace: new("default")}

	t.Run("nil receiver returns error", func(t *testing.T) {
		var s *Selector
		result, err := s.Run(&RunRequest{Job: matchAllJob})
		must.Error(t, err)
		must.Nil(t, result)
	})

	t.Run("nil request returns error", func(t *testing.T) {
		s, err := New(zaptest.NewLogger(t), nil)
		must.NoError(t, err)
		must.NotNil(t, s)
		result, err := s.Run(nil)
		must.Error(t, err)
		must.Nil(t, result)
	})

	t.Run("empty request job returns error", func(t *testing.T) {
		s, err := New(zaptest.NewLogger(t), nil)
		must.NoError(t, err)
		must.NotNil(t, s)
		result, err := s.Run(&RunRequest{Job: nil})
		must.Error(t, err)
		must.Nil(t, result)
	})

	t.Run("no selectors returns match", func(t *testing.T) {
		s, err := New(zaptest.NewLogger(t), nil)
		must.NoError(t, err)
		must.NotNil(t, s)
		result, err := s.Run(&RunRequest{Job: matchAllJob})
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("single matching selector", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "ns-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.Namespace == "platform"`},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.NoError(t, err)
		result, err := s.Run(&RunRequest{Job: matchAllJob})
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("single non-matching selector returns false", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "ns-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.Namespace == "platform"`},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.NoError(t, err)
		result, err := s.Run(&RunRequest{Job: nomatchJob})
		must.NoError(t, err)
		must.False(t, result.Match)
	})

	t.Run("all selectors match returns true", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "ns-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.Namespace == "platform"`},
			},
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "id-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.ID == "match-all"`},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.NoError(t, err)
		result, err := s.Run(&RunRequest{Job: matchAllJob})
		must.NoError(t, err)
		must.True(t, result.Match)
	})

	t.Run("one selector fails returns false", func(t *testing.T) {
		cfgs := []*jobsdk.MethodSelectorConfig{
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "ns-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.Namespace == "platform"`},
			},
			{
				MethodSelectorBaseConfig: &jobsdk.MethodSelectorBaseConfig{
					Provider: "expr",
					Name:     "id-filter",
				},
				ProviderConfig: map[string]any{"expression": `job.ID == "nomatch-id"`},
			},
		}
		s, err := New(zaptest.NewLogger(t), cfgs)
		must.NoError(t, err)
		result, err := s.Run(&RunRequest{Job: matchAllJob})
		must.NoError(t, err)
		must.False(t, result.Match)
	})
}
