package biz

import (
	"context"

	consul "github.com/go-kratos/consul/registry"
	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/conf"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

type TeamsList struct {
	Teams    []*ent.Team
	Paginate *v1.PaginateReply
}

// TeamsUsecase is a Greeter usecase.
type TeamsUsecase struct {
	conf      *conf.Bootstrap
	log       *log.Helper
	discovery *consul.Registry
	jwt       *data.JwtProcessor
	repo      data.TeamsRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewTeamsUsecase(logger log.Logger, c *data.Config, jwt *data.JwtProcessor, repo data.TeamsRepo) (*TeamsUsecase, error) {
	return &TeamsUsecase{
		conf:      c.Bootstrap,
		log:       log.NewHelper(logger),
		discovery: c.GetRegistry(),
		jwt:       jwt,
		repo:      repo,
	}, nil
}

func (uc *TeamsUsecase) CreateTeam(ctx context.Context, dto data.TeamDto) (*ent.Team, error) {
	claims, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}

	if dto.TenantId == claims.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions for create team

	if dto.ParentId != 0 {
		parentTeam, err := uc.repo.GetTeam(ctx, dto.ParentId, claims.TenantId, false)
		if err != nil {
			return nil, err
		}
		dto.TenantId = parentTeam.TenantID

		var parentsIds []int64
		parentTeam.ParentsIds.AssignTo(&parentsIds)
		dto.ParentsIds = append(parentsIds, parentTeam.ID)

		orgTeam := parentTeam
		orgId := dto.ParentsIds[0]
		if orgId != dto.ParentId {
			orgTeam, err = uc.repo.GetTeam(ctx, orgId, claims.TenantId, false)
			if err != nil {
				return nil, err
			}

			// TODO: use rbac
			if orgTeam.TenantID != claims.TenantId {
				return nil, v1.ErrorForbidden("Forbidden")
			}
		}
	}

	return uc.repo.CreateTeam(ctx, dto)
}

func (uc *TeamsUsecase) UpdateTeam(ctx context.Context, teamId int64, dto data.TeamDto) (*ent.Team, error) {
	claims, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}

	if dto.TenantId == claims.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions for delete team
	return uc.repo.UpdateTeam(ctx, teamId, dto)
}

func (uc *TeamsUsecase) DeleteTeam(ctx context.Context, teamId int64) error {
	claims, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	_, err := uc.repo.GetTeam(ctx, teamId, claims.TenantId, false)
	if err != nil {
		return err
	}
	// todo checkPermissions for delete team
	return uc.repo.DeleteTeam(ctx, teamId, claims.TenantId)
}

func (uc *TeamsUsecase) GetTeam(ctx context.Context, teamId int64, getTree bool) (*ent.Team, error) {
	claims, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	return uc.repo.GetTeam(ctx, teamId, claims.TenantId, getTree)
}

func (uc *TeamsUsecase) ListTeams(ctx context.Context, filter data.TeamsListFilter, paginate *v1.PaginateRequest) (*TeamsList, error) {
	claims, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if claims.TenantId != filter.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	if paginate == nil {
		paginate = &v1.PaginateRequest{}
	}

	teams, err := uc.repo.ListTeams(ctx, filter, paginate)
	if err != nil {
		return nil, err
	}

	total, err := uc.repo.CountListTeams(ctx, filter)
	if err != nil {
		return nil, err
	}

	paginateReply := v1.PaginateReply{
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
