// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"sort"
	"testing"

	"github.com/shoenig/test/must"
)

func TestTopologyOverviewSort(t *testing.T) {

	overview := []*TopologyOverview{
		{RegionName: "euw2"},
		{RegionName: "euw1"},
		{RegionName: "eue1"},
		{RegionName: "eue1a"},
		{RegionName: "eue2"},
	}

	sort.Sort(TopologyOverviewSort(overview))

	must.Eq(
		t,
		[]*TopologyOverview{
			{RegionName: "eue1"},
			{RegionName: "eue1a"},
			{RegionName: "eue2"},
			{RegionName: "euw1"},
			{RegionName: "euw2"},
		},
		overview,
	)
}

func TestTopologyDetail_SortNodes(t *testing.T) {

	t.Run("without nodes", func(t *testing.T) {
		topology := &TopologyDetail{Nodes: []*NodeTopology{}}
		topology.SortNodes()
		must.Eq(t, &TopologyDetail{Nodes: []*NodeTopology{}}, topology)
	})

	t.Run("with nodes", func(t *testing.T) {
		topology := &TopologyDetail{
			Nodes: []*NodeTopology{
				{NodePool: "b", Name: "b"},
				{NodePool: "b", Name: "a"},
				{NodePool: "z", Name: "a"},
				{NodePool: "a", Name: "c"},
				{NodePool: "d", Name: "a"},
				{NodePool: "a", Name: "b"},
			},
		}
		topology.SortNodes()
		must.Eq(
			t,
			&TopologyDetail{
				Nodes: []*NodeTopology{
					{NodePool: "a", Name: "b"},
					{NodePool: "a", Name: "c"},
					{NodePool: "b", Name: "a"},
					{NodePool: "b", Name: "b"},
					{NodePool: "d", Name: "a"},
					{NodePool: "z", Name: "a"},
				},
			},
			topology,
		)
	})
}

func TestNodeTopology_SortAllocs(t *testing.T) {

	t.Run("without allocs", func(t *testing.T) {
		topology := &NodeTopology{AllocationTopology: []*AllocationTopology{}}
		topology.SortAllocs()
		must.Eq(t, &NodeTopology{AllocationTopology: []*AllocationTopology{}}, topology)
	})

	t.Run("with allocs", func(t *testing.T) {
		topology := &NodeTopology{
			AllocationTopology: []*AllocationTopology{
				{Namespace: "b", JobID: "z", ID: "a"},
				{Namespace: "a", JobID: "b", ID: "b"},
				{Namespace: "z", JobID: "a", ID: "a"},
				{Namespace: "a", JobID: "a", ID: "a"},
				{Namespace: "a", JobID: "b", ID: "a"},
				{Namespace: "b", JobID: "y", ID: "a"},
			},
		}
		topology.SortAllocs()
		must.Eq(
			t,
			&NodeTopology{
				AllocationTopology: []*AllocationTopology{
					{Namespace: "a", JobID: "a", ID: "a"},
					{Namespace: "a", JobID: "b", ID: "a"},
					{Namespace: "a", JobID: "b", ID: "b"},
					{Namespace: "b", JobID: "y", ID: "a"},
					{Namespace: "b", JobID: "z", ID: "a"},
					{Namespace: "z", JobID: "a", ID: "a"},
				},
			},
			topology,
		)
	})
}
