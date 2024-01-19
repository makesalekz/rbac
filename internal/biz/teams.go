package biz

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
)

type TeamsList struct {
	Teams    []*ent.Team
	Paginate *utils_v1.PaginateReply
}

// TeamsUsecase is a Greeter usecase.
type TeamsUsecase struct {
	repo data.TeamsRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewTeamsUsecase(repo data.TeamsRepo) (*TeamsUsecase, error) {
	return &TeamsUsecase{
		repo: repo,
	}, nil
}

func (uc *TeamsUsecase) CreateTeam(ctx context.Context, dto data.TeamDto) (*ent.Team, error) {
	if dto.ParentId != 0 {
		parentTeam, err := uc.repo.GetTeam(ctx, dto.ParentId, dto.TenantId, false)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, v1.ErrorNotFound("parent team not found")
			}
			return nil, err
		}

		var parentsIds []int64
		parentTeam.ParentsIds.AssignTo(&parentsIds)
		dto.ParentsIds = append(parentsIds, parentTeam.ID)
	}

	return uc.repo.CreateTeam(ctx, dto)
}

func (uc *TeamsUsecase) UpdateTeam(ctx context.Context, team *ent.Team, dto data.TeamDto) (*ent.Team, error) {
	return uc.repo.UpdateTeam(ctx, team, dto)
}

func (uc *TeamsUsecase) DeleteTeam(ctx context.Context, team *ent.Team) error {
	return uc.repo.DeleteTeam(ctx, team)
}

func (uc *TeamsUsecase) GetTeam(ctx context.Context, tenantId, teamId int64, getTree bool) (*ent.Team, error) {
	team, err := uc.repo.GetTeam(ctx, tenantId, teamId, getTree)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("team not found")
		}
		return nil, err
	}

	return team, nil
}

func (uc *TeamsUsecase) ListTeams(ctx context.Context, filter data.TeamsListFilter, paginate *utils_v1.PaginateRequest) (*TeamsList, error) {
	if paginate == nil {
		paginate = &utils_v1.PaginateRequest{}
	}

	teams, err := uc.repo.ListTeams(ctx, filter, paginate)
	if err != nil {
		return nil, err
	}

	total, err := uc.repo.CountListTeams(ctx, filter)
	if err != nil {
		return nil, err
	}

	paginateReply := utils_v1.PaginateReply{
		Total: &total,
	}

	if len(teams) == int(paginate.Limit) {
		paginateReply.FromId = &teams[len(teams)-1].ID
	}

	return &TeamsList{
		Teams:    teams,
		Paginate: &paginateReply,
	}, nil
}
