package biz

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
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

func (u *CheckPermissionsUsecase) appendTeamParents(
	ctx context.Context,
	tenantID int64,
	resources []*v1.Resource,
) ([]*v1.Resource, error) {
	var teamsIDs []int64
	// extract teams from resources
	for i := len(resources) - 1; i >= 0; i-- {
		if resources[i].GetType() == data.ResourceTypeTeam {
			teamsIDs = append(teamsIDs, resources[i].GetId())
			resources = append(resources[:i], resources[i+1:]...)
		}
	}

	if len(teamsIDs) == 0 {
		return resources, nil
	}

	teams, err := u.teamRepo.GetTeams(ctx, tenantID, teamsIDs)
	if err != nil {
		return nil, err
	}

	// append teams & their parents to resources
	for _, team := range teams {
		resources = append(resources, &v1.Resource{
			Id:   team.ID,
			Type: data.ResourceTypeTeam,
		})

		if team.ParentsIds != nil {
			var parentIDs []int64
			err = team.ParentsIds.AssignTo(&parentIDs)
			if err != nil {
				return nil, err
			}

			for _, parentID := range parentIDs {
				resources = append(resources, &v1.Resource{
					Id:   parentID,
					Type: data.ResourceTypeTeam,
				})
			}
		}
	}

	return resources, nil
}

func (u *CheckPermissionsUsecase) getPermissionFields(
	rolePermissions []*ent.RolePermission, value int64,
) map[string]*v1.ListOfFields {
	result := make(map[string]*v1.ListOfFields)
	for _, rolePermission := range rolePermissions {
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

	for _, rolePermission := range rolePermissions {
		if rolePermission.Deny && value >= rolePermission.Value {
			delete(result, rolePermission.PermissionID)
		}
	}

	return result
}

func (u *CheckPermissionsUsecase) CheckPermissions(
	ctx context.Context, tenantID int64, appID string,
	identities []string, permissions []string, resources []*v1.Resource,
	value int64,
) (map[string]*v1.ListOfFields, error) {
	allResources, err := u.appendTeamParents(ctx, tenantID, resources)
	if err != nil {
		return nil, err
	}

	assignedRoles, err := u.repo.CheckRoles(ctx, data.ListRolesDto{
		TenantID:    tenantID,
		IdentityIDs: identities,
		Resources:   allResources,
	})
	if err != nil {
		return nil, err
	}

	roleIDs := data.ExtractUnique(assignedRoles, func(e *ent.ResourceAccess) (int64, bool) { return e.RoleID, true })

	rolePermissions, err := u.roleRepo.ListRolesPermissions(ctx, data.FilterRolePermissions{
		TenantID:    tenantID,
		RoleIDs:     roleIDs,
		Permissions: permissions,
		AppIDs:      []string{appID, "common", "admin"},
	})
	if err != nil {
		return nil, err
	}

	return u.getPermissionFields(rolePermissions, value), nil
}

func (u *CheckPermissionsUsecase) HasPermission(
	ctx context.Context, tenantID int64, appID string,
	identities []string, permission string,
) (*v1.ListOfFields, error) {
	permissions, err := u.CheckPermissions(ctx, tenantID, appID,
		identities, []string{permission}, nil,
		0)
	if err != nil {
		return nil, err
	}

	if len(permissions) == 0 {
		return &v1.ListOfFields{}, nil
	}

	return permissions[permission], nil
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
