// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
)

// Namespace represents a logical grouping for job registration.
type Namespace struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Regions     []string `json:"regions"`
}

// DefaultNamespace creates a namespace with the builtin "default" name,
// description, and all region access via the wildcard.
func DefaultNamespace() *Namespace {
	return &Namespace{
		Name:        NamespaceDefaultName,
		Description: NamespaceDefaultDescription,
		Regions:     []string{NamespaceRegionsWildcard},
	}
}

// validNamespaceRegex is the compiled regex that a namespace name must match
// with. If it does not, it is considered invalid.
var validNamespaceNameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{1,128}$")

const (
	NamespaceDefaultName        = "default"
	NamespaceDefaultDescription = "The builtin default namespace"
	NamespaceRegionsWildcard    = "*"
)

// Validate performs validation of the namespace object. It is safe to call
// without checking whether the namespace object is nil, although this would
// indicate a serious error in the functionality of the caller.
func (n *Namespace) Validate() error {
	if n == nil {
		return errors.New("namespace object is nil")
	}

	var err []error

	if !validNamespaceNameRegex.MatchString(n.Name) {
		err = append(
			err,
			fmt.Errorf("namespace \"name\" must match regex %s", validNamespaceNameRegex.String()),
		)
	}
	if len(n.Description) > 256 {
		err = append(err, errors.New("namespace \"description\" must be under 257 characters"))
	}

	// The namespace needs at least one region reference.
	//
	// If there is more than one referenced, that list cannot include the
	// wildcard identifier. Operators should either use the wildcard, or not.
	if len(n.Regions) < 1 {
		err = append(err, errors.New("namespace \"regions\" needs at least one entry"))
	}
	if len(n.Regions) > 1 && slices.Contains(n.Regions, NamespaceRegionsWildcard) {
		err = append(
			err,
			errors.New("namespace \"regions\" should have single entry when using a wildcard"),
		)
	}

	return errors.Join(err...)
}

// NamespaceStub is a simplified version of Namespace used for list responses.
type NamespaceStub struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Stub creates a NamespaceStub from the Namespace.
func (n *Namespace) Stub() *NamespaceStub {
	if n == nil {
		return nil
	}
	return &NamespaceStub{
		Name:        n.Name,
		Description: n.Description,
	}
}
