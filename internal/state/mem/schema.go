// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mem

import (
	"fmt"

	"github.com/hashicorp/go-memdb"

	"github.com/rasorp/attila/internal/state/mem/index"
)

const (
	// indexID is the primary table index used by all tables. It does not have
	// to directly translate to a field named ID on the stored object, but
	// should be mapped to the field which will be used for default lookups.
	indexID = "id"
)

const (
	regionTableName            = "region"
	jobRegisterMethodTableName = "job_register_method"
	jobRegisterRuleTableName   = "job_register_rule"
	jobRegisterPlanTableName   = "job_register_plan"
)

func newTableSchema() *memdb.DBSchema {

	// Get the list of table schema setup functions, so we can initialize the
	// table mapping to the correct length before iterating.
	tableSchemaFuncs := tableSchemas()

	// Create the base DB scheme and initialize the tables, so they can be
	// populated.
	db := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema, len(tableSchemaFuncs)),
	}

	// Iterate the table schema setup functions and add these to our DB table
	// tracking.
	for _, schemaFn := range tableSchemaFuncs {
		schema := schemaFn()
		if _, ok := db.Tables[schema.Name]; ok {
			panic(fmt.Sprintf("duplicate table name: %s", schema.Name))
		}
		db.Tables[schema.Name] = schema
	}
	return db
}

func tableSchemas() []func() *memdb.TableSchema {
	return []func() *memdb.TableSchema{
		jobRegisterMethodTableSchema,
		jobRegisterPlanTableSchema,
		jobRegisterRuleTableSchema,
		regionTableSchema,
	}
}

func regionTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: regionTableName,
		Indexes: map[string]*memdb.IndexSchema{
			indexID: {
				Name:         indexID,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func jobRegisterMethodTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobRegisterMethodTableName,
		Indexes: map[string]*memdb.IndexSchema{
			indexID: {
				Name:         indexID,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func jobRegisterRuleTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobRegisterRuleTableName,
		Indexes: map[string]*memdb.IndexSchema{
			indexID: {
				Name:         indexID,
				AllowMissing: false,
				Unique:       true,
				Indexer: &memdb.StringFieldIndex{
					Field: "Name",
				},
			},
		},
	}
}

func jobRegisterPlanTableSchema() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: jobRegisterPlanTableName,
		Indexes: map[string]*memdb.IndexSchema{
			indexID: {
				Name:         indexID,
				AllowMissing: false,
				Unique:       true,
				Indexer: &index.SingleIndexer{
					ReadIndex:  index.ReadIndex(index.ReadULIDIndex),
					WriteIndex: index.WriteIndex(index.WriteULIDIndex),
				},
			},
			"job_id": {
				Name:         "job_id",
				AllowMissing: true,
				Unique:       true,
				Indexer: &memdb.CompoundIndex{
					Indexes: []memdb.Indexer{
						&memdb.StringFieldIndex{
							Field: "JobNamespace",
						},

						&memdb.StringFieldIndex{
							Field: "JobID",
						},
					},
				},
			},
		},
	}
}
