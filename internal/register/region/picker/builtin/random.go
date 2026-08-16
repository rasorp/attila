// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"
	"math/rand"
	"time"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

type RandomPickerFactory struct{}

func (RandomPickerFactory) New() jobsdk.RegionPicker { return &RandomPicker{} }

type randomPickerConfig struct {
	base   *jobsdk.RegionPickerBaseConfig
	params *randomPickerProviderConfig
}

type randomPickerProviderConfig struct {
	Seed *int64 `json:"seed,omitempty"`
}

type RandomPicker struct {
	config *randomPickerConfig
}

func (s *RandomPicker) SetConfig(cfg *jobsdk.RegionPickerConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	decodedCfg := randomPickerConfig{base: cfg.RegionPickerBaseConfig}

	if err := decodeParams(cfg.ProviderConfig, &decodedCfg.params); err != nil {
		return err
	}

	if decodedCfg.params == nil {
		return errors.New("random config required")
	}
	if decodedCfg.params.Seed != nil && *decodedCfg.params.Seed < 0 {
		return errors.New("random config option \"seed\" cannot be negative")
	}

	s.config = &decodedCfg
	return nil
}

func (s *RandomPicker) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *RandomPicker) Run(req *jobsdk.RegionPickerRunRequest) (*jobsdk.RegionPickerRunResult, error) {

	var rng *rand.Rand

	if s.config.params.Seed != nil {
		rng = rand.New(rand.NewSource(*s.config.params.Seed))
	} else {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	shuffledCandidates := CopyCandidates(req.RegionCandidates)
	rng.Shuffle(len(shuffledCandidates), func(i, j int) {
		shuffledCandidates[i], shuffledCandidates[j] = shuffledCandidates[j], shuffledCandidates[i]
	})

	return &jobsdk.RegionPickerRunResult{RegionCandidates: shuffledCandidates}, nil
}
