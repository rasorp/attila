// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"
)

const (
	RegionPickerProviderExpr   = "expr"
	RegionPickerProviderFilter = "filter"
	RegionPickerProviderLimit  = "limit"
	RegionPickerProviderRandom = "random"
)

// RegionPickerConfig is the persisted configuration envelope for a single
// picker stage.
type RegionPickerConfig struct {
	*RegionPickerBaseConfig

	// ProviderConfig is intentionally stored as an opaque map, so each
	// implementation can specify their own configuration structure.
	ProviderConfig map[string]any `json:"config,omitempty"`
}

// RegionPickerBaseConfig
type RegionPickerBaseConfig struct {

	// Provider is the region picker provider that will be called to execute the
	// calculation.
	Provider string `json:"provider"`

	// Name provides a human friendly name to the picker being used. This
	// allows the same provider to be used multiple times and gives a useful
	// value in logs.
	Name string `json:"name"`
}

// Validate provides high level validation of the strategy config. It cannot
// provide exact validation on the parameters which is handled by each
// implementation.
//
// The function handles a nil config object, so callers do not need to check
// this.
func (r *RegionPickerConfig) Validate() error {

	// Perform our top level validation first which will block further
	// validation.
	if r == nil {
		return errors.New("region picker config is nil")
	}

	// Gather up errors, so the operator can get the most information possible
	// and fix it in a single pass.
	var errs []error

	if r.Name == "" {
		errs = append(errs, errors.New("region picker \"name\" cannot be empty"))
	}
	if !slices.Contains([]string{
		RegionPickerProviderExpr,
		RegionPickerProviderFilter,
		RegionPickerProviderLimit,
		RegionPickerProviderRandom,
	}, r.Provider) {
		errs = append(errs, fmt.Errorf("unsupported region picker provider %q", r.Provider))
	}

	return errors.Join(errs...)
}

// RegisterRuleRegionCandidate is the strategy-engine view of a region
// candidate. It is designed to be serializable and independent of internal
// domain types so future plugins can operate on a stable contract.
type RegisterRuleRegionCandidate struct {
	Name     string         `json:"name"`
	Group    string         `json:"group,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
	Context  map[string]any `json:"context,omitempty"`
}

type RegionPickerRule struct {
	Name           string                `json:"name"`
	RegionContexts []string              `json:"region_contexts,omitempty"`
	RegionPickers  []*RegionPickerConfig `json:"region_pickers,omitempty"`
	Metadata       *RegionPickerMetadata `json:"metadata,omitempty"`
}

type RegionPickerMetadata struct {
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

// RegionPickerRunRequest is the execution input for a single region picker
// run stage.
type RegionPickerRunRequest struct {
	Job              any                           `json:"job,omitempty"`
	Rule             *RegionPickerRule             `json:"rule,omitempty"`
	RegionCandidates []RegisterRuleRegionCandidate `json:"region_candidates"`
}

// RegionPickerRunResult carries the next candidate set produced by a
// strategy stage.
type RegionPickerRunResult struct {
	RegionCandidates []RegisterRuleRegionCandidate `json:"region_candidates"`
}

// RegionPicker is the interface that defines the region picker functionality
// when handling job registrations.
type RegionPicker interface {

	// Name returns the operator provided name of the region picker.
	Name() string

	// Run executes a run of the picker logic.
	Run(req *RegionPickerRunRequest) (*RegionPickerRunResult, error)

	// SetConfig is used to set all the required and optional configuration of
	// the picker. It should valid all input and do its best to ensure that when
	// run is called, no errors occur that could have been indentified in this
	// function call.
	SetConfig(cfg *RegionPickerConfig) error
}

// RegionPickerFactory constructs an unconfigured concrete picker
// implementation. All implementations should include this functionality, so the
// internal engine can build and run them.
type RegionPickerFactory interface {

	// New is responsible for constructing a new instance of RegionPicker. It
	// does not need to perform any configuration but is currently provided as a
	// single entry point that can be expanded in the future to include logging
	// or similar.
	New() RegionPicker
}

func EncodeParams(v any) (json.RawMessage, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to encode: %w", err)
	}
	return json.RawMessage(bytes), nil
}

func MustEncodeParams(v any) json.RawMessage {
	bytes, err := EncodeParams(v)
	if err != nil {
		panic(err)
	}
	return json.RawMessage(bytes)
}
