// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/nomad/api"
)

// JobToMap normalizes the incoming *api.Job into a map[string]any so that
// selectors can access fields without depending on the Nomad API types.
func JobToMap(job *api.Job) map[string]any {

	if job == nil {
		return make(map[string]any)
	}

	var m map[string]any

	if err := mapstructure.Decode(job, &m); err != nil {
		panic(fmt.Sprintf("failed to decode job: %v", err))
	}

	// Recursively dereference pointer values so CEL can compare them directly.
	unptrMap(m)

	return m
}

// unptrMap recursively dereferences pointer values in a map.
func unptrMap(m map[string]any) {
	for k, v := range m {
		if v == nil {
			continue
		}
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Pointer:
			if rv.IsNil() {
				m[k] = nil
			} else {
				m[k] = unptr(rv.Elem().Interface())
			}
		case reflect.Struct, reflect.Map, reflect.Slice:
			m[k] = unptr(v)
		}
	}
}

// unptr recursively dereferences pointer values.
func unptr(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return nil
		}
		return unptr(rv.Elem().Interface())
	case reflect.Struct:
		var m map[string]any
		if err := mapstructure.Decode(v, &m); err == nil {
			unptrMap(m)
			return m
		}
		return v
	case reflect.Map:
		m := make(map[any]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			mk := iter.Key().Interface()
			mv := iter.Value().Interface()
			m[mk] = unptr(mv)
		}
		return m
	case reflect.Slice:
		s := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			s[i] = unptr(elem)
		}
		return s
	}
	return v
}
