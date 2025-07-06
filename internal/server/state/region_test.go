// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package state

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestRegion_Validate(t *testing.T) {

	testCases := []struct {
		name                string
		inputRegion         *Region
		outputErrorContains string
	}{
		{
			name:                "nil region",
			inputRegion:         nil,
			outputErrorContains: "nil region",
		},
		{
			name: "duplicate address",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "default",
				API: []*RegionAPI{
					{Address: "http://127.0.0.1:4646"},
					{Address: "http://127.0.0.1:4646"},
				},
			},
			outputErrorContains: "duplicate API address found",
		},
		{
			name: "no addresses",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "default",
				API:   []*RegionAPI{},
			},
			outputErrorContains: "API list must have at least one entry",
		},
		{
			name: "multiple default addresses",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "default",
				API: []*RegionAPI{
					{Address: "http://127.0.0.1:4646", Default: true},
					{Address: "http://127.0.0.2:4646", Default: true},
				},
			},
			outputErrorContains: "API list can only have one default",
		},
		{
			name: "empty group",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "",
				API: []*RegionAPI{
					{Address: "http://127.0.0.1:4646", Default: true},
				},
			},
			outputErrorContains: "group cannot be empty",
		},
		{
			name: "invalid address format",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "default",
				API: []*RegionAPI{
					{Address: "127.0.0.1:4646", Default: true},
				},
			},
			outputErrorContains: "invalid address",
		},
		{
			name: "invalid TLS",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "default",
				API: []*RegionAPI{
					{Address: "http://127.0.0.1:4646", Default: true},
				},
				TLS: &RegionTLS{
					CACert:     "cacert",
					ClientCert: "clientcert",
					ClientKey:  "clientkey",
					ServerName: "servername",
					Insecure:   true,
				},
			},
			outputErrorContains: "failed to find any PEM data in certificate input",
		},
		{
			name: "full valid",
			inputRegion: &Region{
				Name:  "euw1",
				Group: "eu",
				API: []*RegionAPI{
					{Address: "http://10.0.0.10:4646", Default: true},
					{Address: "http://10.0.0.11:4646", Default: false},
					{Address: "http://10.0.0.12:4646", Default: false},
				},
				TLS: &RegionTLS{
					CACert:     "-----BEGIN CERTIFICATE-----\nMIIDCTCCAq+gAwIBAgIQKbjUtElJSSdrCrvFQq1uzDAKBggqhkjOPQQDAjCBxzEL\nMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2Nv\nMRowGAYDVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAV\nBgNVBAoTDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1\nTm9tYWQgQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcw\nMTE0MDQwHhcNMjUwNzA1MTEwMzMwWhcNMzAwNzA0MTEwMzMwWjCBxzELMAkGA1UE\nBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRowGAYD\nVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAVBgNVBAoT\nDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1Tm9tYWQg\nQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcwMTE0MDQw\nWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQ3nCAE8gHFQei7qd/JINv8bsNBKrqB\n700SWtXOnugCpB0BOXcBl65sUT1WasbEij6T0Hf6ETPiiNXyssZGTBIho3sweTAO\nBgNVHQ8BAf8EBAMCAYYwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQgdhO0EdKo\n1E9YtyErPLBldP6YqBj7TxNe3tmCjyF9dl0wKwYDVR0jBCQwIoAgdhO0EdKo1E9Y\ntyErPLBldP6YqBj7TxNe3tmCjyF9dl0wCgYIKoZIzj0EAwIDSAAwRQIgBS8zAC+t\nPl6Gprpk7xm5jdvrLnm/aBh88DBxM5crcD8CIQDl510EtK2oFIzSdzydyG5GGA/g\ny0WMYMTB04KOek6zUQ==\n-----END CERTIFICATE-----",
					ClientCert: "-----BEGIN CERTIFICATE-----\nMIICrjCCAlSgAwIBAgIQLQxOsqmULY2y2EfUgKzzBjAKBggqhkjOPQQDAjCBxzEL\nMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkNBMRYwFAYDVQQHEw1TYW4gRnJhbmNpc2Nv\nMRowGAYDVQQJExExMDEgU2Vjb25kIFN0cmVldDEOMAwGA1UEERMFOTQxMDUxFzAV\nBgNVBAoTDkhhc2hpQ29ycCBJbmMuMQ4wDAYDVQQLEwVOb21hZDE+MDwGA1UEAxM1\nTm9tYWQgQWdlbnQgQ0EgNTU0NTgwNDQ2MDM3MzgxODg1NTYwOTU4NTQ1MjQ0MTcw\nMTE0MDQwHhcNMjUwNzA1MTEwNDI5WhcNMjYwNzA1MTEwNDI5WjAeMRwwGgYDVQQD\nExNjbGllbnQuZ2xvYmFsLm5vbWFkMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE\ncO08iQM220O7eBm5VPH9MrUAzaHKY5u5jn3jl3LdCD8oHPUCCwdgRqC+JSaB2ALL\nMwPZuaDWXyQwt3tgkyYmgqOByTCBxjAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwKQYDVR0OBCIEIEj3\n1eswrouGaCvFSptIT8RjfMatMReVHkC3WcTgerCmMCsGA1UdIwQkMCKAIHYTtBHS\nqNRPWLchKzywZXT+mKgY+08TXt7Zgo8hfXZdMC8GA1UdEQQoMCaCE2NsaWVudC5n\nbG9iYWwubm9tYWSCCWxvY2FsaG9zdIcEfwAAATAKBggqhkjOPQQDAgNIADBFAiEA\nvct1V1qVoZLcrvObu8gFvjEpHxDpvhNf53SI3GPvj9kCIHvsjyVlnj82pgjL2yh2\nlbzvj3aoyUCKkFBsjYTcVXND\n-----END CERTIFICATE-----",
					ClientKey:  "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEINRy4CwInWabRFidIR1+vQYXm1KP+GomvOAibzxfVv2HoAoGCCqGSM49\nAwEHoUQDQgAEcO08iQM220O7eBm5VPH9MrUAzaHKY5u5jn3jl3LdCD8oHPUCCwdg\nRqC+JSaB2ALLMwPZuaDWXyQwt3tgkyYmgg==\n-----END EC PRIVATE KEY-----",
					ServerName: "servername",
					Insecure:   true,
				},
			},
			outputErrorContains: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualOutput := tc.inputRegion.Validate()
			if tc.outputErrorContains != "" {
				must.ErrorContains(t, actualOutput, tc.outputErrorContains)
			} else {
				must.NoError(t, actualOutput)
			}
		})
	}
}
