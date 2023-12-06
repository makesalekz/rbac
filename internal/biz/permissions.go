package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

// PermissionsUsecase .
type PermissionsUsecase struct {
	log            *log.Helper
	jwt            *jwt.JwtProcessor
	permissionRepo data.PermissionRepo
}

// NewPermissionUsecase .
func NewPermissionUsecase(logger log.Logger, jwt *jwt.JwtProcessor, permissionRepo data.PermissionRepo) (*PermissionsUsecase, error) {
	return &PermissionsUsecase{
		log:            log.NewHelper(logger),
		jwt:            jwt,
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
