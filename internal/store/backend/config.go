// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package backend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Memory *MemoryConfig `hcl:"memory,block"`
	File   *FileConfig   `hcl:"file,block"`
}

// DefaultConfig returns the default configuration for the Attila storage
// backend. This does not enable any backend, so operators must be aware of this
// when running the server. While this adds some cognitive overhead, it is easy
// enough to supply a single flag to enable the memory backend.
func DefaultConfig() *Config {
	return &Config{}
}

// Validate performs validation on the config object and all nested
// configuration blocks. The function can be called safely without checking if
// the object is nil. The returned error could wrap multiple errors and should
// indicate a terminal error in the process which intends to use the config
// object.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("state config block required")
	}

	var (
		numEnabled int
		errs       []error
	)

	if c.Memory.Enabled() {
		numEnabled++
	}

	if c.File.Enabled() {
		numEnabled++
		if err := c.File.Validate(); err != nil {
			errs = append(errs, err)
		}
	}

	switch numEnabled {
	case 0:
		errs = append(errs, errors.New("no state backend enabled"))
	case 1:
	default:
		errs = append(errs, errors.New("only one storage backend can be enabled"))
	}

	return errors.Join(errs...)
}

func (c *Config) Merge(z *Config) *Config {
	if c == nil {
		return z
	}
	if z == nil {
		return c
	}

	result := *c

	if z.Memory != nil {
		if result.Memory == nil {
			result.Memory = &MemoryConfig{}
		}
		if z.Memory.Enable != nil {
			result.Memory.Enable = z.Memory.Enable
		}
	}

	if z.File != nil {
		if result.File == nil {
			result.File = &FileConfig{}
		}
		if z.File.Enable != nil {
			result.File.Enable = z.File.Enable
		}
		if z.File.Path != "" {
			result.File.Path = z.File.Path
		}
	}

	return &result
}

type MemoryConfig struct {
	Enable *bool `hcl:"enabled"`
}

// Enabled is a helper function that informs the caller if the memory store
// backend is enabled.
func (m *MemoryConfig) Enabled() bool {
	return m != nil && m.Enable != nil && *m.Enable
}

type FileConfig struct {
	Enable *bool  `hcl:"enabled"`
	Path   string `hcl:"path"`
}

// Enabled is a helper function that informs the caller if the file store
// backend is enabled.
func (f *FileConfig) Enabled() bool {
	return f != nil && f.Enable != nil && *f.Enable
}

// Validate performs validation of the file configuration block. If it is not
// enabled, the validation functionality will not run. The returned error could
// wrap multiple errors and should indicate a terminal error in the process
// which intends to use the config object.
func (f *FileConfig) Validate() error {
	if !f.Enabled() {
		return nil
	}

	if f.Path == "" {
		return errors.New("must set path parameter")
	}
	if !filepath.IsAbs(f.Path) {
		return fmt.Errorf("path %q is not an absolute path", f.Path)
	}

	var errs []error

	dir, err := os.Stat(f.Path)
	if err != nil {
		errs = append(errs, err)
	}
	if dir != nil && !dir.IsDir() {
		errs = append(errs, fmt.Errorf("path %q is not a dir", f.Path))
	}

	return errors.Join(errs...)
}
