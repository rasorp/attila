// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

const nomadEUW1Config = `
name = "euw1"
`

const nomadEUW2Config = `
name = "euw2"
ports {
  http = 5646
  rpc  = 5647
  serf = 5648
}
`

const attilaRegionEUW1Config = `
name  = "euw1"
group = "europe"

api {
  address = "http://localhost:4646"
  default = true
}
`

const attilaRegionEUW2Config = `
name  = "euw2"
group = "europe"

api {
  address = "http://localhost:5646"
  default = true
}
`

const attilaJobRegRuleConfig = `
name = "europe-platform"

region_context { kind = "namespace" }

region_picker "europe-region" {
  provider = "expr"

  config {
    expression = "regions.filter(r, r.group == \"europe\")"
  }
}

region_picker "platform-namespace" {
  provider = "filter"

  config {
    expression = "any(region_namespace, {.Name == \"platform\"})"
  }
}
`

const attilaJobRegMethodConfig = `
name     = "europe-platform"
selector = "Namespace == \"platform\""

rule {
  name = "europe-platform"
}
`

const nomadJobConfig = `
job "example" {
  namespace = "platform"
  group "cache" {
    task "redis" {
      driver = "docker"
      config {
        image = "redis:7"
      }
    }
  }
}
`

// writeTempFile writes content to a temp file in t.TempDir() and returns the
// path.
func writeTempFile(t *testing.T, dir string, name string, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	must.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// nomadServer manages a Nomad dev agent lifecycle.
type nomadServer struct {
	name   string
	config string
	port   int
	cmd    *exec.Cmd
}

func (s *nomadServer) start(t *testing.T, ctx context.Context) {
	t.Helper()

	workspace := t.TempDir()
	cfgContent, err := os.ReadFile(s.config)
	must.NoError(t, err)

	cfgPath := filepath.Join(workspace, fmt.Sprintf("nomad_%s.hcl", s.name))
	err = os.WriteFile(cfgPath, cfgContent, 0o644)
	must.NoError(t, err)

	s.cmd = exec.CommandContext(ctx, "nomad", "agent", "-dev", "-config="+cfgPath)
	s.cmd.Dir = workspace
	stdout, _ := s.cmd.StdoutPipe()
	stderr, _ := s.cmd.StderrPipe()

	err = s.cmd.Start()
	must.NoError(t, err)

	go func() { _, _ = io.Copy(io.Discard, stdout) }()
	go func() { _, _ = io.Copy(io.Discard, stderr) }()

	must.NoError(t, waitForNomadHealth(ctx, s.port, 30*time.Second))
}

func (s *nomadServer) stop(t *testing.T) {
	t.Helper()
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_ = s.cmd.Wait()
	}
}

// attilaServer manages the Attila server lifecycle.
type attilaServer struct {
	bin string
	cmd *exec.Cmd
}

func (s *attilaServer) build(ctx context.Context, projRoot string) error {

	buildCmd := exec.CommandContext(ctx, "go", "build", "-tags=e2e", "-o", s.bin, "./internal/cmd")
	buildCmd.Dir = projRoot
	out, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %w\n%s", err, out)
	}
	return nil
}

func (s *attilaServer) start(ctx context.Context) error {
	s.cmd = exec.CommandContext(ctx, s.bin, "server", "run", "--state-memory-enabled")
	stdout, _ := s.cmd.StdoutPipe()
	stderr, _ := s.cmd.StderrPipe()

	err := s.cmd.Start()
	if err != nil {
		return err
	}

	go func() { _, _ = io.Copy(io.Discard, stdout) }()
	go func() { _, _ = io.Copy(io.Discard, stderr) }()
	return nil
}

func (s *attilaServer) stop(t *testing.T) {
	t.Helper()
	if s.cmd != nil && s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
		_ = s.cmd.Wait()
	}
}

// runCLI executes the Attila CLI and returns stdout+stderr as a string.
func runCLI(ctx context.Context, bin string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, bin, args...)
	out, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(out))
	return result, err
}

// parseRegionList extracts region names from `attila region list` output.
func parseRegionList(output string) []string {
	var names []string
	lines := strings.SplitSeq(output, "\n")
	for line := range lines {
		parts := strings.Fields(line)
		if len(parts) > 0 && parts[0] != "Name" { // skip header row
			names = append(names, parts[0])
		}
	}
	return names
}

// parsePlanOutput extracts key fields from a plan output.
type planFields struct {
	ID           string
	NumRegions   int
	JobID        string
	Namespace    string
	PartialError string
}

func parsePlanOutput(output string) planFields {
	var pf planFields

	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "ID"):
			pf.ID = extractValue(line)
		case strings.HasPrefix(line, "Num Regions"):
			pf.NumRegions = extractInt(extractValue(line))
		case strings.HasPrefix(line, "Job ID"):
			pf.JobID = extractValue(line)
		case strings.HasPrefix(line, "Job Namespace"):
			pf.Namespace = extractValue(line)
		case strings.HasPrefix(line, "Partial Error"):
			pf.PartialError = extractValue(line)
		}
	}
	return pf
}

func extractValue(kv string) string {
	if _, after, ok := strings.Cut(kv, "="); ok {
		return strings.TrimSpace(after)
	}
	return ""
}

func extractInt(s string) int {
	var n int
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}

