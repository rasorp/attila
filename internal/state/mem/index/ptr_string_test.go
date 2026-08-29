// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package index

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestStringPtrFieldIndex_FromObject(t *testing.T) {

	type nestedObj struct {
		Inner *string
	}

	type ptrNestedObj struct {
		Deep *string
	}

	type testObj struct {
		Name      *string
		Email     string
		Meta      map[string]*string
		Value     int
		Nested    nestedObj
		PTRNested *ptrNestedObj
		EmptyPtr  *string
	}

	testCases := []struct {
		name        string
		index       StringPtrFieldIndex
		object      any
		expectOk    bool
		expectBytes []byte
		expectErr   string
	}{
		{
			name:        "string value",
			index:       StringPtrFieldIndex{Field: "Name"},
			object:      testObj{Name: new("hello")},
			expectOk:    true,
			expectBytes: []byte("hello"),
		},
		{
			name:        "empty string value",
			index:       StringPtrFieldIndex{Field: "Name"},
			object:      testObj{Name: new("")},
			expectOk:    true,
			expectBytes: []byte(""),
		},
		{
			name:        "non-pointer",
			index:       StringPtrFieldIndex{Field: "Email"},
			object:      testObj{Email: "user@example.com"},
			expectOk:    true,
			expectBytes: []byte("user@example.com"),
		},
		{
			name:     "nil pointer",
			index:    StringPtrFieldIndex{Field: "Name"},
			object:   testObj{Name: nil},
			expectOk: false,
		},
		{
			name:     "uninitialized pointer",
			index:    StringPtrFieldIndex{Field: "EmptyPtr"},
			object:   testObj{},
			expectOk: false,
		},
		{
			name:     "missing field",
			index:    StringPtrFieldIndex{Field: "NonExistent"},
			object:   testObj{Name: new("hello")},
			expectOk: false,
		},
		{
			name:     "nil object",
			index:    StringPtrFieldIndex{Field: "Name"},
			object:   nil,
			expectOk: false,
		},
		{
			name:      "map unexpected type",
			index:     StringPtrFieldIndex{Field: "Meta"},
			object:    testObj{Name: new("hello"), Meta: map[string]*string{"key": new("hello")}},
			expectErr: "field \"Meta\": unexpected type map[string]*string for StringPtrFieldIndex",
		},
		{
			name:      "int unexpected type",
			index:     StringPtrFieldIndex{Field: "Value"},
			object:    testObj{Name: new("hello"), Value: 42},
			expectErr: "field \"Value\": unexpected type int for StringPtrFieldIndex",
		},
		{
			name:        "nested pointer",
			index:       StringPtrFieldIndex{Field: "Nested.Inner"},
			object:      testObj{Nested: nestedObj{Inner: new("hello")}},
			expectOk:    true,
			expectBytes: []byte("hello"),
		},
		{
			name:     "nested nil pointer",
			index:    StringPtrFieldIndex{Field: "Nested.Inner"},
			object:   testObj{Nested: nestedObj{Inner: nil}},
			expectOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, bytes, err := tc.index.FromObject(tc.object)

			if tc.expectErr != "" {
				must.ErrorContains(t, err, tc.expectErr)
			} else {
				must.Eq(t, tc.expectOk, ok)
				must.Eq(t, tc.expectBytes, bytes)
			}
		})
	}
}

func TestStringPtrFieldIndex_FromArgs(t *testing.T) {
	testCases := []struct {
		name        string
		index       StringPtrFieldIndex
		args        []any
		expectBytes []byte
		expectErr   string
	}{
		{
			name:        "[]byte argument",
			index:       StringPtrFieldIndex{Field: "Name"},
			args:        []any{[]byte("test")},
			expectBytes: []byte("test"),
		},
		{
			name:        "string argument",
			index:       StringPtrFieldIndex{Field: "Name"},
			args:        []any{"test"},
			expectBytes: []byte("test"),
		},
		{
			name:        "empty string argument",
			index:       StringPtrFieldIndex{Field: "Name"},
			args:        []any{""},
			expectBytes: []byte(""),
		},
		{
			name:        "empty []byte argument",
			index:       StringPtrFieldIndex{Field: "Name"},
			args:        []any{[]byte{}},
			expectBytes: []byte{},
		},
		{
			name:      "no arguments",
			index:     StringPtrFieldIndex{Field: "Name"},
			args:      []any{},
			expectErr: "StringPtrFieldIndex expects exactly 1 arg",
		},
		{
			name:      "too many arguments",
			index:     StringPtrFieldIndex{Field: "Name"},
			args:      []any{"a", "b"},
			expectErr: "StringPtrFieldIndex expects exactly 1 arg",
		},
		{
			name:      "integer argument",
			index:     StringPtrFieldIndex{Field: "Name"},
			args:      []any{42},
			expectErr: `field "Name": unexpected arg type int`,
		},
		{
			name:      "nil argument",
			index:     StringPtrFieldIndex{Field: "Name"},
			args:      []any{nil},
			expectErr: `field "Name": unexpected arg type <nil>`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bytes, err := tc.index.FromArgs(tc.args...)

			if tc.expectErr != "" {
				must.ErrorContains(t, err, tc.expectErr)
			} else {
				must.NoError(t, err)
				must.Eq(t, tc.expectBytes, bytes)
			}
		})
	}
}
