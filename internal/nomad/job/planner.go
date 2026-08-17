// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"errors"
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/domain"
	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/register/region/picker"
	pickercontext "github.com/rasorp/attila/internal/register/region/picker/context"
	"github.com/rasorp/attila/internal/store"
	jobsdk "github.com/rasorp/attila/pkg/job"
)

type Planner struct {
	logger *zap.Logger

	clients *client.Clients
	job     *api.Job
	state   store.State

	plan *domain.JobRegisterPlan
}

type PlannerReq struct {
	Clients *client.Clients
	Job     *api.Job
	State   store.State
}

func NewPlanner(logger *zap.Logger, req *PlannerReq) *Planner {
	return &Planner{
		clients: req.Clients,
		job:     req.Job,
		logger: logger.With(
			zap.String("job_id", *req.Job.ID),
			zap.String("job_namespace", *req.Job.Namespace),
		).Named("job_plan"),
		plan:  domain.NewJobRegisterPlan(*req.Job.ID, *req.Job.Namespace),
		state: req.State,
	}
}

func (p *Planner) Run() (*domain.JobRegisterPlan, error) {
	listResp, err := p.state.JobRegister().Method().List(nil)
	if err != nil {
		return nil, err
	}
	if len(listResp.Methods) == 0 {
		return nil, errors.New("found zero job register methods")
	}

	var ruleLinks []*domain.JobRegisterMethodRuleLink

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

	var rules []*domain.JobRegisterRule

	for _, ruleLink := range ruleLinks {
		regRule, err := p.state.JobRegister().Rule().Get(&store.JobRegisterRuleGetReq{Name: ruleLink.Name})
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
		pickedRegions, err := p.runRegisterPlanPicker(rule, regionListResp.Regions)
		if err != nil {
			return nil, err
		}

		if err := p.generatePlanResult(rule.Name, pickedRegions); err != nil {
			return nil, err
		}
	}

	return p.plan, nil
}

// runRegisterPlanPicker executes the job registration plan strategy pipeline and
// returns the set of regions selected for plan generation.
func (p *Planner) runRegisterPlanPicker(
	rule *domain.JobRegisterRule, regions []*domain.Region) ([]*domain.Region, error) {
	p.logger.Debug(
		"performing execution of rule region picker",
		zap.String("rule_name", rule.Name),
		zap.Int("num_regions", len(regions)),
	)

	if len(rule.RegionPickers) == 0 {
		return regions, nil
	}

	candidates, err := pickercontext.BuildCandidates(regions, func(region *domain.Region) (map[string]any, error) {
		regionClient, err := p.clients.Get(region.Name)
		if err != nil {
			return nil, err
		}

		ctx := make(map[string]any)
		if err := populateRegionContext(rule, regionClient, ctx); err != nil {
			return nil, err
		}
		return ctx, nil
	})
	if err != nil {
		return nil, err
	}

	picker, err := picker.New(rule.RegionPickers)
	if err != nil {
		return nil, err
	}

	pickedCandidates, err := picker.Process(&jobsdk.RegionPickerRunRequest{
		Job:              p.job,
		Rule:             domainJobRegisterRuleToPickerRule(rule),
		RegionCandidates: candidates,
	})
	if err != nil {
		return nil, err
	}

	regionByName := make(map[string]*domain.Region, len(regions))
	for _, region := range regions {
		regionByName[region.Name] = region
	}

	pickedRegions := make([]*domain.Region, 0, len(pickedCandidates))
	for _, candidate := range pickedCandidates {
		region, ok := regionByName[candidate.Name]
		if !ok {
			return nil, fmt.Errorf("strategy selected unknown region %q", candidate.Name)
		}
		pickedRegions = append(pickedRegions, region)
	}

	return pickedRegions, nil
}

// generatePlanResult iterates the selected region slice and perform a Nomad
// job plan for each. The Nomad plan and region name will then be added to the
// Attila plan result.
//
// Any failure in calling the Nomad API will result in a failure of the whole
// function.
func (p *Planner) generatePlanResult(ruleName string, regions []*domain.Region) error {
	for _, pickedRegion := range regions {
		nomadClient, err := p.clients.Get(pickedRegion.Name)
		if err != nil {
			return fmt.Errorf("failed to get Nomad client, %w", err)
		}

		// TODO(jrasem): add support for job plan diff.
		planResp, _, err := nomadClient.Jobs().PlanOpts(p.job, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to call Nomad job plan, %w", err)
		}

		p.plan.AddRegion(pickedRegion, planResp)

		p.logger.Info(
			"region picked by rule picker",
			zap.String("rule_name", ruleName),
			zap.String("region_name", pickedRegion.Name),
		)
	}

	return nil
}

func domainJobRegisterRuleToPickerRule(rule *domain.JobRegisterRule) *jobsdk.RegionPickerRule {

	regionContexts := make([]string, 0, len(rule.RegionContexts))
	for _, regionContext := range rule.RegionContexts {
		regionContexts = append(regionContexts, regionContext.Kind)
	}

	var metadata *jobsdk.RegionPickerMetadata
	if rule.Metadata != nil {
		metadata = &jobsdk.RegionPickerMetadata{
			CreateTime: rule.Metadata.CreateTime,
			UpdateTime: rule.Metadata.UpdateTime,
		}
	}

	return &jobsdk.RegionPickerRule{
		Name:           rule.Name,
		RegionContexts: regionContexts,
		RegionPickers:  rule.RegionPickers,
		Metadata:       metadata,
	}
}

func populateRegionContext(
	rule *domain.JobRegisterRule, client *api.Client, ctx map[string]any) error {
	for _, regionContext := range rule.RegionContexts {
		switch regionContext.Kind {
		case domain.JobRegisterRuleContextKindNamespace:
			namespaceList, _, err := client.Namespaces().List(nil)
			if err != nil {
				return err
			}
			ctx["region_namespace"] = namespaceList

		case domain.JobRegisterRuleContextKindNodepool:
			nodepoolList, _, err := client.NodePools().List(nil)
			if err != nil {
				return err
			}
			ctx["region_nodepool"] = nodepoolList
		}
	}

	return nil
}
