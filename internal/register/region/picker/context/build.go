// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package context

import (
	"github.com/rasorp/attila/internal/domain"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

type BuildRegionContextFunc func(*domain.Region) (map[string]any, error)

func BuildCandidates(regions []*domain.Region, buildContext BuildRegionContextFunc) ([]jobsdk.RegisterRuleRegionCandidate, error) {
	if len(regions) == 0 {
		return nil, nil
	}

	candidates := make([]jobsdk.RegisterRuleRegionCandidate, 0, len(regions))
	for _, region := range regions {
		candidate := jobsdk.RegisterRuleRegionCandidate{
			Name:  region.Name,
			Group: region.Group,
		}
		if buildContext != nil {
			ctx, err := buildContext(region)
			if err != nil {
				return nil, err
			}
			candidate.Context = ctx
		}
		candidates = append(candidates, candidate)
	}

	return candidates, nil
}
