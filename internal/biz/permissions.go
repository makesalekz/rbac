package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

// PermissionsUsecase .
type PermissionsUsecase struct {
	permissionRepo data.PermissionRepo
	roleRepo       data.RoleRepo
	assignedRepo   data.AssignedRolesRepo
}

// NewPermissionsUsecase .
func NewPermissionsUsecase(
	logger log.Logger,
	permissionRepo data.PermissionRepo,
	roleRepo data.RoleRepo,
	assignedRepo data.AssignedRolesRepo,
) (*PermissionsUsecase, error) {
	return &PermissionsUsecase{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		assignedRepo:   assignedRepo,
	}, nil
}

func (uc *PermissionsUsecase) GetPermissionByID(ctx context.Context, permissionID string) (*ent.Permission, error) {
	permission, err := uc.permissionRepo.GetPermissionByID(ctx, permissionID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("permission not found")
		}
		return nil, err
	}

	return permission, nil
}

func (uc *PermissionsUsecase) CreatePermission(
	ctx context.Context,
	data data.CreatePermissionDto,
) (*ent.Permission, error) {
	return uc.permissionRepo.CreatePermission(ctx, data)
}

func (uc *PermissionsUsecase) UpdatePermission(
	ctx context.Context,
	permissionID string,
	data data.UpdatePermissionDto,
) (*ent.Permission, error) {
	return uc.permissionRepo.UpdatePermission(ctx, permissionID, data)
}

func (uc *PermissionsUsecase) DeletePermission(ctx context.Context, permissionID string) error {
	return uc.permissionRepo.DeletePermission(ctx, permissionID)
}

func (uc *PermissionsUsecase) GetPermissions(
	ctx context.Context,
	appID string,
	permissionIDs []string,
) ([]*ent.Permission, error) {
	return uc.permissionRepo.GetPermissions(ctx, appID, permissionIDs)
}

func (uc *PermissionsUsecase) GetDeniedPermissions(
	ctx context.Context,
	tenantID int64,
	identities []string,
) (map[string]bool, error) {
	assignedRoles, err := uc.assignedRepo.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    tenantID,
		IdentityIDs: identities,
	})
	if err != nil {
		return nil, err
	}

	rolesIDs := make([]int64, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		rolesIDs[i] = assignedRole.RoleID
	}

	permissions, err := uc.roleRepo.ListRolesPermissions(ctx, data.FilterRolePermissions{
		TenantId:   tenantID,
		RolesIds:   rolesIDs,
		DeniedOnly: true,
	})
	if err != nil {
		return nil, err
	}

	excludePermissions := make(map[string]bool)
	for _, permission := range permissions {
		excludePermissions[permission.PermissionID] = true
	}

	return excludePermissions, nil
}

func (uc *PermissionsUsecase) GetGroupedPermissions(
	ctx context.Context,
	tenantID int64,
	identities []string,
	filter data.FilterPermissions,
) ([]*ent.PermissionGroup, error) {
	groups, err := uc.permissionRepo.GetGroupedPermissions(ctx, filter)
	if err != nil {
		return nil, err
	}

	if !filter.WithDenied {
		excludePermissions, err := uc.GetDeniedPermissions(ctx, tenantID, identities)
		if err != nil {
			return nil, err
		}

		// filter denied permissions & empty groups
		for k := len(groups) - 1; k >= 0; k-- {
			group := groups[k]
			for i := len(group.Edges.Permissions) - 1; i >= 0; i-- {
				if _, ok := excludePermissions[group.Edges.Permissions[i].ID]; ok {
					group.Edges.Permissions = append(group.Edges.Permissions[:i], group.Edges.Permissions[i+1:]...)
				}
			}

			if len(group.Edges.Permissions) == 0 {
				groups = append(groups[:k], groups[k+1:]...)
			}
		}
	}

	return groups, nil
}
