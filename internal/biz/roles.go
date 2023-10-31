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
type RolesUsecase struct {
	log       *log.Helper
	discovery registry.Discovery
	roleRepo  data.RoleRepo
}

// NewRolesUsecase .
func NewRolesUsecase(logger log.Logger, c *data.Config, usersRepo data.RoleRepo) (*RolesUsecase, error) {
	return &RolesUsecase{
		log:       log.NewHelper(logger),
		discovery: c.GetRegistry(),
		roleRepo:  usersRepo,
	}, nil
}

func (uc *RolesUsecase) GetRoleById(ctx context.Context, roleId int64) (*ent.Role, error) {
	return uc.roleRepo.GetRoleById(ctx, roleId)
}

func (uc *RolesUsecase) UpdateRole(ctx context.Context, userId int64, data data.UpdateRoleDto) (*ent.Role, error) {
	return uc.roleRepo.UpdateRole(ctx, userId, data)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, roleId int64) error {
	return uc.roleRepo.DeleteRole(ctx, roleId)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, teamId int64, name string) ([]*ent.Role, error) {
	return uc.roleRepo.GetRolesList(ctx, teamId, name)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, createRoleDto data.CreateRoleDto) (*ent.Role, error) {
	return uc.roleRepo.CreateRole(ctx, createRoleDto)
}

func (uc *RolesUsecase) AddPermissionToRole(ctx context.Context, roleId int64, permissionId string, fields []string) (*ent.Permission, error) {
	return uc.roleRepo.AddPermissionToRole(ctx, roleId, permissionId, fields)
}

func (uc *RolesUsecase) RemovePermissionFromRole(ctx context.Context, roleId int64, permissionId string) error {
	return uc.roleRepo.RemovePermissionFromRole(ctx, roleId, permissionId)
}

func (uc *RolesUsecase) ListRolePermissions(ctx context.Context, roleId int64) ([]*ent.Permission, error) {
	return uc.roleRepo.ListRolePermissions(ctx, roleId)
}
