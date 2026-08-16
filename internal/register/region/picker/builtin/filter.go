// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"
	"fmt"
	"maps"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type FilterPickerFactory struct{}

func (FilterPickerFactory) New() jobsdk.RegionPicker { return &FilterPicker{} }

type filterPickerConfig struct {
	base   *jobsdk.RegionPickerBaseConfig
	params *filterPickerProviderConfig
}

type filterPickerProviderConfig struct {
	Expression string `json:"expression"`
}

type FilterPicker struct {
	config  *filterPickerConfig
	program *vm.Program
}

func (s *FilterPicker) SetConfig(cfg *jobsdk.RegionPickerConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	var decodedCfg filterPickerConfig
	decodedCfg.base = cfg.RegionPickerBaseConfig

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
		return fmt.Errorf("failed to compile filter picker expression: %w", err)
	}

	s.config = &decodedCfg
	s.program = program
	return nil
}

func (s *FilterPicker) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *FilterPicker) Run(req *jobsdk.RegionPickerRunRequest) (*jobsdk.RegionPickerRunResult, error) {

	filteredCandidates := make([]jobsdk.RegisterRuleRegionCandidate, 0, len(req.RegionCandidates))

	for _, candidate := range req.RegionCandidates {
		ctx := map[string]any{
			"job":       req.Job,
			"rule":      req.Rule,
			"candidate": candidate,
			"region":    candidate,
		}
		maps.Copy(ctx, candidate.Context)

		result, err := expr.Run(s.program, ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to run filter picker selector: %w", err)
		}

		selected, ok := result.(bool)
		if !ok {
			return nil, fmt.Errorf("filter picker selector returned incorrect type: %T", result)
		}
		if selected {
			filteredCandidates = append(filteredCandidates, candidate)
		}
	}

	return &jobsdk.RegionPickerRunResult{RegionCandidates: filteredCandidates}, nil
}