// waitForHealth waits for an HTTP health check to succeed.
func waitForNomadHealth(ctx context.Context, port int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("health check timed out for port %d after %v", port, time.Since(start))
		}

		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/v1/agent/self", port))
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			_ = resp.Body.Close()
		}

		time.Sleep(200 * time.Millisecond)

		// Allow some time for startup before aggressively checking.
		if time.Since(start) < 2*time.Second {
			continue
		}
	}
}

func TestE2E(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	testRoot, err := os.Getwd()
	must.NoError(t, err)
	projRoot := filepath.Join(testRoot, "..")

	// Create a shared temp directory for all config files.
	workspace := t.TempDir()

	// Write all demo configs to temp files.
	nomadEUW1Path := writeTempFile(t, workspace, "nomad_euw1.hcl", nomadEUW1Config)
	nomadEUW2Path := writeTempFile(t, workspace, "nomad_euw2.hcl", nomadEUW2Config)
	attilaRegionEUW1Path := writeTempFile(t, workspace, "attila_region_euw1.hcl", attilaRegionEUW1Config)
	attilaRegionEUW2Path := writeTempFile(t, workspace, "attila_region_euw2.hcl", attilaRegionEUW2Config)
	attilaJobRegRulePath := writeTempFile(t, workspace, "attila_job_reg_rule.hcl", attilaJobRegRuleConfig)
	attilaJobRegMethodPath := writeTempFile(t, workspace, "attila_job_reg_method.hcl", attilaJobRegMethodConfig)
	nomadJobPath := writeTempFile(t, workspace, "nomad_job.nomad.hcl", nomadJobConfig)

	ns1 := &nomadServer{
		name:   "euw1",
		config: nomadEUW1Path,
		port:   4646,
	}
	ns2 := &nomadServer{
		name:   "euw2",
		config: nomadEUW2Path,
		port:   5646,
	}

	ns1.start(t, ctx)
	defer ns1.stop(t)
	ns2.start(t, ctx)
	defer ns2.stop(t)

	atBin := t.TempDir() + "/attila"
	srv := &attilaServer{bin: atBin}

	must.NoError(t, srv.build(ctx, projRoot))

	err = srv.start(ctx)
	must.NoError(t, err)
	defer srv.stop(t)

	// Create euw1 and euw2 regions within Attila.
	euw1RegionCreate, err := runCLI(ctx, atBin, "region", "create", attilaRegionEUW1Path)
	must.NoError(t, err)
	must.StrContains(t, euw1RegionCreate, "Name")
	must.StrContains(t, euw1RegionCreate, "euw1")

	euw2RegionCreate, err := runCLI(ctx, atBin, "region", "create", attilaRegionEUW2Path)
	must.NoError(t, err)
	must.StrContains(t, euw2RegionCreate, "Name")
	must.StrContains(t, euw2RegionCreate, "euw2")

	// Perform a list and get on the region objects.
	regionList, err := runCLI(ctx, atBin, "region", "list")
	must.NoError(t, err)
	names := parseRegionList(regionList)
	must.SliceContains(t, names, "euw1")
	must.SliceContains(t, names, "euw2")

	regionGet, err := runCLI(ctx, atBin, "region", "get", "euw1")
	must.NoError(t, err)
	must.StrContains(t, regionGet, "Name")
	must.StrContains(t, regionGet, "euw1")
	must.StrContains(t, regionGet, "europe")
	must.StrContains(t, regionGet, "localhost:4646")

	ruleCreate, err := runCLI(ctx, atBin, "job", "register", "rule", "create", attilaJobRegRulePath)
	must.NoError(t, err)
	must.StrContains(t, ruleCreate, "europe-platform")

	methodCreate, err := runCLI(ctx, atBin, "job", "register", "method", "create", attilaJobRegMethodPath)
	must.NoError(t, err)
	must.StrContains(t, methodCreate, "europe-platform")

	// Generate a plan which should not match any regions because the target
	// namespace does not exist in either Nomad region.
	planCreate1, err := runCLI(ctx, atBin, "job", "register", "plan", "create", nomadJobPath)
	must.NoError(t, err)
	must.StrContains(t, planCreate1, "Job ID")
	must.StrContains(t, planCreate1, "example")

	// Create the required Nomad namespace in a single region.
	nsCmd := exec.CommandContext(ctx, "nomad", "namespace", "apply", "-address=http://127.0.0.1:4646", "platform")
	namespaceCreate, err := nsCmd.CombinedOutput()
	must.NoError(t, err)
	must.StrContains(t, strings.TrimSpace(string(namespaceCreate)), "Successfully applied namespace")

	planCreate2, err := runCLI(ctx, atBin, "job", "register", "plan", "create", nomadJobPath)
	must.NoError(t, err)
	must.StrContains(t, planCreate2, "Num Regions")
	must.StrContains(t, planCreate2, "euw1")

	fields := parsePlanOutput(planCreate2)
	must.NotEq(t, "", fields.ID)

	planRun, err := runCLI(ctx, atBin, "job", "register", "plan", "run", fields.ID, nomadJobPath)
	must.NoError(t, err)
	must.StrContains(t, planRun, "Num Regions")
	must.StrContains(t, planRun, "euw1")
	must.StrContains(t, planRun, "Eval ID")
}
