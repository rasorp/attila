// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestNamespace_Validate(t *testing.T) {
	testCases := []struct {
		name           string
		inputNamespace *Namespace
		expectError    bool
	}{
		{
			name:           "nil namespace",
			inputNamespace: nil,
			expectError:    true,
		},
		{
			name: "valid namespace with dashes",
			inputNamespace: &Namespace{
				Name:        "valid-name-123",
				Description: "A valid description",
				Regions:     []string{"euw1"},
			},
			expectError: false,
		},
		{
			name: "valid namespace with underscores",
			inputNamespace: &Namespace{
				Name:    "namespace_456",
				Regions: []string{"euw1"},
			},
			expectError: false,
		},
		{
			name: "valid name at max length",
			inputNamespace: &Namespace{
				Name:    "a1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345",
				Regions: []string{"euw1"},
			},
			expectError: false,
		},
		{
			name: "name too long",
			inputNamespace: &Namespace{
				Name:    "a12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456712345678901",
				Regions: []string{"euw1"},
			},
			expectError: true,
		},
		{
			name: "empty name",
			inputNamespace: &Namespace{
				Name:    "",
				Regions: []string{"euw1"},
			},
			expectError: true,
		},
		{
			name: "name with invalid characters",
			inputNamespace: &Namespace{
				Name:    "invalid name!",
				Regions: []string{"euw1"},
			},
			expectError: true,
		},
		{
			name: "description max length",
			inputNamespace: &Namespace{
				Name:        "platform",
				Description: "a12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
				Regions:     []string{"euw1"},
			},
			expectError: false,
		},
		{
			name: "description too long",
			inputNamespace: &Namespace{
				Name:        "platform",
				Description: "a123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901a123456789012345678901234567901a12345678912345678901234",
				Regions:     []string{"euw1"},
			},
			expectError: true,
		},
		{
			name: "valid with multiple regions",
			inputNamespace: &Namespace{
				Name:    "platform",
				Regions: []string{"euw1", "euw2"},
			},
			expectError: false,
		},
		{
			name: "valid with region wildcard",
			inputNamespace: &Namespace{
				Name:    "scoped",
				Regions: []string{"*"},
			},
			expectError: false,
		},
		{
			name: "invalid with region wildcard",
			inputNamespace: &Namespace{
				Name:    "scoped",
				Regions: []string{"euw1", "*"},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualError := tc.inputNamespace.Validate()

			if tc.expectError {
				must.Error(t, actualError)
			} else {
				must.NoError(t, actualError)
			}
		})
	}
}

func TestNamespace_Stub(t *testing.T) {
	testCases := []struct {
		name           string
		inputNamespace *Namespace
		expectedOutput *NamespaceStub
	}{
		{
			name:           "nil",
			inputNamespace: nil,
			expectedOutput: nil,
		},
		{
			name: "full",
			inputNamespace: &Namespace{
				Name:        "platform",
				Description: "platforms namespace",
				Regions:     []string{"euw1", "euw2"},
			},
			expectedOutput: &NamespaceStub{
				Name:        "platform",
				Description: "platforms namespace",
			},
		},
		{
			name: "partial",
			inputNamespace: &Namespace{
				Name:        "platform",
				Description: "",
				Regions:     []string{"euw1", "euw2"},
			},
			expectedOutput: &NamespaceStub{
				Name:        "platform",
				Description: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			must.Eq(t, tc.expectedOutput, tc.inputNamespace.Stub())
		})
	}
}
