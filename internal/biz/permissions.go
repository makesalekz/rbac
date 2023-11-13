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
	log            *log.Helper
	jwt            *data.JwtProcessor
	permissionRepo data.PermissionRepo
}

// NewPermissionUsecase .
func NewPermissionUsecase(logger log.Logger, jwt *data.JwtProcessor, permissionRepo data.PermissionRepo) (*PermissionsUsecase, error) {
	return &PermissionsUsecase{
		log:            log.NewHelper(logger),
		jwt:            jwt,
		permissionRepo: permissionRepo,
	}, nil
}

func (uc *PermissionsUsecase) GetPermissionById(ctx context.Context, permissionId string) (*ent.Permission, error) {
	_, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("GetPermissionById", "permissionId", permissionId)
	return uc.permissionRepo.GetPermissionById(ctx, permissionId)
}

func (uc *PermissionsUsecase) CreatePermission(ctx context.Context, data data.CreatePermissionDto) (*ent.Permission, error) {
	_, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("CreatePermission", "data", data)
	return uc.permissionRepo.CreatePermission(ctx, data)
}

func (uc *PermissionsUsecase) UpdatePermission(ctx context.Context, permissionId string, data data.UpdatePermissionDto) (*ent.Permission, error) {
	_, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("UpdatePermission", "permissionId", permissionId, "data", data)
	return uc.permissionRepo.UpdatePermission(ctx, permissionId, data)
}

func (uc *PermissionsUsecase) DeletePermission(ctx context.Context, permissionId string) error {
	_, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("DeletePermission", "permissionId", permissionId)

	_, err := uc.GetPermissionById(ctx, permissionId)
	if err != nil {
		return v1.ErrorNotFound("Permission with given ID not found")
	}
	return uc.permissionRepo.DeletePermission(ctx, permissionId)
}

func (uc *PermissionsUsecase) GetPermissions(ctx context.Context, appId string, permissionIds []string) ([]*ent.Permission, error) {
	_, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("GetPermissions", "appId", appId, "permissionIds", permissionIds)
	return uc.permissionRepo.GetPermissions(ctx, appId, permissionIds)
}
