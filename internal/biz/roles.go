package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
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
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	// todo checkPermissions can view role details
	uc.log.Debug("GetRoleById", "tenant", tenant)
	return uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
}

func (uc *RolesUsecase) UpdateRole(ctx context.Context, roleId int64, data data.UpdateRoleDto) (*ent.Role, error) {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	// todo checkPermissions can update role
	uc.log.Debug("UpdateRole", "tenant", tenant, "data", data)
	entry, err := uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, v1.ErrorNotFound("Role not found")
		}
		return nil, v1.ErrorDatabaseQuery("Internal error")
	}
	if entry.IsSystem {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	return uc.roleRepo.UpdateRole(ctx, roleId, data)
}

func (uc *RolesUsecase) DeleteRole(ctx context.Context, roleId int64) error {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	// todo checkPermissions can delete role
	uc.log.Debug("DeleteRole", "tenant", tenant)
	entry, err := uc.roleRepo.GetRoleById(ctx, roleId, tenant.TenantId)
	if err != nil {
		if ent.IsNotFound(err) {
			return v1.ErrorNotFound("Role not found")
		}
		return v1.ErrorDatabaseQuery("Internal error")
	}
	if entry.IsSystem {
		return v1.ErrorForbidden("Forbidden")
	}
	return uc.roleRepo.DeleteRole(ctx, roleId)
}

func (uc *RolesUsecase) GetRoles(ctx context.Context, tenantId int64, name string) ([]*ent.Role, error) {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions can view role list
	uc.log.Debug("GetRoles", "tenant", tenant)
	return uc.roleRepo.GetRolesList(ctx, tenantId, name)
}

func (uc *RolesUsecase) CreateRole(ctx context.Context, createRoleDto data.CreateRoleDto) (*ent.Role, error) {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != createRoleDto.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions can create a role
	uc.log.Debug("CreateRole", "tenant", tenant, "data", createRoleDto)
	return uc.roleRepo.CreateRole(ctx, createRoleDto)
}

func (uc *RolesUsecase) AddPermissionToRole(ctx context.Context, createDto data.CreateRolePermissionDto) (*ent.RolePermission, error) {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != createDto.TenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions can add permission to role
	uc.log.Debug("AddPermissionToRole", "tenant", tenant, "data", createDto)
	return uc.roleRepo.AddPermissionToRole(ctx, createDto)
}

func (uc *RolesUsecase) RemovePermissionFromRole(ctx context.Context, roleId, tenantId int64, permissionId string) error {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions can remove permission from role
	uc.log.Debug("RemovePermissionFromRole", "tenant", tenant, "roleId", roleId, "permissionId", permissionId)
	return uc.roleRepo.RemovePermissionFromRole(ctx, roleId, tenantId, permissionId)
}

func (uc *RolesUsecase) ListRolePermissions(ctx context.Context, roleId, tenantId int64) ([]*ent.RolePermission, error) {
	tenant, ok := uc.jwt.GetTenantClaimsFromContext(ctx)
	if !ok {
		return nil, v1.ErrorUnauthorized("Unauthorized")
	}
	if tenant.TenantId != tenantId {
		return nil, v1.ErrorForbidden("Forbidden")
	}
	// todo checkPermissions can view role permissions
	uc.log.Debug("ListRolePermissions", "tenant", tenant, "roleId", roleId)
	return uc.roleRepo.ListRolePermissions(ctx, roleId, tenantId)
}
