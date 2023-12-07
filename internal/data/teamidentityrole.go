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
	TenantId    int64
	IdentityIDs []string
	TeamsIDs    []int64
}

// TeamIdentityRoleRepo
type TeamIdentityRoleRepo interface {
	AssignRole(ctx context.Context, dto AssignRoleDto) error
	GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.TeamIdentityRole, error)
	DeleteIdentityRole(ctx context.Context, assignedRole *ent.TeamIdentityRole) error
	ListRoles(ctx context.Context, dto ListRolesDto) ([]*ent.TeamIdentityRole, error)
}

type teamIdentityRoleRepo struct {
	db *ent.Client
}

func (t *teamIdentityRoleRepo) AssignRole(ctx context.Context, dto AssignRoleDto) error {
	query := t.db.TeamIdentityRole.Create().
		SetTenantID(dto.TenantId).
		SetIdentityID(dto.IdentityId).
		SetRoleID(dto.RoleId)

	if dto.TeamId != 0 {
		query.SetTeamID(dto.TeamId)
	}

	return query.Exec(ctx)
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

	if len(dto.IdentityIDs) > 0 {
		identityIDs := append(dto.IdentityIDs, "")
		query.Where(teamidentityrole.IdentityIDIn(identityIDs...))
	}

	if len(dto.TeamsIDs) > 0 {
		query.Where(
			teamidentityrole.Or(
				teamidentityrole.TeamIDIn(dto.TeamsIDs...),
				teamidentityrole.TeamIDIsNil(),
			),
		)
	}

	return query.All(ctx)
}

// NewTeamIdentityRoleRepo .
func NewTeamIdentityRoleRepo(d *Data) TeamIdentityRoleRepo {
	return &teamIdentityRoleRepo{
		db: d.db,
	}
}
