package data

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgtype"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/team"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"

	_ "github.com/lib/pq"
)

// TeamsRepo.
type TeamsRepo interface {
	CreateTeam(ctx context.Context, dto TeamDto) (*ent.Team, error)
	UpdateTeam(ctx context.Context, team *ent.Team, dto TeamDto) (*ent.Team, error)
	DeleteTeam(ctx context.Context, team *ent.Team) error
	GetTeam(ctx context.Context, tenantID, teamID int64, getTree bool) (*ent.Team, error)
	GetTeams(ctx context.Context, tenantID int64, teamIDs []int64) ([]*ent.Team, error)
	ListTeams(ctx context.Context, filter TeamsListFilter, paginate *utils_v1.PaginateRequest) ([]*ent.Team, error)
	CountListTeams(ctx context.Context, filter TeamsListFilter) (int32, error)
}

type teamsRepo struct {
	db *ent.Client
}

// NewTeamsRepo .
func NewTeamsRepo(d *Data) TeamsRepo {
	return &teamsRepo{
		db: d.db,
	}
}

func (r *teamsRepo) CreateTeam(ctx context.Context, dto TeamDto) (*ent.Team, error) {
	query := r.db.Team.Create().
		SetTenantID(dto.TenantID).
		SetName(dto.Name).
		SetDescription(dto.Description)

	if dto.ParentID != 0 {
		query.SetParentID(dto.ParentID)
	}

	parentsIDs := &pgtype.Int8Array{Status: pgtype.Present}

	if len(dto.ParentsIDs) > 0 {
		err := parentsIDs.Set(dto.ParentsIDs)
		if err != nil {
			return nil, err
		}
	}
	return query.SetParentsIds(parentsIDs).Save(ctx)
}

func (r *teamsRepo) UpdateTeam(ctx context.Context, team *ent.Team, dto TeamDto) (*ent.Team, error) {
	return r.db.Team.UpdateOne(team).
		SetName(dto.Name).
		SetDescription(dto.Description).
		Save(ctx)
}

func (r *teamsRepo) DeleteTeam(ctx context.Context, team *ent.Team) error {
	return r.db.Team.DeleteOne(team).Exec(ctx)
}

func (r *teamsRepo) GetTeam(ctx context.Context, tenantID, teamID int64, getTree bool) (*ent.Team, error) {
	team, err := r.db.Team.Query().Where(team.ID(teamID), team.TenantID(tenantID)).Only(ctx)
	if err != nil {
		return nil, err
	}

	if getTree {
		subs, err := r.db.Team.Query().Where(func(s *sql.Selector) {
			s.Where(sql.ExprP("parents_ids @> ARRAY[$1]::bigint[]", teamID))
		}).All(ctx)
		if err != nil {
			return nil, err
		}

		findChildren(team, subs)
	}

	return team, err
}

func (r *teamsRepo) GetTeams(ctx context.Context, tenantID int64, teamIDs []int64) ([]*ent.Team, error) {
	return r.db.Team.Query().
		Where(
			team.IDIn(teamIDs...),
			team.TenantID(tenantID),
		).All(ctx)
}

func findChildren(team *ent.Team, subs []*ent.Team) {
	for _, sub := range subs {
		if *sub.ParentID == team.ID {
			team.Edges.Children = append(team.Edges.Children, sub)
			findChildren(sub, subs)
		}
	}
}

func (r *teamsRepo) ListTeams(
	ctx context.Context,
	filter TeamsListFilter,
	paginate *utils_v1.PaginateRequest,
) ([]*ent.Team, error) {
	query := r.db.Team.Query()

	if filter.TenantID != 0 {
		query.Where(team.TenantID(filter.TenantID))
	}

	if filter.ParentID != 0 {
		query.Where(team.ParentID(filter.ParentID))
	}

	if paginate.GetFromId() != 0 {
		query.Where(team.IDGT(paginate.GetFromId()))
	}

	if paginate.GetLimit() == 0 {
		paginate.Limit = 100
	}

	return query.Limit(int(paginate.GetLimit())).Order(ent.Asc(team.FieldID)).All(ctx)
}

func (r *teamsRepo) CountListTeams(ctx context.Context, filter TeamsListFilter) (int32, error) {
	query := r.db.Team.Query()

	if filter.TenantID != 0 {
		query.Where(team.TenantID(filter.TenantID))
	}

	if filter.ParentID != 0 {
		query.Where(team.ParentID(filter.ParentID))
	}

	count, err := query.Count(ctx)

	return int32(count), err
}
