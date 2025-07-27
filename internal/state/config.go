// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
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
// the object is nil. The returned error could be a multierror and should
// indicate a terminal error in the process which intends to use the config
// object.
func (c *Config) Validate() error {

	if c == nil {
		return errors.New("state config block required")
	}

	var (
		numEnabled int
		mErr       *multierror.Error
	)

	if c.Memory.Enabled() {
		numEnabled++
	}

	if c.File.Enabled() {
		numEnabled++

		if err := c.File.Validate(); err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	// We must have one configured backend.
	switch numEnabled {
	case 0:
		mErr = multierror.Append(mErr, errors.New("no state backend enabled"))
	case 1:
	default:
		mErr = multierror.Append(mErr, errors.New("only one storage backend can be enabled"))
	}

	return mErr.ErrorOrNil()
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
		if z.Memory.Enable != nil {
			result.Memory.Enable = z.Memory.Enable
		}
	}

	if z.File != nil {
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

// Enabled is a helper function that informs the caller if the memory state
// backend is enabled. It should be used instead of directly querying the
// enabled boolean pointer.
func (m *MemoryConfig) Enabled() bool {
	return m != nil && m.Enable != nil && *m.Enable
}

type FileConfig struct {
	Enable *bool  `hcl:"enabled"`
	Path   string `hcl:"path"`
}

// Enabled is a helper function that informs the caller if the file state
// backend is enabled. It should be used instead of directly querying the
// enable boolean pointer.
func (f *FileConfig) Enabled() bool {
	return f != nil && f.Enable != nil && *f.Enable
}

// Validate performs validation of the file configuration block. If it is not
// enabled, the validation functionality will not run. The returned error could
// be a multierror and should indicate a terminal error in the process which
// intends to use the config object.
func (f *FileConfig) Validate() error {

	// If the file backend is not enabled, then do not perform any further
	// validation as we will not be using any of the values.
	if !f.Enabled() {
		return nil
	}

	// Check the directory value is as expected. Without confirming this we
	// cannot reliably perform further checks, so do not use a multierror here.
	if f.Path == "" {
		return errors.New("must set path parameter")
	}
	if !filepath.IsAbs(f.Path) {
		return fmt.Errorf("path %q is not an absolute path", f.Path)
	}

	var mErr *multierror.Error

	// Reaching this point we have an absolute path configured. This must exist
	// on the filesystem and be a directory.
	dir, err := os.Stat(f.Path)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	if dir != nil && !dir.IsDir() {
		mErr = multierror.Append(mErr, fmt.Errorf("path %q is not a dir", f.Path))
	}

	return mErr.ErrorOrNil()
}
