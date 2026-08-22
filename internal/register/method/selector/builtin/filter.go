// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	"github.com/rasorp/attila/internal/helper/convert"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

type FilterSelectorFactory struct{}

func (FilterSelectorFactory) New() jobsdk.MethodSelector { return &FilterSelector{} }

type filterSelectorConfig struct {
	base   *jobsdk.MethodSelectorBaseConfig
	params *filterSelectorProviderConfig
}

type filterSelectorProviderConfig struct {
	Expression string `json:"expression"`
}

type FilterSelector struct {
	config  *filterSelectorConfig
	program *vm.Program
}

func (s *FilterSelector) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *FilterSelector) Provider() string { return jobsdk.MethodSelectorProviderFilter }

func (s *FilterSelector) Run(req *jobsdk.MethodSelectorRunRequest) (*jobsdk.MethodSelectorRunResult, error) {

	if req == nil || req.Job == nil {
		return &jobsdk.MethodSelectorRunResult{Match: false}, nil
	}

	result, err := expr.Run(s.program, map[string]any{"job": convert.JobToMap(req.Job)})
	if err != nil {
		return nil, fmt.Errorf("failed to run filter selector expression: %w", err)
	}

	matched, ok := result.(bool)
	if !ok {
		return nil, fmt.Errorf("filter selector returned incorrect type: %T", result)
	}

	return &jobsdk.MethodSelectorRunResult{Match: matched}, nil
}

func (s *FilterSelector) SetConfig(cfg *jobsdk.MethodSelectorConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	var decodedCfg filterSelectorConfig
	decodedCfg.base = cfg.MethodSelectorBaseConfig

	if err := decodeParams(cfg.ProviderConfig, &decodedCfg.params); err != nil {
		return err
	}

	if decodedCfg.params == nil {
		return errors.New("filter config required")
	}
	if strings.TrimSpace(decodedCfg.params.Expression) == "" {
		return errors.New("param \"expression\" cannot be empty")
	}

	program, err := expr.Compile(decodedCfg.params.Expression, expr.AsBool())
	if err != nil {
		return fmt.Errorf("failed to compile filter selector expression: %w", err)
	}

	s.config = &decodedCfg
	s.program = program
	return nil
}
