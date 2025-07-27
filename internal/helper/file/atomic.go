// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"errors"
	"os"
	"path/filepath"
)

// AtomicWrite performs an atomic file write by writing to a temporary file
// before moving it to the desired location.
func AtomicWrite(filename string, data []byte, perm os.FileMode) error {

	existingFile, err := os.Stat(filename)
	if err == nil && !existingFile.Mode().IsRegular() {
		return errors.New("file is not a regular file")
	}

	// Generate a temporary file using the desired name as a prefix. It writes
	// this to the same target directory as the end file.
	f, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename)+".tmp")
	if err != nil {
		return err
	}

	tmpFileName := f.Name()

	var writeErr error

	// Run a cleanup function that will perform a best effort removal of the
	// temporary file if we encounter an error.
	defer func(writeErr error) {
		if writeErr != nil {
			_ = f.Close()
			_ = os.Remove(tmpFileName)
		}
	}(writeErr)

	if _, writeErr = f.Write(data); writeErr != nil {
		return writeErr
	}
	if writeErr = f.Chmod(perm); writeErr != nil {
		return writeErr
	}
	if writeErr = f.Sync(); writeErr != nil {
		return writeErr
	}
	if writeErr = f.Close(); writeErr != nil {
		return writeErr
	}
	return os.Rename(tmpFileName, filename)
}
