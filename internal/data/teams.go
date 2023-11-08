package data

import (
	"context"

	teams_v1 "rbac/api/rbac/v1"
	"rbac/ent"
	"rbac/ent/team"

	"entgo.io/ent/dialect/sql"
	"github.com/jackc/pgtype"
	_ "github.com/lib/pq"
)

type TeamDto struct {
	TenantId    int64
	Name        string
	Description string
	ParentId    int64
	ParentsIds  []int64
}

type TeamsListFilter struct {
	TenantId int64
	ParentId int64
}

// TeamsRepo
type TeamsRepo interface {
	CreateTeam(ctx context.Context, dto TeamDto) (*ent.Team, error)
	UpdateTeam(ctx context.Context, teamId int64, dto TeamDto) (*ent.Team, error)
	DeleteTeam(ctx context.Context, teamId, tenantId int64) error
	GetTeam(ctx context.Context, teamId, tenantId int64, getTree bool) (*ent.Team, error)
	ListTeams(ctx context.Context, filter TeamsListFilter, paginate *teams_v1.PaginateRequest) ([]*ent.Team, error)
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
		SetTenantID(dto.TenantId).
		SetName(dto.Name).
		SetDescription(dto.Description)

	if dto.ParentId != 0 {
		query.SetParentID(dto.ParentId)
	}

	parentsIds := &pgtype.Int8Array{Status: pgtype.Present}

	if len(dto.ParentsIds) > 0 {
		parentsIds.Set(dto.ParentsIds)
	}
	return query.SetParentsIds(parentsIds).Save(ctx)
}

func (r *teamsRepo) UpdateTeam(ctx context.Context, teamId int64, dto TeamDto) (*ent.Team, error) {
	return r.db.Team.UpdateOneID(teamId).
		Where(team.TenantID(dto.TenantId)).
		SetName(dto.Name).
		SetDescription(dto.Description).
		Save(ctx)
}

func (r *teamsRepo) DeleteTeam(ctx context.Context, teamId, tenantId int64) error {
	return r.db.Team.DeleteOneID(teamId).Where(team.TenantID(tenantId)).Exec(ctx)
}

func (r *teamsRepo) GetTeam(ctx context.Context, teamId, tenantId int64, getTree bool) (*ent.Team, error) {
	team, err := r.db.Team.Query().Where(team.ID(teamId), team.TenantID(tenantId)).Only(ctx)
	if err != nil {
		return nil, err
	}

	if getTree {
		subs, err := r.db.Team.Query().Where(func(s *sql.Selector) {
			s.Where(sql.ExprP("parents_ids @> ARRAY[$1]::bigint[]", teamId))
		}).All(ctx)
		if err != nil {
			return nil, err
		}

		findChildren(team, subs)
	}

	return team, err
}

func findChildren(team *ent.Team, subs []*ent.Team) {
	for _, sub := range subs {
		if *sub.ParentID == team.ID {
			team.Edges.Children = append(team.Edges.Children, sub)
			findChildren(sub, subs)
		}
	}
}

func (r *teamsRepo) ListTeams(ctx context.Context, filter TeamsListFilter, paginate *teams_v1.PaginateRequest) ([]*ent.Team, error) {
	query := r.db.Team.Query()

	if filter.TenantId != 0 {
		query.Where(team.TenantID(filter.TenantId))
	}

	if filter.ParentId != 0 {
		query.Where(team.ParentID(filter.ParentId))
	}

	if paginate.FromId != 0 {
		query.Where(team.IDGT(paginate.FromId))
	}

	if paginate.Limit == 0 {
		paginate.Limit = 100
	}

	return query.Limit(int(paginate.Limit)).Order(ent.Asc(team.FieldID)).All(ctx)
}

func (r *teamsRepo) CountListTeams(ctx context.Context, filter TeamsListFilter) (int32, error) {
	query := r.db.Team.Query()

	if filter.TenantId != 0 {
		query.Where(team.TenantID(filter.TenantId))
	}

	if filter.ParentId != 0 {
		query.Where(team.ParentID(filter.ParentId))
	}

	count, err := query.Count(ctx)

	return int32(count), err
}
