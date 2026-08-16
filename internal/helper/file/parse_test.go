// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/pkg/api"
)

func TestParseConfig_HCL_NormalizesRegionPickerConfig(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "rule.hcl")

	must.NoError(t, os.WriteFile(testFile, []byte(`
name            = "platform_namespace"
region_contexts = ["namespace"]

region_picker "europe-region" {
  provider = "expr"

  config {
    selector = "filter(regions, .Group == \"europe\")"
  }
}

region_picker "random-order" {
  provider = "random"

  config {
    seed = 42
  }
}
`), 0o600))

	var rule api.JobRegisterRule
	must.NoError(t, ParseConfig(testFile, &rule))
	must.Len(t, 2, rule.RegionPickers)

	selector, ok := rule.RegionPickers[0].Config["selector"].(string)
	must.True(t, ok)
	must.Eq(t, `filter(regions, .Group == "europe")`, selector)

	seed, ok := rule.RegionPickers[1].Config["seed"].(float64)
	must.True(t, ok)
	must.Eq(t, float64(42), seed)
}
