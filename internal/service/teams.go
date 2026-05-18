package service

import (
	"context"
	"time"

	v1 "github.com/makesalekz/rbac/api/rbac/v1"
	"github.com/makesalekz/rbac/ent"
	"github.com/makesalekz/rbac/internal/biz"
	"github.com/makesalekz/rbac/internal/data"
	utils_v1 "github.com/makesalekz/utils/api/utils/v1"
	"github.com/makesalekz/utils/v2/auth"
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
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	dto := data.TeamDto{
		TenantID:    tenantID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		ParentID:    req.GetParentId(),
	}
	err := dto.Validate()
	if err != nil {
		return nil, err
	}

	team, err := s.tu.CreateTeam(ctx, dto)
	if err != nil {
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) UpdateTeam(ctx context.Context, req *v1.UpdateTeamRequest) (*v1.TeamReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetTeamId() == 0 {
		return nil, v1.ErrorBadRequest("empty team id")
	}

	team, err := s.tu.GetTeam(ctx, tenantID, req.GetTeamId(), false)
	if err != nil {
		return nil, err
	}

	updated, err := s.tu.UpdateTeam(ctx, team, data.TeamDto{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(updated),
	}, nil
}

func (s *TeamsService) DeleteTeam(ctx context.Context, req *v1.TeamRequest) (*utils_v1.EmptyReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetTeamId() == 0 {
		return nil, v1.ErrorBadRequest("empty team id")
	}

	team, err := s.tu.GetTeam(ctx, tenantID, req.GetTeamId(), false)
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
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetTeamId() == 0 {
		return nil, v1.ErrorBadRequest("empty team id")
	}

	team, err := s.tu.GetTeam(ctx, tenantID, req.GetTeamId(), req.GetWithTree())
	if err != nil {
		return nil, err
	}

	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) GetTeams(ctx context.Context, req *v1.GetTeamsRequest) (*v1.GetTeamsReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	teams, err := s.tu.GetTeams(ctx, tenantID, req.GetTeamIds())
	if err != nil {
		return nil, err
	}

	return &v1.GetTeamsReply{Teams: replyTeams(teams)}, nil
}

func (s *TeamsService) ListTeams(ctx context.Context, req *v1.ListTeamsRequest) (*v1.ListTeamsReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	list, err := s.tu.ListTeams(ctx, data.TeamsListFilter{
		TenantID: tenantID,
		ParentID: req.GetParentId(),
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

	if team.ParentsIds != nil && len(team.ParentsIds.Elements) > 0 {
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
