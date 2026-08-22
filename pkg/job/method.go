// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"errors"
	"fmt"
	"slices"

	"github.com/hashicorp/nomad/api"
)

const (
	MethodSelectorProviderExpr   = "expr"
	MethodSelectorProviderFilter = "filter"
)

// MethodSelectorConfig is the persisted configuration envelope for a single
// method selector stage.
type MethodSelectorConfig struct {
	*MethodSelectorBaseConfig

	// ProviderConfig is intentionally stored as an opaque map, so each
	// implementation can specify their own configuration structure.
	ProviderConfig map[string]any `json:"config,omitempty"`
}

// MethodSelectorBaseConfig holds the common base for all method selectors.
type MethodSelectorBaseConfig struct {

	// Provider is the selector type that will be used to evaluate jobs.
	Provider string `json:"provider"`

	// Name provides a human friendly name for the selector, useful in logs.
	Name string `json:"name"`
}

// Validate provides high level validation of the selector config.
func (m *MethodSelectorConfig) Validate() error {

	if m == nil || m.MethodSelectorBaseConfig == nil {
		return errors.New("method selector config is nil")
	}

	var errs []error

	if m.Name == "" {
		errs = append(errs, errors.New("method selector \"name\" cannot be empty"))
	}
	if !slices.Contains([]string{
		MethodSelectorProviderExpr,
		MethodSelectorProviderFilter,
	}, m.Provider) {
		errs = append(errs, fmt.Errorf("unsupported method selector provider %q", m.Provider))
	}
	return errors.Join(errs...)
}

// MethodSelectorRunRequest is the execution input for a method selector run.
type MethodSelectorRunRequest struct {

	// Job is the incoming job being evaluated against the job registration
	// method.
	Job *api.Job `json:"job,omitempty"`
}

// Validate ensures the request object is valid.
func (m *MethodSelectorRunRequest) Validate() error {
	if m == nil {
		return errors.New("request is empty")
	}
	if m.Job == nil {
		return errors.New("request job is empty")
	}
	return nil
}

// MethodSelectorRunResult carries the boolean outcome of a selector evaluation.
type MethodSelectorRunResult struct {

	// Match indicates whether the job matches this selector and the method
	// should continue with its execution path.
	Match bool `json:"match"`
}

// MethodSelector is the interface that defines method selector functionality
// when handling job planning and registration.
type MethodSelector interface {

	// Name returns the operator provided name of the method selector.
	Name() string

	// Provider returns the underlying provider name of the method selector.
	Provider() string

	// Run evaluates whether a job should be handled by this method's rules.
	Run(req *MethodSelectorRunRequest) (*MethodSelectorRunResult, error)

	// SetConfig configures the selector. It validates all input and ensures
	// that when Run is called, no errors occur that could have been identified
	// during configuration.
	SetConfig(cfg *MethodSelectorConfig) error
}

// MethodSelectorFactory constructs an unconfigured concrete method selector
// implementation.
type MethodSelectorFactory interface {

	// New creates a new unconfigured MethodSelector instance.
	New() MethodSelector
}
