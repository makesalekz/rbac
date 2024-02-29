package biz

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

// CheckPermissionsUsecase .
type CheckPermissionsUsecase struct {
	repo     data.AssignedRolesRepo
	roleRepo data.RoleRepo
	teamRepo data.TeamsRepo
}

// NewCheckPermissionsUsecase .
func NewCheckPermissionsUsecase(
	repo data.AssignedRolesRepo,
	roleRepo data.RoleRepo,
	teamRepo data.TeamsRepo,
) (*CheckPermissionsUsecase, error) {
	return &CheckPermissionsUsecase{
		repo:     repo,
		roleRepo: roleRepo,
		teamRepo: teamRepo,
	}, nil
}

func (u *CheckPermissionsUsecase) CheckPermissions(ctx context.Context, tenantId int64, identities []string, teamId int64, permissions []string) (map[string]*v1.ListOfFields, error) {
	var teamsIds []int64
	if tenantId != 0 && teamId != 0 {
		team, err := u.teamRepo.GetTeam(ctx, tenantId, teamId, false)
		if err != nil {
			return nil, err
		}

		if team.ParentsIds != nil {
			err = team.ParentsIds.AssignTo(&teamsIds)
			if err != nil {
				return nil, err
			}
		}

		teamsIds = append(teamsIds, team.ID)
	}

	assignedRoles, err := u.repo.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: identities,
		TeamsIDs:    teamsIds,
	})
	if err != nil {
		return nil, err
	}

	rolesIds := make([]int64, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		rolesIds[i] = assignedRole.RoleID
	}

	rolesPermissions, err := u.roleRepo.ListRolesPermissions(ctx, data.FilterRolePermissions{
		TenantId:    tenantId,
		RolesIds:    rolesIds,
		Permissions: permissions,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]*v1.ListOfFields)
	for _, rolePermission := range rolesPermissions {
		if _, ok := result[rolePermission.PermissionID]; !ok {
			result[rolePermission.PermissionID] = &v1.ListOfFields{
				Fields: rolePermission.Fields,
			}
		} else {
			result[rolePermission.PermissionID].Fields = mergeFields(result[rolePermission.PermissionID].Fields, rolePermission.Fields)
		}
	}

	for _, rolePermission := range rolesPermissions {
		if rolePermission.Deny {
			delete(result, rolePermission.PermissionID)
		}
	}

	return result, nil
}

func (u *CheckPermissionsUsecase) HasPermission(ctx context.Context, tenantId int64, identities []string, permission string) (*v1.ListOfFields, error) {
	permissionsMap, err := u.CheckPermissions(ctx, tenantId, identities, 0, []string{permission})
	if err != nil {
		return nil, err
	}

	if len(permissionsMap) == 0 {
		return nil, nil
	}

	return permissionsMap[permission], nil
}

func mergeFields(fields1 []string, fields2 []string) []string {
	if len(fields1) > len(fields2) {
		fields1, fields2 = fields2, fields1
	}
	for _, field := range fields2 {
		if !contains(fields1, field) {
			fields1 = append(fields1, field)
		}
	}
	return fields1
}

func contains(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}
