package biz

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

// AssignedRolesUsecase .
type AssignedRolesUsecase struct {
	repo     data.AssignedRolesRepo
	roleRepo data.RoleRepo
	teamRepo data.TeamsRepo
}

// NewAssignedRolesUsecase .
func NewAssignedRolesUsecase(
	repo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamRepo data.TeamsRepo,
) (*AssignedRolesUsecase, error) {
	return &AssignedRolesUsecase{
		repo:     repo,
		roleRepo: roleRepo,
		teamRepo: teamRepo,
	}, nil
}

func (u *AssignedRolesUsecase) AssignRole(ctx context.Context, dto data.AssignRoleDto) error {
	_, err := u.roleRepo.GetRoleById(ctx, dto.TenantId, dto.RoleId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("role not found")
		}
		return v1.ErrorDatabaseQuery("get role failed: %v", err)
	}

	if dto.TeamId != 0 {
		_, err = u.teamRepo.GetTeam(ctx, dto.TenantId, dto.TeamId, false)
		if err != nil {
			if ent.IsNotFound(err) {
				return v1.ErrorNotFound("team not found")
			}
			return v1.ErrorDatabaseQuery("get team failed: %v", err)
		}
	}

	err = u.repo.AssignRole(ctx, dto)
	if err != nil {
		if ent.IsConstraintError(err) {
			return v1.ErrorAlreadyExists("role already assigned")
		}
		return v1.ErrorDatabaseQuery("assign role failed: %v", err)
	}

	return nil
}

func (u *AssignedRolesUsecase) UnassignRole(ctx context.Context, tenantId, assignId int64) error {
	assignedRole, err := u.repo.GetAssignedRoleById(ctx, tenantId, assignId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("assigned role not found")
		}
		return v1.ErrorDatabaseQuery("get assigned role failed: %v", err)
	}

	return u.repo.UnassignRole(ctx, assignedRole)
}

func (u *AssignedRolesUsecase) ListIdentityRoles(ctx context.Context, tenantId int64, identityId string) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: []string{identityId},
	})
}

func (u *AssignedRolesUsecase) ListAssignedRoles(ctx context.Context, dto data.ListRolesDto) ([]*ent.TeamIdentityRole, error) {
	return u.repo.ListAssignedRoles(ctx, dto)
}
