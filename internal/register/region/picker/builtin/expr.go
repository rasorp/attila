// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type ExprPickerFactory struct{}

func (ExprPickerFactory) New() jobsdk.RegionPicker { return &ExprPicker{} }

type exprPickerConfig struct {
	base   *jobsdk.RegionPickerBaseConfig
	params *exprPickerProviderConfig
}

type exprPickerProviderConfig struct {
	Expression string `json:"expression"`
}

type ExprPicker struct {
	config  *exprPickerConfig
	program cel.Program
}

func (s *ExprPicker) SetConfig(cfg *jobsdk.RegionPickerConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	decodedCfg := exprPickerConfig{base: cfg.RegionPickerBaseConfig}

	if err := decodeParams(cfg.ProviderConfig, &decodedCfg.params); err != nil {
		return err
	}

	if decodedCfg.params == nil {
		return errors.New("expr config required")
	}
	if strings.TrimSpace(decodedCfg.params.Expression) == "" {
		return errors.New("param \"expression\" cannot be empty")
	}

	env, err := newCelEnv()
	if err != nil {
		return err
	}

	ast, issues := env.Compile(decodedCfg.params.Expression)
	if issues.Err() != nil {
		return fmt.Errorf("failed to compile expr expression: %w", issues.Err())
	}

	program, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("failed to create expr program: %w", err)
	}

	s.config = &decodedCfg
	s.program = program
	return nil
}

func (s *ExprPicker) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *ExprPicker) Run(req *jobsdk.RegionPickerRunRequest) (*jobsdk.RegionPickerRunResult, error) {

	regionCandidates := candidatesToCelMaps(req.RegionCandidates)

	result, _, err := s.program.Eval(map[string]any{
		"regions":    regionCandidates,
		"candidates": regionCandidates,
		"job":        anyToCelMap(req.Job),
		"rule":       anyToCelMap(req.Rule),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to run expr picker expression: %w", err)
	}

	pickedCandidates, err := candidatesFromCELResult(result)
	if err != nil {
		return nil, fmt.Errorf("expr picker returned incorrect type: %w", err)
	}

	return &jobsdk.RegionPickerRunResult{RegionCandidates: CopyCandidates(pickedCandidates)}, nil
}

func newCelEnv() (*cel.Env, error) {
	env, err := cel.NewEnv(
		cel.Variable("regions", cel.ListType(cel.DynType)),
		cel.Variable("candidates", cel.ListType(cel.DynType)),
		cel.Variable("job", cel.DynType),
		cel.Variable("rule", cel.DynType),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create expr env: %w", err)
	}
	return env, nil
}

func anyToCelMap(v any) any {
	if v == nil {
		return map[string]any{}
	}
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return v
}

func candidatesToCelMaps(candidates []jobsdk.RegisterRuleRegionCandidate) []map[string]any {
	if candidates == nil {
		return nil
	}

	result := make([]map[string]any, 0, len(candidates))
	for _, candidate := range candidates {
		result = append(result, map[string]any{
			"name":     candidate.Name,
			"group":    candidate.Group,
			"metadata": candidate.Metadata,
			"context":  candidate.Context,
		})
	}
	return result
}

func candidatesFromCELResult(result ref.Val) ([]jobsdk.RegisterRuleRegionCandidate, error) {
	if result == nil {
		return nil, nil
	}

	if candidates, err := convertCELSlice(result); err == nil {
		return candidates, nil
	}
	if candidate, err := convertCELCandidate(result); err == nil {
		return []jobsdk.RegisterRuleRegionCandidate{candidate}, nil
	}

	return nil, fmt.Errorf("%T", result.Value())
}

func convertCELSlice(result ref.Val) ([]jobsdk.RegisterRuleRegionCandidate, error) {
	native, err := result.ConvertToNative(reflect.TypeFor[[]jobsdk.RegisterRuleRegionCandidate]())
	if err == nil {
		return native.([]jobsdk.RegisterRuleRegionCandidate), nil
	}

	nativeAny, err := result.ConvertToNative(reflect.TypeFor[[]any]())
	if err != nil {
		return nil, err
	}

	raw := nativeAny.([]any)
	candidates := make([]jobsdk.RegisterRuleRegionCandidate, 0, len(raw))
	for i, item := range raw {
		switch v := item.(type) {
		case jobsdk.RegisterRuleRegionCandidate:
			candidates = append(candidates, v)
		case map[string]any:
			candidates = append(candidates, candidateFromMap(v))
		default:
			return nil, fmt.Errorf("element %d has type %T", i, item)
		}
	}
	return candidates, nil
}

func convertCELCandidate(result ref.Val) (jobsdk.RegisterRuleRegionCandidate, error) {
	native, err := result.ConvertToNative(reflect.TypeFor[jobsdk.RegisterRuleRegionCandidate]())
	if err == nil {
		return native.(jobsdk.RegisterRuleRegionCandidate), nil
	}

	nativeMap, err := result.ConvertToNative(reflect.TypeFor[map[string]any]())
	if err != nil {
		return jobsdk.RegisterRuleRegionCandidate{}, err
	}
	return candidateFromMap(nativeMap.(map[string]any)), nil
}

func candidateFromMap(raw map[string]any) jobsdk.RegisterRuleRegionCandidate {
	candidate := jobsdk.RegisterRuleRegionCandidate{}

	if name, ok := raw["name"].(string); ok {
		candidate.Name = name
	}
	if name, ok := raw["Name"].(string); ok && candidate.Name == "" {
		candidate.Name = name
	}
	if group, ok := raw["group"].(string); ok {
		candidate.Group = group
	}
	if group, ok := raw["Group"].(string); ok && candidate.Group == "" {
		candidate.Group = group
	}
	if metadata, ok := raw["metadata"].(map[string]any); ok {
		candidate.Metadata = metadata
	}
	if metadata, ok := raw["Metadata"].(map[string]any); ok && candidate.Metadata == nil {
		candidate.Metadata = metadata
	}
	if context, ok := raw["context"].(map[string]any); ok {
		candidate.Context = context
	}
	if context, ok := raw["Context"].(map[string]any); ok && candidate.Context == nil {
		candidate.Context = context
	}

	return candidate
}
