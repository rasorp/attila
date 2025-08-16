// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package mock

import (
	"github.com/hashicorp/nomad/api"
	"github.com/oklog/ulid/v2"

	"github.com/rasorp/attila/internal/helper/pointer"
	"github.com/rasorp/attila/internal/server/state"
)

func Region() *state.Region {
	return &state.Region{
		Name:  "mock-" + ulid.Make().String(),
		Group: "europe",
		Auth: &state.RegionAuth{
			Token: "very-secret-token",
		},
		API: []*state.RegionAPI{
			{
				Address: "http://10.0.0.10:4646",
				Default: true,
			},
			{
				Address: "http://10.0.0.11:4646",
				Default: false,
			},
			{
				Address: "http://10.0.0.12:4646",
				Default: false,
			},
		},
		TLS: &state.RegionTLS{
			CACert:     "-----BEGIN CERTIFICATE-----\nMIIDCTCCAq+gAwIBAgIQKbjUtElJSSdrCrvFQq1uzDAKBggqhkjOPQQDAjCBxzEL\nMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2Nv\nMRowGAYDVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAV\nBgNVBAoTDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1\nTm9tYWQgQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcw\nMTE0MDQwHhcNMjUwNzA1MTEwMzMwWhcNMzAwNzA0MTEwMzMwWjCBxzELMAkGA1UE\nBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRowGAYD\nVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAVBgNVBAoT\nDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1Tm9tYWQg\nQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcwMTE0MDQw\nWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQ3nCAE8gHFQei7qd/JINv8bsNBKrqB\n700SWtXOnugCpB0BOXcBl65sUT1WasbEij6T0Hf6ETPiiNXyssZGTBIho3sweTAO\nBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQgdhO0EdKo\n1E9YtyErPLBldP6YqBj7TxNe3tmCjyF9dl0wKwYDVR0jBCQwIoAgdhO0EdKo1E9Y\ntyErPLBldP6YqBj7TxNe3tmCjyF9dl0wCgYIKoZIzj0EAwIDSAAwRQIgBS8zAC+t\nPl6Gprpk7xm5jdvrLnm/aBh88DBxM5crcD8CIQDl510EtK2oFIzSdzydyG5GGA/g\ny0WMYMTB04KOek6zUQ==\n-----END CERTIFICATE-----",
			ClientCert: "-----BEGIN CERTIFICATE-----\nMIICrjCCAlSgAwIBAgIQLQxOsqmULY2y2EfUgKzzBjAKBggqhkjOPQQDAjCBxzEL\nMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2Nv\nMRowGAYDVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAV\nBgNVBAoTDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1\nTm9tYWQgQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcw\nMTE0MDQwHhcNMjUwNzA1MTEwNDI5WhcNMjYwNzA1MTEwNDI5WjAeMRwwGgYDVQQD\nExNjbGllbnQuZ2xvYmFsLm5vbWFkMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\ncO08iQM220O7eBm5VPH9MrUAzaHKY5u5jn3jl3LdCD8oHPUCCwdgRqC+JSaB2ALL\nMwPZuaDWXyQwt3tgkyYmgqOByTCBxjAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwKQYDVR0OBCIEIEj3\n1eswrouGaCvFSptIT8RjfMatMReVHkC3WcTgerCmMCsGA1UdIwQkMCKAIHYTtBHS\nqNRPWLchKzywZXT+mKgY+08TXt7Zgo8hfXZdMC8GA1UdEQQoMCaCE2NsaWVudC5n\nbG9iYWwubm9tYWSCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiEA\nvct1V1qVoZLcrvObu8gFvjEpHxDpvhNf53SI3GPvj9kCIHvsjyVlnj82pgjL2yh2\nlbzvj3aoyUCKkFBsjYTcVXND\n-----END CERTIFICATE-----",
			ClientKey:  "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEINRy4CwInWabRFidIR1+vQYXm1KP+GomvOAibzxfVv2HoAoGCCqGSM49\nAwEHoUQDQgAEcO08iQM220O7eBm5VPH9MrUAzaHKY5u5jn3jl3LdCD8oHPUCCwdg\nRqC+JSaB2ALLMwPZuaDWXyQwt3tgkyYmgg==\n-----END EC PRIVATE KEY-----",
			ServerName: "servername",
			Insecure:   true,
		},
	}
}

func JobRegistrationMethod() *state.JobRegisterMethod {
	return &state.JobRegisterMethod{
		Name:     "mock-" + ulid.Make().String(),
		Selector: "Namespace == \"platform\"",
		Rules: []*state.JobRegisterMethodRuleLink{
			{Name: "mock-" + ulid.Make().String()},
			{Name: "mock-" + ulid.Make().String()},
			{Name: "mock-" + ulid.Make().String()},
		},
	}
}

func JobRegistrationPlan() *state.JobRegisterPlan {
	return &state.JobRegisterPlan{
		ID: ulid.Make(),
		Job: &api.Job{
			ID:        pointer.Of("mock-" + ulid.Make().String()),
			Namespace: pointer.Of("platform"),
		},
	}
}

func JobRegistrationRule() *state.JobRegisterRule {
	return &state.JobRegisterRule{
		Name: "mock-" + ulid.Make().String(),
		RegionContexts: []state.JobRegisterRuleRegionContext{
			state.JobRegisterRuleContextNamespace,
			state.JobRegisterRuleContextNodePool,
		},
		RegionFilter: &state.JobRegisterRuleFilter{
			Expression: &state.JobRegisterRuleFilterExpression{
				Selector: "any(region_namespace, {.Name == \"platform\"})",
			},
		},
		RegionPicker: &state.JobRegisterRulePicker{
			Expression: &state.JobRegisterRuleFilterExpression{
				Selector: "filter(regions, .Group == \"europe\")",
			},
		},
	}
}
