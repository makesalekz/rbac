package biz

import (
	"context"
	_ "embed"

	"rbac/ent"
	"rbac/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
)

// PermissionsUsecase .
type PermissionsUsecase struct {
	log            *log.Helper
	discovery      registry.Discovery
	permissionRepo data.PermissionRepo
}

// NewPermissionUsecase .
func NewPermissionUsecase(logger log.Logger, c *data.Config, permissionRepo data.PermissionRepo) (*PermissionsUsecase, error) {
	return &PermissionsUsecase{
		log:            log.NewHelper(logger),
		discovery:      c.GetRegistry(),
		permissionRepo: permissionRepo,
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
