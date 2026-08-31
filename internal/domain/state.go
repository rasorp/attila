// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

// State provides access to all storage-backed sub-resources within the system.
// Implementations should be thread-safe, as they are typically shared across
// multiple goroutines.
type State interface {

	// JobRegister returns a state accessor for job registration resources.
	JobRegister() JobRegisterState

	// Namespace returns a state accessor for namespace resources.
	Namespace() NamespaceState

	// Region returns a state accessor for region resources.
	Region() RegionState

	// Name returns the name of the storage backend provider.
	Name() string
}

// JobRegisterState provides access to job registration sub-resources.
type JobRegisterState interface {

	// Plan returns a state accessor for job register plans.
	Plan() JobRegisterPlanState

	// Method returns a state accessor for job register methods.
	Method() JobRegisterMethodState

	// Rule returns a state accessor for job register rules.
	Rule() JobRegisterRuleState

	// Run returns a state accessor for job register runs.
	Run() JobRegisterRunState
}

// StateError is returned by state methods to carry structured error
// information, including the wrapped underlying error and an HTTP status code
// for API responses.
type StateError interface {

	// Error returns the human-readable error message.
	Error() string

	// Err returns the underlying wrapped error.
	Err() error

	// StatusCode returns the HTTP status code associated with this error.
	StatusCode() int
}

// StateErrorResp is a serializable error response used by state
// implementations. It embeds ErrorBody and implements the StateError interface.
type StateErrorResp struct {
	ErrorBody `json:"error"`
}

type ErrorBody struct {
	Msg  string `json:"message"`
	Code int    `json:"code"`
	err  error
}

// NewStateErrorResp creates a new state error response from an underlying error
// and a HTTP status code.
func NewStateErrorResp(e error, c int) *StateErrorResp {
	return &StateErrorResp{
		ErrorBody: ErrorBody{
			err:  e,
			Code: c,
			Msg:  e.Error(),
		},
	}
}

// Error returns the human-readable error message.
func (e *StateErrorResp) Error() string { return e.Msg }

// Err returns the underlying wrapped error.
func (e *StateErrorResp) Err() error { return e.err }

// StatusCode returns the HTTP status code associated with this error.
func (e *StateErrorResp) StatusCode() int { return e.Code }
