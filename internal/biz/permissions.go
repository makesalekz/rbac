package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	rbac_v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

// PermissionsUsecase .
type PermissionsUsecase struct {
	log            *log.Helper
	jwt            *jwt.JwtProcessor
	permissionRepo data.PermissionRepo
	roleRepo       data.RoleRepo
	assignedRepo   data.TeamIdentityRoleRepo
}

// NewPermissionUsecase .
func NewPermissionUsecase(
	logger log.Logger,
	jwt *jwt.JwtProcessor,
	permissionRepo data.PermissionRepo,
	roleRepo data.RoleRepo,
	assignedRepo data.TeamIdentityRoleRepo,
) (*PermissionsUsecase, error) {
	return &PermissionsUsecase{
		log:            log.NewHelper(logger),
		jwt:            jwt,
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		assignedRepo:   assignedRepo,
	}, nil
}

func (uc *PermissionsUsecase) GetPermissionById(ctx context.Context, permissionId string) (*ent.Permission, error) {
	return uc.permissionRepo.GetPermissionById(ctx, permissionId)
}

func (uc *PermissionsUsecase) CreatePermission(ctx context.Context, data data.CreatePermissionDto) (*ent.Permission, error) {
	return uc.permissionRepo.CreatePermission(ctx, data)
}

func (uc *PermissionsUsecase) UpdatePermission(ctx context.Context, permissionId string, data data.UpdatePermissionDto) (*ent.Permission, error) {
	return uc.permissionRepo.UpdatePermission(ctx, permissionId, data)
}

func (uc *PermissionsUsecase) DeletePermission(ctx context.Context, permissionId string) error {
	return uc.permissionRepo.DeletePermission(ctx, permissionId)
}

func (uc *PermissionsUsecase) GetPermissions(ctx context.Context, appId string, permissionIds []string) ([]*ent.Permission, error) {
	return uc.permissionRepo.GetPermissions(ctx, appId, permissionIds)
}

func (uc *PermissionsUsecase) GetGroupedPermissions(ctx context.Context, filter data.FilterPermissions) ([]*ent.PermissionGroup, error) {
	groups, err := uc.permissionRepo.GetGroupedPermissions(ctx, filter)
	if err != nil {
		return nil, err
	}

	if !filter.WithDenied {
		// filter permissions, that denied for the user
		claims, ok := uc.jwt.GetClaimsFromContext(ctx)
		if !ok || !claims.IsUserTenantRequest() {
			return nil, rbac_v1.ErrorUnauthorized("invalid token")
		}

		assignedRoles, err := uc.assignedRepo.ListRoles(ctx, data.ListRolesDto{
			TenantId:    claims.GetTenantId(),
			IdentityIDs: claims.GetIdentities(),
		})
		if err != nil {
			return nil, err
		}

		rolesIds := make([]int64, len(assignedRoles))
		for i, assignedRole := range assignedRoles {
			rolesIds[i] = assignedRole.RoleID
		}

		permissions, err := uc.roleRepo.ListRolesPermissions(ctx, data.FilterRolePermissions{
			TenantId:   claims.GetTenantId(),
			RolesIds:   rolesIds,
			DeniedOnly: true,
		})
		if err != nil {
			return nil, err
		}

		excludePermissions := make(map[string]bool)
		for _, permission := range permissions {
			excludePermissions[permission.PermissionID] = true
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
