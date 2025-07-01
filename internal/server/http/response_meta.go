// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package http

type internalResponseMeta interface {
	StatusCode() int
}

type internalResponseMetaImpl struct {
	code int
}

func newInternalResponseMeta(c int) internalResponseMetaImpl {
	return internalResponseMetaImpl{
		code: c,
	}
}

func (r internalResponseMetaImpl) StatusCode() int {
	return r.code
}
