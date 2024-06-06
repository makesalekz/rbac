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

func (u *CheckPermissionsUsecase) getParentIDs(ctx context.Context, tenantID, teamID int64) ([]int64, error) {
	if tenantID != 0 && teamID != 0 {
		team, err := u.teamRepo.GetTeam(ctx, tenantID, teamID, false)
		if err != nil {
			return nil, err
		}

		if team.ParentsIds != nil {
			var parentIDs []int64
			err = team.ParentsIds.AssignTo(&parentIDs)
			if err != nil {
				return nil, err
			}
			return parentIDs, nil
		}
	}
	return nil, nil
}

func (u *CheckPermissionsUsecase) CheckPermissions(
	ctx context.Context,
	tenantID int64,
	identities []string,
	permissions []string,
	resources []*v1.Resource,
) (map[string]*v1.ListOfFields, error) {
	var teamsIDs []int64
	for i := len(resources) - 1; i >= 0; i-- {
		if resources[i].GetType() == data.RESOURCE_TYPE_TEAM {
			parentIDs, err := u.getParentIDs(ctx, tenantID, resources[i].GetId())
			if err != nil {
				return nil, err
			}

			teamsIDs = append(teamsIDs, resources[i].GetId())
			if parentIDs != nil {
				teamsIDs = append(teamsIDs, parentIDs...)
			}

			resources = append(resources[:i], resources[i+1:]...)
		}
	}

	assignedRoles, err := u.repo.CheckRoles(ctx, data.ListRolesDto{
		TenantId:    tenantID,
		IdentityIDs: identities,
		TeamsIDs:    teamsIDs,
		Resources:   resources,
	})
	if err != nil {
		return nil, err
	}

	rolesIDs := make([]int64, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		rolesIDs[i] = assignedRole.RoleID
	}

	rolesPermissions, err := u.roleRepo.ListRolesPermissions(ctx, data.FilterRolePermissions{
		TenantID:    tenantID,
		RolesIDs:    rolesIDs,
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
			result[rolePermission.PermissionID].Fields = mergeFields(
				result[rolePermission.PermissionID].GetFields(),
				rolePermission.Fields,
			)
		}
	}

	for _, rolePermission := range rolesPermissions {
		if rolePermission.Deny {
			delete(result, rolePermission.PermissionID)
		}
	}

	return result, nil
}

func (u *CheckPermissionsUsecase) HasPermission(
	ctx context.Context,
	tenantID int64,
	identities []string,
	permission string,
) (*v1.ListOfFields, error) {
	permissionsMap, err := u.CheckPermissions(ctx, tenantID, identities, []string{permission}, nil)
	if err != nil {
		return nil, err
	}

	if len(permissionsMap) == 0 {
		return &v1.ListOfFields{}, nil
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
