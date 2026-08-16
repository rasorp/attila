// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// ParseConfig loads a JSON or HCL config file into obj.
//
// HCL decoding needs an extra normalization pass after hclsimple.DecodeFile.
// Most strongly-typed fields decode directly into their final Go values, but
// dynamic fields such as map[string]any can retain HCL syntax nodes
// (for example *hcl.Attribute) rather than the evaluated scalar/list/map value.
//
// The normalization step walks the decoded object graph, evaluates any HCL
// attributes that were captured inside dynamic maps, and rewrites them into
// plain Go values so downstream JSON marshaling and validation behave as
// expected.
func ParseConfig(path string, obj any) error {

	if _, err := os.Stat(path); err != nil {
		return err
	}

	fileExt := filepath.Ext(path)

	switch fileExt {
	case ".json":
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if err := json.Unmarshal(fileBytes, obj); err != nil {
			return fmt.Errorf("failed to unmarshal file: %w", err)
		}
	case ".hcl":
		evalCtx := hclEvalCtx(filepath.Dir(path))
		if err := hclsimple.DecodeFile(path, evalCtx, obj); err != nil {
			return fmt.Errorf("failed to decode file: %w", err)
		}
		if err := normalizeDecodedHCL(obj, evalCtx); err != nil {
			return fmt.Errorf("failed to normalize decoded file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %q", fileExt)
	}

	return nil
}

// normalizeDecodedHCL starts a post-decode walk of obj and rewrites any HCL
// syntax objects that were preserved inside dynamically-typed fields into plain
// Go values.
func normalizeDecodedHCL(obj any, evalCtx *hcl.EvalContext) error {
	return normalizeDecodedHCLValue(reflect.ValueOf(obj), evalCtx)
}

// normalizeDecodedHCLValue recursively traverses structs, slices, and maps.
//
// The interesting case is maps: when HCL decodes into map[string]any, values in
// a nested block may remain as *hcl.Attribute rather than string/bool/number
// values. Those entries are normalized in-place so callers see the same shapes
// they would expect from JSON decoding.
func normalizeDecodedHCLValue(v reflect.Value, evalCtx *hcl.EvalContext) error {
	if !v.IsValid() {
		return nil
	}

	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return normalizeDecodedHCLValue(v.Elem(), evalCtx)
	}

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := normalizeDecodedHCLValue(v.Field(i), evalCtx); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if err := normalizeDecodedHCLValue(v.Index(i), evalCtx); err != nil {
				return err
			}
		}
	case reflect.Map:
		// Map entries are where dynamic HCL values can survive decoding as syntax
		// nodes, so normalize each element before writing it back.
		for _, key := range v.MapKeys() {
			normalizedValue, err := normalizeDecodedHCLMapValue(v.MapIndex(key).Interface(), evalCtx)
			if err != nil {
				return err
			}

			reflectedValue, err := normalizedMapValue(v.Type().Elem(), normalizedValue)
			if err != nil {
				return err
			}
			v.SetMapIndex(key, reflectedValue)
		}
	}

	return nil
}

// normalizeDecodedHCLMapValue converts values stored in a dynamic map into
// JSON-compatible Go values.
//
// HCL blocks decoded into map[string]any may contain *hcl.Attribute values,
// nested maps, or slices. This helper evaluates attributes and recursively
// normalizes nested collections.
func normalizeDecodedHCLMapValue(value any, evalCtx *hcl.EvalContext) (any, error) {
	switch typedValue := value.(type) {
	case *hcl.Attribute:
		return evaluateHCLAttribute(typedValue, evalCtx)
	case map[string]any:
		normalizedMap := make(map[string]any, len(typedValue))
		for key, nestedValue := range typedValue {
			normalizedValue, err := normalizeDecodedHCLMapValue(nestedValue, evalCtx)
			if err != nil {
				return nil, err
			}
			normalizedMap[key] = normalizedValue
		}
		return normalizedMap, nil
	case []any:
		normalizedSlice := make([]any, len(typedValue))
		for i, nestedValue := range typedValue {
			normalizedValue, err := normalizeDecodedHCLMapValue(nestedValue, evalCtx)
			if err != nil {
				return nil, err
			}
			normalizedSlice[i] = normalizedValue
		}
		return normalizedSlice, nil
	default:
		return value, nil
	}
}

// evaluateHCLAttribute evaluates a decoded HCL attribute using the same eval
// context used during file parsing and converts the resulting cty.Value into a
// plain Go value via JSON round-tripping.
//
// The JSON conversion is intentional: downstream consumers of ParseConfig work
// with ordinary Go values that later marshal cleanly into JSON request bodies.
func evaluateHCLAttribute(attribute *hcl.Attribute, evalCtx *hcl.EvalContext) (any, error) {
	attributeValue, diags := attribute.Expr.Value(evalCtx)
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to evaluate HCL attribute %q: %s", attribute.Name, diags.Error())
	}

	simpleJSONValue := ctyjson.SimpleJSONValue{}
	simpleJSONValue.Value = attributeValue

	attributeBytes, err := json.Marshal(simpleJSONValue)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HCL attribute %q: %w", attribute.Name, err)
	}

	var normalizedValue any
	if err := json.Unmarshal(attributeBytes, &normalizedValue); err != nil {
		return nil, fmt.Errorf("failed to unmarshal HCL attribute %q: %w", attribute.Name, err)
	}
	return normalizedValue, nil
}

// normalizedMapValue adapts a normalized Go value so it can be written back
// into a reflected map whose element type may be broader than the concrete
// normalized value.
func normalizedMapValue(mapValueType reflect.Type, normalizedValue any) (reflect.Value, error) {
	if normalizedValue == nil {
		return reflect.Zero(mapValueType), nil
	}

	reflectedValue := reflect.ValueOf(normalizedValue)
	if reflectedValue.Type().AssignableTo(mapValueType) {
		return reflectedValue, nil
	}
	if reflectedValue.Type().ConvertibleTo(mapValueType) {
		return reflectedValue.Convert(mapValueType), nil
	}
	return reflect.Value{}, fmt.Errorf("cannot assign normalized value type %T into map value type %s", normalizedValue, mapValueType)
}
