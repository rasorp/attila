// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package index

import (
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/domain"
)

// ReadULIDIndex can be used as a memdb.Indexer query via ReadIndex and
// allows querying by a ULID.
func ReadULIDIndex(arg any) ([]byte, error) {
	id, ok := arg.(ulid.ULID)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T for ULIDQuery index", arg)
	}

	return id.Bytes(), nil
}

func WriteULIDIndex(raw any) ([]byte, error) {
	plan, ok := raw.(*domain.JobRegisterPlan)
	if !ok {
		return nil, fmt.Errorf("unexpected type %T for job register plan", raw)
	}
	return plan.ID.Bytes(), nil
}
