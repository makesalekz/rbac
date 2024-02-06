package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
)

type RolesService struct {
	v1.UnimplementedRolesServer

	sh *ServiceHelper
	uc *biz.RolesUsecase
	pu *biz.PermissionsUsecase
	au *biz.TeamIdentityUsecase
}

func NewRolesService(
	sh *ServiceHelper,
	uc *biz.RolesUsecase,
	pu *biz.PermissionsUsecase,
	au *biz.TeamIdentityUsecase,
) *RolesService {
	return &RolesService{
		sh: sh,
		uc: uc,
		pu: pu,
		au: au,
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleReply, error) {
	permissionToCheck := "admin.role.create"
	if req.IsSystem {
		permissionToCheck = "admin.role_system.create"
	}

	claims, _, err := s.sh.HasPermission(ctx, permissionToCheck)
	if err != nil {
		return nil, err
	}

	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		TenantId:    claims.GetTenantId(),
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    req.IsSystem,
		Allow:       req.Allow,
		Deny:        req.Deny,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return &v1.RoleReply{
		Role: s.roleReply(role),
	}, nil
}

func (s *RolesService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.RoleReply, error) {
	claims, err := s.sh.GetClaims(ctx)
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	permissionToCheck := "admin.role.update"
	if role.IsSystem {
		permissionToCheck = "admin.role_system.update"
	}

	_, _, err = s.sh.HasPermission(ctx, permissionToCheck)
	if err != nil {
		return nil, err
	}

	updated, err := s.uc.UpdateRole(ctx, role, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
		Allow:       req.Allow,
		Deny:        req.Deny,
	})
	if err != nil {
		return nil, err
	}

	return &v1.RoleReply{
		Role: s.roleReply(updated),
	}, nil
}

func (s *RolesService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*utils_v1.EmptyReply, error) {
	claims, err := s.sh.GetClaims(ctx)
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	permissionToCheck := "admin.role.delete"
	if role.IsSystem {
		permissionToCheck = "admin.role_system.delete"
	}

	_, _, err = s.sh.HasPermission(ctx, permissionToCheck)
	if err != nil {
		return nil, err
	}

	err = s.uc.DeleteRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.RoleReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.role.read")
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}
	return &v1.RoleReply{
		Role: s.roleReply(role),
	}, nil
}

func (s *RolesService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.role.read")
	if err != nil {
		return nil, err
	}

	roles, err := s.uc.GetRoles(ctx, claims.GetTenantId(), req.Search)
	if err != nil {
		return nil, err
	}

	result := make([]*v1.Role, len(roles))
	for i, role := range roles {
		result[i] = s.roleReply(role)
	}
	return &v1.ListRolesReply{
		Roles: result,
	}, nil
}

func (s *RolesService) SetRolePermission(ctx context.Context, req *v1.SetRolePermissionRequest) (*utils_v1.EmptyReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.role.update")
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	permission, err := s.pu.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}

	err = s.uc.SetRolePermission(ctx, role, permission, data.CreateRolePermissionDto{
		Fields: req.Fields,
		Deny:   *req.Deny,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) DeleteRolePermission(ctx context.Context, req *v1.DeleteRolePermissionRequest) (*utils_v1.EmptyReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.role.update")
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	permission, err := s.pu.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}

	err = s.uc.RemovePermissionFromRole(ctx, role, permission)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) ListRolePermissions(ctx context.Context, req *v1.ListRolePermissionsRequest) (*v1.RolePermissionsReply, error) {
	claims, _, err := s.sh.HasPermission(ctx, "admin.role.read")
	if err != nil {
		return nil, err
	}

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	rolePermissions, err := s.uc.ListRolePermissions(ctx, role)
	if err != nil {
		return nil, err
	}

	permissions := make([]*v1.RolePermission, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = s.rolePermissionReply(rp)
	}
	return &v1.RolePermissionsReply{
		Permissions: permissions,
	}, nil
}

func (s *RolesService) roleReply(role *ent.Role) *v1.Role {
	return &v1.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}
}

func (s *RolesService) rolePermissionReply(rolePermission *ent.RolePermission) *v1.RolePermission {
	return &v1.RolePermission{
		Id:     rolePermission.PermissionID,
		Fields: rolePermission.Fields,
		Deny:   rolePermission.Deny,
	}
}
