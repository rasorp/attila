// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package state

type Error interface {
	Error() string
	Err() error
	StatusCode() int
	String() string
}
