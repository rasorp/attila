// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package convert

import (
	"testing"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/shoenig/test/must"
)

func TestJobToMap(t *testing.T) {

	t.Run("nil job", func(t *testing.T) {
		result := JobToMap(nil)
		must.MapLen(t, 0, result)
	})

	t.Run("full job", func(t *testing.T) {
		result := JobToMap(&api.Job{
			Name:      new("example"),
			Namespace: new("default"),
			Type:      new("service"),
			TaskGroups: []*api.TaskGroup{
				{
					Name: new("cache"),
					Tasks: []*api.Task{
						{
							Name:   "redis",
							Driver: "docker",
							Config: map[string]any{
								"image": "redis:7",
								"ports": []string{"db"},
							},
						},
					},
				},
			},
		})

		must.MapEq(t, map[string]any{
			"Meta":              map[any]any{},
			"consul_namespace":  nil,
			"StatusDescription": nil,
			"CreateIndex":       nil,
			"Type":              "service",
			"all_at_once":       nil,
			"Constraints":       []any{},
			"Multiregion":       nil,
			"Spreads":           []any{},
			"Migrate":           nil,
			"Payload":           []any{},
			"SubmitTime":        nil,
			"Region":            nil,
			"Name":              "example",
			"Datacenters":       []any{},
			"Update":            nil,
			"Reschedule":        nil,
			"UI":                nil,
			"ParentID":          nil,
			"Stable":            nil,
			"Periodic":          nil,
			"ParameterizedJob":  nil,
			"Dispatched":        false,
			"Namespace":         "default",
			"ID":                nil,
			"vault_namespace":   nil,
			"Status":            nil,
			"ModifyIndex":       nil,
			"JobModifyIndex":    nil,
			"node_pool":         nil,
			"Affinities":        []any{},
			"TaskGroups": []any{
				map[string]any{
					"Migrate":                      nil,
					"Consul":                       nil,
					"Disconnect":                   nil,
					"Meta":                         map[any]any{},
					"Services":                     []any{},
					"Spreads":                      []any{},
					"Volumes":                      map[any]any{},
					"RestartPolicy":                nil,
					"EphemeralDisk":                nil,
					"Networks":                     []any{},
					"shutdown_delay":               nil,
					"PreventRescheduleOnLost":      nil,
					"Name":                         "cache",
					"Count":                        nil,
					"ReschedulePolicy":             nil,
					"stop_after_client_disconnect": nil,
					"max_client_disconnect":        nil,
					"Scaling":                      nil,
					"Constraints":                  []any{},
					"Affinities":                   []any{},
					"max_run_duration":             nil,
					"Tasks": []any{
						map[string]any{
							"Lifecycle":       nil,
							"Meta":            map[any]any{},
							"Artifacts":       []any{},
							"Actions":         []any{},
							"User":            "",
							"Env":             map[any]any{},
							"Vault":           nil,
							"Templates":       []any{},
							"csi_plugin":      nil,
							"Kind":            "",
							"Identity":        nil,
							"Affinities":      []any{},
							"Resources":       nil,
							"logs":            nil,
							"VolumeMounts":    []any{},
							"Leader":          false,
							"Name":            "redis",
							"Constraints":     []any{},
							"shutdown_delay":  time.Duration(0),
							"ScalingPolicies": []any{},
							"kill_timeout":    nil,
							"Driver":          "docker",
							"Config": map[any]any{
								"image": "redis:7",
								"ports": []any{"db"},
							},
							"Services":        []any{},
							"RestartPolicy":   nil,
							"Consul":          nil,
							"DispatchPayload": nil,
							"kill_signal":     "",
							"Identities":      []any{},
							"Schedule":        nil,
							"Secrets":         []any{},
						},
					},
					"Update": nil,
				},
			},
			"Stop":                     nil,
			"VersionTag":               nil,
			"Priority":                 nil,
			"DispatchIdempotencyToken": nil,
			"nomad_token_id":           nil,
			"Version":                  nil,
		}, result)
	})
}
