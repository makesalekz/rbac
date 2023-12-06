package data

import (
	"context"

	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/ent/teamidentityrole"

	_ "github.com/lib/pq"
)

type AssignRoleDto struct {
	RoleId     int64
	TeamId     int64
	IdentityId string
	TenantId   int64
}

type ListRolesDto struct {
	TeamId      int64
	TenantId    int64
	IdentityIDs []string
}

// TeamIdentityRoleRepo
type TeamIdentityRoleRepo interface {
	AssignRole(ctx context.Context, dto AssignRoleDto) (*ent.TeamIdentityRole, error)
	GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.TeamIdentityRole, error)
	DeleteIdentityRole(ctx context.Context, assignedRole *ent.TeamIdentityRole) error
	ListRoles(ctx context.Context, dto ListRolesDto) ([]*ent.TeamIdentityRole, error)
}

type teamIdentityRoleRepo struct {
	db *ent.Client
}

func (t *teamIdentityRoleRepo) AssignRole(ctx context.Context, dto AssignRoleDto) (*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Create().
		SetTenantID(dto.TenantId).
		SetTeamID(dto.TeamId).
		SetIdentityID(dto.IdentityId).
		SetRoleID(dto.RoleId).
		Save(ctx)
}

func (t *teamIdentityRoleRepo) GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Query().
		Where(
			teamidentityrole.TenantID(tenantId),
			teamidentityrole.ID(assignId),
		).
		Only(ctx)
}

func (t *teamIdentityRoleRepo) DeleteIdentityRole(ctx context.Context, assignedRole *ent.TeamIdentityRole) error {
	return t.db.TeamIdentityRole.DeleteOne(assignedRole).Exec(ctx)
}

func (t *teamIdentityRoleRepo) ListRoles(ctx context.Context, dto ListRolesDto) ([]*ent.TeamIdentityRole, error) {
	query := t.db.TeamIdentityRole.Query().
		Where(teamidentityrole.TenantID(dto.TenantId))

	if dto.TeamId != 0 {
		query.Where(teamidentityrole.TeamID(dto.TeamId))
	}

	if len(dto.IdentityIDs) > 0 {
		query.Where(teamidentityrole.IdentityIDIn(dto.IdentityIDs...))
	}

	return query.All(ctx)
}

// NewTeamIdentityRoleRepo .
func NewTeamIdentityRoleRepo(d *Data) TeamIdentityRoleRepo {
	return &teamIdentityRoleRepo{
		db: d.db,
	}
}
