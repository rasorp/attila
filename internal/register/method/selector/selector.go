// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package selector

import (
	"errors"
	"fmt"

	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/register/method/selector/builtin"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

// RunRequest is the selector object passed when triggering a run. It is
// specifically distinct from the jobsdk object to decouple the packages and
// allow for independent changes.
type RunRequest struct {

	// Job is the incoming job being evaluated against the job registration
	// method.
	Job *api.Job
}

// Validate ensures the request object is valid.
func (r *RunRequest) Validate() error {
	if r == nil {
		return errors.New("request is empty")
	}
	if r.Job == nil {
		return errors.New("request job is empty")
	}
	return nil
}

// RunResult is the selector object return when a run has finished its
// processing. It is specifically distinct from the jobsdk object to decouple
// the packages and allow for independent changes.
type RunResult struct {

	// Match indicates whether the job matches this selector and the method
	// should continue with its execution path.
	Match bool
}

// Selector manages a set of method selectors keyed by their operator-provided
// name. It processes incoming jobs through the selector pipeline and returns
// methods whose rules should apply.
type Selector struct {
	logger *zap.Logger

	// providers is keyed by the operator provided name and not protected by a
	// mutex as they are expected and intended to be accessed serially.
	providers map[string]jobsdk.MethodSelector
}

// New creates a new Selector from the provided configuration slices. Each
// config is validated and factory-mapped to its builtin provider type.
func New(log *zap.Logger, cfgs []*jobsdk.MethodSelectorConfig) (*Selector, error) {

	var errs []error

	providers := make(map[string]jobsdk.MethodSelector)

	for _, cfg := range cfgs {

		if err := cfg.Validate(); err != nil {
			errs = append(errs, err)
			continue
		}

		var factory jobsdk.MethodSelectorFactory

		switch cfg.Provider {
		case jobsdk.MethodSelectorProviderExpr:
			factory = builtin.ExprSelectorFactory{}
		case jobsdk.MethodSelectorProviderFilter:
			factory = builtin.FilterSelectorFactory{}
		default:
			errs = append(errs, fmt.Errorf("unsupported method selector provider %q", cfg.Provider))
		}

		selector := factory.New()

		if err := selector.SetConfig(cfg); err != nil {
			errs = append(errs, err)
		} else {
			providers[cfg.Name] = selector
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &Selector{
		logger:    log.Named("method_selector"),
		providers: providers,
	}, nil
}

// Run runs all configured selectors against the given job and returns whether
// the job passed all the matching runs.
func (s *Selector) Run(req *RunRequest) (*RunResult, error) {

	if err := req.Validate(); err != nil {
		return nil, err
	}
	if s == nil || s.providers == nil {
		return nil, errors.New("selector is not configured")
	}

	// Start with matched set to true which makes the rest of the logic easier
	// to manage while staying accurate.
	matched := true

	for _, selector := range s.providers {

		s.logger.Info("running job register method selector",
			zap.String("selector_name", selector.Name()),
			zap.String("selector_provider", selector.Provider()),
		)

		result, err := selector.Run(&jobsdk.MethodSelectorRunRequest{Job: req.Job})
		if err != nil {
			return nil, fmt.Errorf("failed to run selector: %w", err)
		}

		s.logger.Info("successfully ran job register method selector",
			zap.String("selector_name", selector.Name()),
			zap.String("selector_provider", selector.Provider()),
			zap.Bool("selector_result_match", result.Match),
		)

		// If the selector didn't return a match, set the value to indicate this
		// and short-circuit. There is no need to continue processing as no
		// match is considered terminal.
		if !result.Match {
			matched = false
			break
		}
	}

	return &RunResult{Match: matched}, nil
}
