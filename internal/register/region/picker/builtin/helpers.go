// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	jobsdk "github.com/rasorp/attila/pkg/job"
)

func decodeParams(raw map[string]any, dst any) error {
	if len(raw) == 0 {
		return nil
	}
	if err := mapstructure.Decode(raw, &dst); err != nil {
		return fmt.Errorf("failed to decode params: %w", err)
	}
	return nil
}

func CopyCandidates(candidates []jobsdk.RegisterRuleRegionCandidate) []jobsdk.RegisterRuleRegionCandidate {
	out := make([]jobsdk.RegisterRuleRegionCandidate, len(candidates))
	copy(out, candidates)
	return out
}
