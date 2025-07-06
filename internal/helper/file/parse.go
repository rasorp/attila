// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

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
		if err := hclsimple.DecodeFile(path, hclEvalCtx(filepath.Dir(path)), obj); err != nil {
			return fmt.Errorf("failed to decode file: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %q", fileExt)
	}

	return nil
}
