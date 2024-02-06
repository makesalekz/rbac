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

// AssignedRolesRepo
type AssignedRolesRepo interface {
	AssignRole(ctx context.Context, dto AssignRoleDto) error
	UnassignRole(ctx context.Context, assignedRole *ent.TeamIdentityRole) error
	GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.TeamIdentityRole, error)
	ListAssignedRoles(ctx context.Context, dto ListRolesDto) ([]*ent.TeamIdentityRole, error)
}

type assignedRolesRepo struct {
	db *ent.Client
}

func (t *assignedRolesRepo) AssignRole(ctx context.Context, dto AssignRoleDto) error {
	query := t.db.TeamIdentityRole.Create().
		SetTenantID(dto.TenantId).
		SetIdentityID(dto.IdentityId).
		SetRoleID(dto.RoleId)

	if dto.TeamId != 0 {
		query.SetTeamID(dto.TeamId)
	}

	return query.Exec(ctx)
}

func (t *assignedRolesRepo) GetAssignedRoleById(ctx context.Context, tenantId, assignId int64) (*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Query().
		Where(
			teamidentityrole.TenantID(tenantId),
			teamidentityrole.ID(assignId),
		).
		Only(ctx)
}

func (t *assignedRolesRepo) UnassignRole(ctx context.Context, assignedRole *ent.TeamIdentityRole) error {
	return t.db.TeamIdentityRole.DeleteOne(assignedRole).Exec(ctx)
}

func (t *assignedRolesRepo) ListAssignedRoles(ctx context.Context, dto ListRolesDto) ([]*ent.TeamIdentityRole, error) {
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

// NewAssignedRolesRepo .
func NewAssignedRolesRepo(d *Data) AssignedRolesRepo {
	return &assignedRolesRepo{
		db: d.db,
	}
}
