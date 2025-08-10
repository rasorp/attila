// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package nomad

import (
	"time"

	"github.com/hashicorp/nomad/api"
)

type Topology struct {
	Overview *Overview `json:"overview"`
	Detail   *Detail   `json:"detail"`

	// CreateTime marks the time at which the topology collection was created
	// and can help callers identify how stale the data is.
	CreateTime time.Time `json:"create_time"`
}

type Overview struct {
	RegionName        string `json:"region_name"`
	NumServers        int    `json:"num_servers"`
	NumClients        int    `json:"num_clients"`
	NumAllocs         int    `json:"num_allocs"`
	CPUAllocatable    int64  `json:"cpu_allocatable"`
	CPUAllocated      int64  `json:"cpu_allocated"`
	MemoryAllocatable int64  `json:"memory_allocatable"`
	MemoryAllocated   int64  `json:"memory_allocated"`
}

type Detail struct {
	Servers []*Server `json:"servers,omitempty"`
	Nodes   []*Node   `json:"nodes,omitempty"`
}

type Server struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Version     string `json:"version"`
	RaftVersion string `json:"raft_version"`
}

type Node struct {
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	NodePool           string        `json:"node_pool"`
	Status             string        `json:"status"`
	CPUAllocatable     int64         `json:"cpu_allocatable"`
	CPUAllocated       int64         `json:"cpu_allocated"`
	MemoryAllocatable  int64         `json:"memory_allocatable"`
	MemoryAllocated    int64         `json:"memory_allocated"`
	AllocationTopology []*Allocation `json:"allocations"`
}

type Allocation struct {
	ID        string `json:"id"`
	JobID     string `json:"job_id"`
	Namespace string `json:"namespace"`
	CPU       int64  `json:"cpu"`
	Memory    int64  `json:"memory"`
}

func NewTopology(name string) *Topology {
	return &Topology{
		Overview: &Overview{
			RegionName: name,
		},
		Detail:     &Detail{},
		CreateTime: time.Now(),
	}
}

// AddNode takes the node object and allocations from the node and adds them to
// the topology object. The caller must ensure the allocation array comes from
// the passed node, as the function does not perform any checks.
func (r *Topology) AddNode(node *api.NodeListStub, allocs []*api.Allocation) {

	allocatableCPU := node.NodeResources.Cpu.CpuShares - int64(node.ReservedResources.Cpu.CpuShares)
	allocatableMem := node.NodeResources.Memory.MemoryMB - int64(node.ReservedResources.Memory.MemoryMB)

	r.Overview.NumClients++
	r.Overview.CPUAllocatable += allocatableCPU
	r.Overview.MemoryAllocatable += allocatableMem

	nt := Node{
		ID:                 node.ID,
		Name:               node.Name,
		NodePool:           node.NodePool,
		Status:             node.Status,
		CPUAllocatable:     allocatableCPU,
		MemoryAllocatable:  allocatableMem,
		AllocationTopology: make([]*Allocation, 0),
	}

	for _, alloc := range allocs {

		if alloc.ClientTerminalStatus() {
			continue
		}

		r.Overview.NumAllocs++

		r.Overview.CPUAllocated += int64(*alloc.Resources.CPU)
		nt.CPUAllocated += int64(*alloc.Resources.CPU)
		r.Overview.MemoryAllocated += int64(*alloc.Resources.MemoryMB)
		nt.MemoryAllocated += int64(*alloc.Resources.MemoryMB)

		nt.AllocationTopology = append(
			nt.AllocationTopology,
			&Allocation{
				ID:        alloc.ID,
				JobID:     alloc.JobID,
				Namespace: alloc.Namespace,
				CPU:       int64(*alloc.Resources.CPU),
				Memory:    int64(*alloc.Resources.MemoryMB),
			},
		)
	}

	r.Detail.Nodes = append(r.Detail.Nodes, &nt)
}

// AddServer adds the passed server objects to the topology tracking. It is the
// caller's responsibility to ensure the server belongs to the named region the
// topology is tracking.
func (r *Topology) AddServer(server *api.AgentMember) {
	r.Overview.NumServers++

	r.Detail.Servers = append(
		r.Detail.Servers,
		&Server{
			ID:          server.Tags["id"],
			Name:        server.Name,
			Status:      server.Status,
			Version:     server.Tags["build"],
			RaftVersion: server.Tags["raft_vsn"],
		},
	)
}
