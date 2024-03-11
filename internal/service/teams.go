package service

import (
	"context"
	"time"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
)

type TeamsService struct {
	v1.UnimplementedTeamsServer

	sh *ServiceHelper
	tu *biz.TeamsUsecase
}

func NewTeamsService(
	sh *ServiceHelper,
	tu *biz.TeamsUsecase,
) *TeamsService {
	return &TeamsService{
		sh: sh,
		tu: tu,
	}
}

func (s *TeamsService) CreateTeam(ctx context.Context, req *v1.CreateTeamRequest) (*v1.TeamReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.team.create")
	if err != nil {
		return nil, err
	}

	team, err := s.tu.CreateTeam(ctx, data.TeamDto{
		TenantId:    tenantId,
		Name:        req.Name,
		Description: req.Description,
		ParentId:    req.ParentId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) UpdateTeam(ctx context.Context, req *v1.UpdateTeamRequest) (*v1.TeamReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.team.update")
	if err != nil {
		return nil, err
	}

	team, err := s.tu.GetTeam(ctx, tenantId, req.GetTeamId(), false)
	if err != nil {
		return nil, err
	}

	updated, err := s.tu.UpdateTeam(ctx, team, data.TeamDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(updated),
	}, nil
}

func (s *TeamsService) DeleteTeam(ctx context.Context, req *v1.TeamRequest) (*utils_v1.EmptyReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.team.delete")
	if err != nil {
		return nil, err
	}

	team, err := s.tu.GetTeam(ctx, tenantId, req.GetTeamId(), false)
	if err != nil {
		return nil, err
	}

	err = s.tu.DeleteTeam(ctx, team)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *TeamsService) GetTeam(ctx context.Context, req *v1.TeamRequest) (*v1.TeamReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.team.read")
	if err != nil {
		return nil, err
	}

	team, err := s.tu.GetTeam(ctx, tenantId, req.GetTeamId(), req.GetWithTree())
	if err != nil {
		return nil, err
	}

	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) ListTeams(ctx context.Context, req *v1.ListTeamsRequest) (*v1.ListTeamsReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.team.read")
	if err != nil {
		return nil, err
	}

	list, err := s.tu.ListTeams(ctx, data.TeamsListFilter{
		TenantId: tenantId,
		ParentId: req.GetParentId(),
	}, req.GetPaginate())
	if err != nil {
		return nil, err
	}

	return &v1.ListTeamsReply{
		Teams:    replyTeams(list.Teams),
		Paginate: list.Paginate,
	}, nil
}

func replyTeam(team *ent.Team) *v1.Team {
	result := v1.Team{
		Id:          team.ID,
		OwnerId:     team.TenantID,
		ParentId:    team.ParentID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   team.UpdatedAt.Format(time.RFC3339),
	}

	if len(team.Edges.Children) > 0 {
		result.Subs = replyTeams(team.Edges.Children)
	}

	if len(team.ParentsIds.Elements) > 0 {
		_ = team.ParentsIds.AssignTo(&result.ParentsIds)
	}

	return &result
}

func replyTeams(teams []*ent.Team) []*v1.Team {
	result := make([]*v1.Team, len(teams))
	for i, team := range teams {
		result[i] = replyTeam(team)
	}
	return result
}
