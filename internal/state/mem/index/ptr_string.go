// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package index

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-memdb"
)

// StringPtrFieldIndex is a memdb.Indexer that reads pointer-to-string fields
// from structs. It exists because go-memdb's built-in StringFieldIndex cannot
// read *string fields as it uses reflect.Value.String() which only works on
// concrete string types, not pointers to strings.
type StringPtrFieldIndex struct {

	// Field is the dot-separated field path to extract from the object. e.g.
	// "Job.Namespace" or just "Name".
	Field string
}

var _ memdb.Indexer = (*StringPtrFieldIndex)(nil)

// FromObject extracts the index value from the given object. It dereferences
// pointer fields that StringFieldIndex cannot handle. Returns false when the
// field is missing or nil, which is treated as AllowMissing (the row won't be
// indexed for this key).
func (i *StringPtrFieldIndex) FromObject(obj any) (bool, []byte, error) {
	val := ptrFieldValue(i.Field, obj)

	// If val is nil, the field is missing — treat as "allow missing".
	if val == nil {
		return false, nil, nil
	}

	switch v := val.(type) {
	case *string:
		return true, []byte(*v), nil
	case string:
		return true, []byte(v), nil
	default:
		return false, nil, fmt.Errorf("field %q: unexpected type %T for StringPtrFieldIndex", i.Field, val)
	}
}

// FromArgs implements the optional SingleIndexer interface so that this indexer
// can be used as a read query target. Supports raw []byte or string values.
func (i *StringPtrFieldIndex) FromArgs(args ...any) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("StringPtrFieldIndex expects exactly 1 arg")
	}

	switch v := args[0].(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("field %q: unexpected arg type %T", i.Field, args[0])
	}
}

// ptrFieldValue reads a dot-separated field path from obj using reflection. It
// auto-dereferences pointers at every level of the path. Returns nil if any
// part of the path is missing or nil.
func ptrFieldValue(path string, obj any) any {
	if obj == nil {
		return nil
	}

	objVal := reflect.ValueOf(obj)
	if objVal.Kind() == reflect.Pointer {
		if objVal.IsNil() {
			return nil
		}
		objVal = objVal.Elem()
	}
	if objVal.Kind() != reflect.Struct {
		return nil
	}

	var currentVal = objVal

	for part := range strings.SplitSeq(path, ".") {
		if currentVal.Kind() == reflect.Pointer {
			if currentVal.IsNil() {
				return nil
			}
			currentVal = currentVal.Elem()
		}
		if currentVal.Kind() != reflect.Struct {
			return nil
		}

		f := currentVal.FieldByName(part)
		if !f.IsValid() {
			return nil
		}
		currentVal = f
	}

	if currentVal.IsZero() || (currentVal.Kind() == reflect.Pointer && currentVal.IsNil()) {
		return nil
	}

	switch currentVal.Kind() {
	case reflect.String, reflect.Pointer:
		return currentVal.Interface()
	default:
		// Not a string or pointer — StringFieldIndex can't handle this type.
		return currentVal.Interface()
	}
}
