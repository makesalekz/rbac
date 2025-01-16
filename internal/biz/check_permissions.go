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
) map[string][]string {
	result := make(map[string][]string)
	for _, rolePermission := range rolePermissions {
		if _, ok := result[rolePermission.PermissionID]; !ok {
			result[rolePermission.PermissionID] = rolePermission.Fields
		} else {
			result[rolePermission.PermissionID] = mergeFields(
				result[rolePermission.PermissionID],
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

func (u *CheckPermissionsUsecase) getPermissionResources(
	rolePermissions []*ent.RolePermission, assignedRoles []*ent.ResourceAccess,
) map[string][]*v1.Resource {
	result := make(map[string][]*v1.Resource)

	roleMap := make(map[int64][]*ent.ResourceAccess)

	for _, role := range assignedRoles {
		roleMap[role.RoleID] = append(roleMap[role.RoleID], role)
	}

	for _, permission := range rolePermissions {
		for _, role := range roleMap[permission.RoleID] {
			var resource *v1.Resource
			//nolint: gocritic // it suggest rewriting to switch, which is not the case
			if role.ResourceType == nil || *role.ResourceType == "" {
				resource = &v1.Resource{
					Type: "*",
					Id:   0,
				}
			} else if role.ResourceID == nil || *role.ResourceID == 0 {
				resource = &v1.Resource{
					Type: *role.ResourceType,
					Id:   0,
				}
			} else {
				resource = &v1.Resource{
					Type: *role.ResourceType,
					Id:   *role.ResourceID,
				}
			}

			result[permission.PermissionID] = append(result[permission.PermissionID], resource)
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

	fields := u.getPermissionFields(rolePermissions, value)
	allowedResources := u.getPermissionResources(rolePermissions, assignedRoles)

	return buildCheckPermissionsReply(fields, allowedResources), nil
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

func buildCheckPermissionsReply(
	fields map[string][]string, resources map[string][]*v1.Resource,
) map[string]*v1.ListOfFields {
	result := make(map[string]*v1.ListOfFields)

	for permissionID, field := range fields {
		if _, ok := result[permissionID]; !ok {
			result[permissionID] = &v1.ListOfFields{
				Fields: field,
			}

			continue
		}

		result[permissionID].Fields = append(result[permissionID].Fields, field...)
	}

	for permissionID, resource := range resources {
		if _, ok := result[permissionID]; !ok {
			result[permissionID] = &v1.ListOfFields{
				Resources: resource,
			}

			continue
		}

		result[permissionID].Resources = append(result[permissionID].Resources, resource...)
	}

	for permissionID, accesses := range result {
		if len(accesses.GetResources()) == 0 {
			delete(result, permissionID)
		}
	}

	return result
}
