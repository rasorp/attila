// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type LimitPickerFactory struct{}

func (LimitPickerFactory) New() jobsdk.RegionPicker { return &LimitPicker{} }

type limitPickerConfig struct {
	base   *jobsdk.RegionPickerBaseConfig
	params *limitPickerProviderConfig
}

type limitPickerProviderConfig struct {
	Num int `json:"num"`
}

type LimitPicker struct {
	config *limitPickerConfig
}

func (s *LimitPicker) SetConfig(cfg *jobsdk.RegionPickerConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	decodedCfg := limitPickerConfig{base: cfg.RegionPickerBaseConfig}

	if err := decodeParams(cfg.ProviderConfig, &decodedCfg.params); err != nil {
		return err
	}

	if decodedCfg.params == nil {
		return errors.New("limit config required")
	}
	if decodedCfg.params.Num < 0 {
		return errors.New("limit config option \"num\" cannot be negative")
	}

	s.config = &decodedCfg
	return nil
}

func (s *LimitPicker) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *LimitPicker) Run(req *jobsdk.RegionPickerRunRequest) (*jobsdk.RegionPickerRunResult, error) {
	if s.config.params.Num == 0 {
		return &jobsdk.RegionPickerRunResult{RegionCandidates: []jobsdk.RegisterRuleRegionCandidate{}}, nil
	}
	if s.config.params.Num >= len(req.RegionCandidates) {
		return &jobsdk.RegionPickerRunResult{RegionCandidates: CopyCandidates(req.RegionCandidates)}, nil
	}
	return &jobsdk.RegionPickerRunResult{RegionCandidates: CopyCandidates(req.RegionCandidates[:s.config.params.Num])}, nil
}
