package data

import (
	"context"
	_ "github.com/lib/pq"
	"rbac/ent"
	"rbac/ent/teamidentityrole"
)

type AssignRoleDto struct {
	RoleId     int64
	TeamId     int64
	IdentityId int64
	TenantId   int64
}

type DeleteRoleDto struct {
	AssignId int64
}

type ListIdentityRolesDto struct {
	IdentityId int64
	TenantId   int64
}

type ListTeamRolesDto struct {
	TeamId   int64
	TenantId int64
}

// TeamIdentityRoleRepo
type TeamIdentityRoleRepo interface {
	AssignRole(ctx context.Context, dto AssignRoleDto) (*ent.TeamIdentityRole, error)
	DeleteIdentityRole(ctx context.Context, dto DeleteRoleDto) error
	ListIdentityRoles(ctx context.Context, dto ListIdentityRolesDto) ([]*ent.TeamIdentityRole, error)
	ListTeamRoles(ctx context.Context, dto ListTeamRolesDto) ([]*ent.TeamIdentityRole, error)
}

type teamIdentityRoleRepo struct {
	db *ent.Client
}

func (t *teamIdentityRoleRepo) ListTeamRoles(ctx context.Context, dto ListTeamRolesDto) ([]*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Query().Where(
		teamidentityrole.TeamID(dto.TeamId),
		teamidentityrole.TenantID(dto.TenantId),
	).All(ctx)
}

func (t *teamIdentityRoleRepo) ListIdentityRoles(ctx context.Context, dto ListIdentityRolesDto) ([]*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Query().Where(
		teamidentityrole.IdentityID(dto.IdentityId),
		teamidentityrole.TenantID(dto.TenantId),
	).All(ctx)
}

func (t *teamIdentityRoleRepo) DeleteIdentityRole(ctx context.Context, dto DeleteRoleDto) error {
	return t.db.TeamIdentityRole.DeleteOneID(dto.AssignId).Exec(ctx)
}

func (t *teamIdentityRoleRepo) AssignRole(ctx context.Context, dto AssignRoleDto) (*ent.TeamIdentityRole, error) {
	return t.db.TeamIdentityRole.Create().
		SetTeamID(dto.TeamId).
		SetIdentityID(dto.IdentityId).
		SetRoleID(dto.RoleId).
		SetTenantID(dto.TenantId).
		Save(ctx)
}

// NewTeamIdentityRoleRepo .
func NewTeamIdentityRoleRepo(d *Data) TeamIdentityRoleRepo {
	return &teamIdentityRoleRepo{
		db: d.db,
	}
}
