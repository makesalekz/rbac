package service

import (
	"context"
	"time"

	v1 "rbac/api/rbac/v1"
	"rbac/ent"
	"rbac/internal/biz"
	"rbac/internal/data"
)

type TeamsService struct {
	v1.UnimplementedTeamsServer

	tu *biz.TeamsUsecase
}

func NewTeamsService(tu *biz.TeamsUsecase) *TeamsService {
	return &TeamsService{
		tu: tu,
	}
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
		team.ParentsIds.AssignTo(&result.ParentsIds)
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

func (s *TeamsService) CreateTeam(ctx context.Context, req *v1.CreateTeamRequest) (*v1.TeamReply, error) {
	team, err := s.tu.CreateTeam(ctx, data.TeamDto{
		Name:        req.Name,
		Description: req.Description,
		ParentId:    req.ParentId,
	})
	if err != nil {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) UpdateTeam(ctx context.Context, req *v1.UpdateTeamRequest) (*v1.TeamReply, error) {
	team, err := s.tu.UpdateTeam(ctx, req.TeamId, data.TeamDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("Team not found")
		}
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) DeleteTeam(ctx context.Context, req *v1.TeamRequest) (*v1.EmptyReply, error) {
	err := s.tu.DeleteTeam(ctx, req.TeamId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("Team not found")
		}
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}

func (s *TeamsService) GetTeam(ctx context.Context, req *v1.TeamRequest) (*v1.TeamReply, error) {
	team, err := s.tu.GetTeam(ctx, req.TeamId, false)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("Team not found")
		}
		return nil, err
	}

	result := &v1.TeamReply{
		Team: replyTeam(team),
	}

	return result, nil
}

func (s *TeamsService) GetTeamTree(ctx context.Context, req *v1.TeamRequest) (*v1.TeamReply, error) {
	team, err := s.tu.GetTeam(ctx, req.TeamId, true)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("Team not found")
		}
		return nil, err
	}
	return &v1.TeamReply{
		Team: replyTeam(team),
	}, nil
}

func (s *TeamsService) ListTeams(ctx context.Context, req *v1.ListTeamsRequest) (*v1.ListTeamsReply, error) {
	list, err := s.tu.ListTeams(ctx, data.TeamsListFilter{
		TenantId: req.OwnerId,
		ParentId: req.ParentId,
	}, req.Paginate)
	if err != nil {
		return nil, err
	}
	return &v1.ListTeamsReply{
		Teams:    replyTeams(list.Teams),
		Paginate: list.Paginate,
	}, nil
}
