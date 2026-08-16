// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package picker

import (
	"errors"
	"fmt"

	"github.com/rasorp/attila/internal/register/region/picker/builtin"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

type Picker struct {

	// providers is keyed by the operator provided name and not protected by a
	// mutex as they are expected and intended to be access serially.
	providers map[string]jobsdk.RegionPicker
}

func New(cfgs []*jobsdk.RegionPickerConfig) (*Picker, error) {

	// Collect all errors, so we can provide the most feedback in a single go to
	// the caller.
	var errs []error

	pickers := make(map[string]jobsdk.RegionPicker)

	for _, cfg := range cfgs {

		// Perform the high-level validation, so that we do not try and create a
		// provider that doesn't exist or simply will not work.
		if err := cfg.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}

		var factory jobsdk.RegionPickerFactory

		switch cfg.Provider {
		case jobsdk.RegionPickerProviderExpr:
			factory = builtin.ExprPickerFactory{}
		case jobsdk.RegionPickerProviderFilter:
			factory = builtin.FilterPickerFactory{}
		case jobsdk.RegionPickerProviderLimit:
			factory = builtin.LimitPickerFactory{}
		case jobsdk.RegionPickerProviderRandom:
			factory = builtin.RandomPickerFactory{}

		// There are a number of validation points before this default case, so
		// we should never hit it.
		default:
			errs = append(errs, fmt.Errorf("unsupported region picker provider %q", cfg.Provider))
		}

		// Create the new provider and set the configuration which is our way to
		// validate the provider specific config. If the validation passes, we
		// add this provider to out tracking for use.
		pickerProvider := factory.New()

		if err := pickerProvider.SetConfig(cfg); err != nil {
			errs = append(errs, err)
		} else {
			pickers[cfg.Name] = pickerProvider
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Picker{providers: pickers}, nil
}

func (p *Picker) Process(req *jobsdk.RegionPickerRunRequest) ([]jobsdk.RegisterRuleRegionCandidate, error) {

	if req == nil {
		return nil, errors.New("empy picker request")
	}
	if p == nil || p.providers == nil {
		return nil, errors.New("picker is not configured")
	}

	current := builtin.CopyCandidates(req.RegionCandidates)

	for _, picker := range req.Rule.RegionPickers {

		provider, ok := p.providers[picker.Name]
		if !ok {
			return nil, fmt.Errorf("picker provider %q not configured", picker.Provider)
		}

		nextReq := *req
		nextReq.RegionCandidates = current

		result, err := provider.Run(&nextReq)
		if err != nil {
			return nil, err
		}
		if result == nil {
			return nil, fmt.Errorf("region picker %q returned nil result", picker.Name)
		}
		current = builtin.CopyCandidates(result.RegionCandidates)
	}

	return current, nil
}
