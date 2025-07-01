// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package index

import (
	"errors"
	"fmt"
)

// SingleIndexer implements both memdb.Indexer and memdb.SingleIndexer. It may
// be used in a memdb.IndexSchema to specify functions that generate the index
// value for memdb.Txn operations.
type SingleIndexer struct {

	// readIndex is used by memdb for Txn.Get, Txn.First, and other operations
	// that read data.
	ReadIndex

	// writeIndex is used by memdb for Txn.Insert, Txn.Delete, and other
	// operations that write data to the index.
	WriteIndex
}

// ReadIndex implements memdb.Indexer. It exists so that a function can be used
// to provide the interface. Unlike memdb.Indexer, a readIndex function accepts
// only a single argument.
type ReadIndex func(arg any) ([]byte, error)

func (f ReadIndex) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("index supports only a single arg")
	}
	return f(args[0])
}

var ErrMissingValueForIndex = errors.New("object is missing a value for this index")

// WriteIndex implements memdb.SingleIndexer. It exists so that a function
// can be used to provide this interface.
//
// Instead of a bool return value, writeIndex expects errMissingValueForIndex to
// indicate that an index could not be build for the object. It will translate
// this error into a false value to satisfy the memdb.SingleIndexer interface.
type WriteIndex func(raw any) ([]byte, error)

func (f WriteIndex) FromObject(raw any) (bool, []byte, error) {
	v, err := f(raw)
	if errors.Is(err, ErrMissingValueForIndex) {
		return false, nil, nil
	}
	return err == nil, v, err
}
