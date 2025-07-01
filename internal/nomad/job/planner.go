// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/expr-lang/expr"
	"github.com/hashicorp/go-set/v3"
	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/server/state"
)

type Planner struct {
	logger zerolog.Logger

	clients *client.Clients
	job     *api.Job
	state   state.State

	//
	plan *state.JobRegisterPlan
}

type PlannerReq struct {
	Clients *client.Clients
	Job     *api.Job
	State   state.State
}

func NewPlanner(logger zerolog.Logger, req *PlannerReq) *Planner {
	return &Planner{
		clients: req.Clients,
		job:     req.Job,
		logger: logger.With().
			Str("job_id", *req.Job.ID).
			Str("job_namespace", *req.Job.Namespace).
			Str("component", "job_register_planner").
			Logger(),
		plan:  state.NewJobRegisterPlan(*req.Job.ID, *req.Job.Namespace),
		state: req.State,
	}
}

func (p *Planner) Run() (*state.JobRegisterPlan, error) {

	listResp, err := p.state.JobRegister().Method().List(nil)
	if err != nil {
		return nil, err
	}
	if len(listResp.Methods) == 0 {
		return nil, errors.New("found zero job register methods")
	}

	var ruleLinks []*state.JobRegisterMethodRuleLink

	for _, method := range listResp.Methods {

		exprProgram, err := expr.Compile(method.Selector, expr.AsBool())
		if err != nil {
			return nil, fmt.Errorf("failed to compile method selector: %w", err)
		}

		resultBool, err := expr.Run(exprProgram, p.job)
		if err != nil {
			return nil, fmt.Errorf("failed to run method selector: %w", err)
		}

		if resultBool.(bool) {
			ruleLinks = append(ruleLinks, method.Rules...)
		}
	}

	var rules []*state.JobRegisterRule

	for _, ruleLink := range ruleLinks {
		regRule, err := p.state.JobRegister().Rule().Get(&state.JobRegisterRuleGetReq{Name: ruleLink.Name})
		if err != nil {
			return nil, err
		}
		if regRule == nil {
			return nil, fmt.Errorf("job registration rule not found: %q", ruleLink.Name)
		}
		rules = append(rules, regRule.Rule)
	}

	regionListResp, err := p.state.Region().List(nil)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {

		filteredRegions, err := p.runRegisterPlanRule(rule, regionListResp.Regions)
		if err != nil {
			return nil, err
		}

		if err := p.runRegisterPlanPicker(rule, filteredRegions); err != nil {
			return nil, err
		}
	}

	return p.plan, nil
}

func (p *Planner) runRegisterPlanRule(
	rule *state.JobRegisterRule, regions []*state.Region) ([]*state.Region, error) {

	filteredRegions := set.New[*state.Region](0)

	for _, region := range regions {

		p.logger.Debug().
			Str("rule_name", rule.Name).
			Str("region_name", region.Name).
			Msg("performing execution of rule region filter")

		regionClient, err := p.clients.Get(region.Name)
		if err != nil {
			return nil, err
		}

		context := map[string]any{"job": p.job, "region": region}

		if err := populateRegionContext(rule, regionClient, context); err != nil {
			return nil, err
		}

		exprProgram, err := expr.Compile(rule.RegionFilter.Expression.Selelctor, expr.AsBool())
		if err != nil {
			return nil, fmt.Errorf("failed to compile method selector: %w", err)
		}

		resultBool, err := expr.Run(exprProgram, context)
		if err != nil {
			return nil, fmt.Errorf("failed to run method selector: %w", err)
		}

		//
		if resultBool.(bool) {

			filteredRegions.Insert(region)

			p.logger.Debug().
				Str("rule_name", rule.Name).
				Str("region_name", region.Name).
				Msg("region passed rule region filter")
		}
	}

	return filteredRegions.Slice(), nil
}

// runRegisterPlanPicker executes the job registration plan from the picker
// stage onwards and includes populating the result entries for any picked and
// planned regions.
func (p *Planner) runRegisterPlanPicker(
	rule *state.JobRegisterRule, regions []*state.Region) error {

	p.logger.Debug().
		Str("rule_name", rule.Name).
		Int("num_regions", len(regions)).
		Msg("performing execution of rule region picker")

	regionContext := map[string]any{"regions": regions}

	exprProgram, err := expr.Compile(rule.RegionPicker.Expression.Selelctor, expr.AsKind(reflect.Slice))
	if err != nil {
		return fmt.Errorf("failed to compile picker expression selector: %w", err)
	}

	pickerResult, err := expr.Run(exprProgram, regionContext)
	if err != nil {
		return fmt.Errorf("failed to run picker expression selector: %w", err)
	}

	pickedRegions, ok := pickerResult.([]any)
	if !ok {
		return fmt.Errorf("picker expression selector returned incorrect type: %T", pickerResult)
	}

	return p.generatePlanResult(rule.Name, pickedRegions)
}

// generatePlanResult iterates the selected region slice and perform a Nomad
// job plan for each. The Nomad plan and region name will then be added to the
// Attila plan result.
//
// Any failure in calling the Nomad API will result in a failure of the whole
// function.
func (p *Planner) generatePlanResult(ruleName string, regions []any) error {

	for _, pickpickedRegions := range regions {

		r := pickpickedRegions.(*state.Region)

		client, err := p.clients.Get(r.Name)
		if err != nil {
			return fmt.Errorf("failed to get Nomad client, %w", err)
		}

		// TODO(jrasell): add support for job plan diff.
		planResp, _, err := client.Jobs().PlanOpts(p.job, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to call Nomad job plan, %w", err)
		}

		p.plan.AddRegion(r, planResp)

		p.logger.Info().
			Str("rule_name", ruleName).
			Str("region_name", r.Name).
			Msg("region picked by rule picker")
	}

	return nil
}

func populateRegionContext(
	rule *state.JobRegisterRule, client *api.Client, ctx map[string]any) error {

	for _, regionContext := range rule.RegionContexts {
		switch regionContext {
		case state.JobRegisterRuleContextNamespace:
			namespaceList, _, err := client.Namespaces().List(nil)
			if err != nil {
				return err
			}
			ctx["region_namespace"] = namespaceList

		case state.JobRegisterRuleContextNodePool:
			nodepoolList, _, err := client.NodePools().List(nil)
			if err != nil {
				return err
			}
			ctx["region_nodepool"] = nodepoolList
		}
	}

	return nil
}
