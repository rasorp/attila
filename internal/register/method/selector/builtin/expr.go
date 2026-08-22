// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package builtin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"

	"github.com/rasorp/attila/internal/helper/convert"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

type ExprSelectorFactory struct{}

func (ExprSelectorFactory) New() jobsdk.MethodSelector { return &ExprSelector{} }

type exprSelectorConfig struct {
	base   *jobsdk.MethodSelectorBaseConfig
	params *exprSelectorProviderConfig
}

type exprSelectorProviderConfig struct {
	Expression string `json:"expression"`
}

type ExprSelector struct {
	config  *exprSelectorConfig
	program cel.Program
}

func (s *ExprSelector) Name() string {
	if s.config == nil || s.config.base == nil {
		return ""
	}
	return s.config.base.Name
}

func (s *ExprSelector) Provider() string { return jobsdk.MethodSelectorProviderExpr }

func (s *ExprSelector) Run(req *jobsdk.MethodSelectorRunRequest) (*jobsdk.MethodSelectorRunResult, error) {

	if req == nil || req.Job == nil {
		return &jobsdk.MethodSelectorRunResult{Match: false}, nil
	}

	result, _, err := s.program.Eval(map[string]any{"job": convert.JobToMap(req.Job)})
	if err != nil {
		return nil, fmt.Errorf("failed to run expr selector expression: %w", err)
	}

	matched, ok := result.Value().(bool)
	if !ok {
		return &jobsdk.MethodSelectorRunResult{Match: false}, nil
	}

	return &jobsdk.MethodSelectorRunResult{Match: matched}, nil
}

func (s *ExprSelector) SetConfig(cfg *jobsdk.MethodSelectorConfig) error {

	if err := cfg.Validate(); err != nil {
		return err
	}

	decodedCfg := exprSelectorConfig{base: cfg.MethodSelectorBaseConfig}

	if err := decodeParams(cfg.ProviderConfig, &decodedCfg.params); err != nil {
		return err
	}

	if decodedCfg.params == nil {
		return errors.New("expr config required")
	}
	if strings.TrimSpace(decodedCfg.params.Expression) == "" {
		return errors.New("param \"expression\" cannot be empty")
	}

	env, err := newCelEnv()
	if err != nil {
		return err
	}

	ast, issues := env.Compile(decodedCfg.params.Expression)
	if issues.Err() != nil {
		return fmt.Errorf("failed to compile expr expression: %w", issues.Err())
	}

	program, err := env.Program(ast)
	if err != nil {
		return fmt.Errorf("failed to create expr program: %w", err)
	}

	s.config = &decodedCfg
	s.program = program
	return nil
}

func newCelEnv() (*cel.Env, error) {
	env, err := cel.NewEnv(cel.Variable("job", cel.DynType))
	if err != nil {
		return nil, fmt.Errorf("failed to create expr env: %w", err)
	}
	return env, nil
}
