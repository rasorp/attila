// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package pointer

func Of[A any](a A) *A { return &a }
