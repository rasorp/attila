// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"cmp"
	"context"
	"net/http"
	"slices"
	"sort"
	"time"
)

type TopologiesGetReq struct {
	RegionName string
}

type TopologiesGetResp struct {
	Topology *Topology `json:"topology"`
}

type TopologiesListReq struct{}

type TopologiesListResp struct {
	TopologyOverviews []*TopologyOverview `json:"topologies"`
}

type TopologyOverview struct {
	RegionName        string `json:"region_name"`
	NumServers        int    `json:"num_servers"`
	NumClients        int    `json:"num_clients"`
	NumAllocs         int    `json:"num_allocs"`
	CPUAllocatable    int64  `json:"cpu_allocatable"`
	CPUAllocated      int64  `json:"cpu_allocated"`
	MemoryAllocatable int64  `json:"memory_allocatable"`
	MemoryAllocated   int64  `json:"memory_allocated"`
}

// TopologyOverviewSort allows the list of region topology overviews to be
// sorted by region name.
type TopologyOverviewSort []*TopologyOverview

func (t TopologyOverviewSort) Len() int { return len(t) }

func (t TopologyOverviewSort) Less(a, b int) bool { return t[a].RegionName < t[b].RegionName }

func (t TopologyOverviewSort) Swap(a, b int) { t[a], t[b] = t[b], t[a] }

type Topology struct {
	Overview *TopologyOverview `json:"overview"`
	Detail   *TopologyDetail   `json:"detail"`

	// CreateTime marks the time at which the topology collection was created
	// and can help callers identify how stale the data is.
	CreateTime time.Time `json:"create_time"`
}

type TopologyDetail struct {
	Servers []*ServerTopology `json:"servers"`
	Nodes   []*NodeTopology   `json:"nodes"`
}

// SortNodes sorts the node array in the topology detail by node pool and then
// name. It is not applied to the API response object; the function is provided
// as a convenience if the caller wants it.
func (t *TopologyDetail) SortNodes() {
	slices.SortFunc(t.Nodes, func(a, b *NodeTopology) int {
		return cmp.Or(
			cmp.Compare(a.NodePool, b.NodePool),
			cmp.Compare(a.Name, b.Name),
		)
	})
}

type ServerTopology struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Version     string `json:"version"`
	RaftVersion string `json:"raft_version"`
}

type NodeTopology struct {
	ID                 string                `json:"id"`
	Name               string                `json:"name"`
	NodePool           string                `json:"node_pool"`
	Status             string                `json:"status"`
	CPUAllocatable     int64                 `json:"cpu_allocatable"`
	CPUAllocated       int64                 `json:"cpu_allocated"`
	MemoryAllocatable  int64                 `json:"memory_allocatable"`
	MemoryAllocated    int64                 `json:"memory_allocated"`
	AllocationTopology []*AllocationTopology `json:"allocations"`
}

// SortAllocs sorts the allocation topology array in the node topology. It sorts
// by namespace, job ID, then ID. It is not applied to the API response object;
// the function is provided as a convenience if the caller wants it.
func (nt *NodeTopology) SortAllocs() {
	slices.SortFunc(nt.AllocationTopology, func(a, b *AllocationTopology) int {
		return cmp.Or(
			cmp.Compare(a.Namespace, b.Namespace),
			cmp.Compare(a.JobID, b.JobID),
			cmp.Compare(a.ID, b.ID),
		)
	})
}

type AllocationTopology struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Namespace string `json:"namespace"`
	CPU       int64  `json:"cpu"`
	Memory    int64  `json:"memory"`
}

type Topologies struct {
	client *Client
}

func (c *Client) Topologies() *Topologies {
	return &Topologies{client: c}
}

func (t *Topologies) Get(ctx context.Context, req *TopologiesGetReq) (*TopologiesGetResp, *Response, error) {

	var resp TopologiesGetResp

	httpReq, err := t.client.NewRequest(http.MethodGet, "/v1alpha1/topologies/"+req.RegionName, nil)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := t.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, httpResp, err
	}

	return &resp, httpResp, nil
}

func (t *Topologies) List(ctx context.Context, _ *TopologiesListReq) (*TopologiesListResp, *Response, error) {

	var resp TopologiesListResp

	httpReq, err := t.client.NewRequest(http.MethodGet, "/v1alpha1/topologies", nil)
	if err != nil {
		return nil, nil, err
	}

	httpResp, err := t.client.Do(ctx, httpReq, &resp)
	if err != nil {
		return nil, httpResp, err
	}

	sort.Sort(TopologyOverviewSort(resp.TopologyOverviews))

	return &resp, httpResp, nil
}
