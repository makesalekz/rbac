package biz

import (
	"context"
	v1 "rbac/api/rbac/v1"

	"rbac/ent"
	"rbac/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

// RolesUsecase .
type RolesUsecase struct {
	log      *log.Helper
	jwt      *data.JwtProcessor
	roleRepo data.RoleRepo
}

// NewRolesUsecase .
func NewRolesUsecase(logger log.Logger, jwt *data.JwtProcessor, usersRepo data.RoleRepo) (*RolesUsecase, error) {
	return &RolesUsecase{
		log:      log.NewHelper(logger),
		jwt:      jwt,
		roleRepo: usersRepo,
	}, nil
}

func (uc *RolesUsecase) GetRoleById(ctx context.Context, roleId int64) (*ent.Role, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("GetRoleById", "userId", userId, "tenant", tenant)
	return uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
}

func (uc *RolesUsecase) UpdateRole(ctx context.Context, roleId int64, data data.UpdateRoleDto) (*ent.Role, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("UpdateRole", "userId", userId, "tenant", tenant, "data", data)
	entry, err := uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
	if err != nil {
		return nil, v1.ErrorNotFound("Role not found")
	}
	if entry.IsSystem {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	return uc.roleRepo.UpdateRole(ctx, roleId, data)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, roleId int64) error {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	uc.log.Debug("DeleteRole", "userId", userId, "tenant", tenant)
	entry, err := uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
	if err != nil {
		return v1.ErrorNotFound("Role not found")
	}
	if entry.IsSystem {
		return v1.ErrorForbidden("Forbidden")
	}
	return uc.roleRepo.DeleteRole(ctx, roleId)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, tenantId int64, name string) ([]*ent.Role, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	uc.log.Debug("GetRoles", "userId", userId, "tenant", tenant)
	return uc.roleRepo.GetRolesList(ctx, tenantId, name)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, createRoleDto data.CreateRoleDto) (*ent.Role, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != createRoleDto.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	uc.log.Debug("CreateRole", "userId", userId, "tenant", tenant, "data", createRoleDto)
	return uc.roleRepo.CreateRole(ctx, createRoleDto)
}

func (uc *RolesUsecase) AddPermissionToRole(ctx context.Context, createDto data.CreateRolePermissionDto) (*ent.RolePermission, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != createDto.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	uc.log.Debug("AddPermissionToRole", "userId", userId, "tenant", tenant, "data", createDto)
	return uc.roleRepo.AddPermissionToRole(ctx, createDto)
}

func (uc *RolesUsecase) RemovePermissionFromRole(ctx context.Context, roleId, tenantId int64, permissionId string) error {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return v1.ErrorForbidden("Forbidden")
	}
	uc.log.Debug("RemovePermissionFromRole", "userId", userId, "tenant", tenant, "roleId", roleId, "permissionId", permissionId)
	return uc.roleRepo.RemovePermissionFromRole(ctx, roleId, tenantId, permissionId)
}

func (uc *RolesUsecase) ListRolePermissions(ctx context.Context, roleId, tenantId int64) ([]*ent.RolePermission, error) {
	userId, tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	uc.log.Debug("ListRolePermissions", "userId", userId, "tenant", tenant, "roleId", roleId)
	return uc.roleRepo.ListRolePermissions(ctx, roleId, tenantId)
}
