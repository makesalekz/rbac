package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
	"gitlab.calendaria.team/services/utils/v2/auth"
)

type RolesService struct {
	v1.UnimplementedRolesServer

	sh    *ServiceHelper
	uc    *biz.RolesUsecase
	pu    *biz.PermissionsUsecase
	au    *biz.AssignedRolesUsecase
	check *biz.CheckPermissionsUsecase
}

func NewRolesService(
	sh *ServiceHelper,
	uc *biz.RolesUsecase,
	pu *biz.PermissionsUsecase,
	au *biz.AssignedRolesUsecase,
	check *biz.CheckPermissionsUsecase,
) *RolesService {
	return &RolesService{
		sh:    sh,
		uc:    uc,
		pu:    pu,
		au:    au,
		check: check,
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	dto := data.CreateRoleDto{
		TenantID:    tenantID,
		Name:        req.GetName(),
		Description: req.GetDescription(),
		IsSystem:    req.GetIsSystem(),
		Allow:       req.GetAllow(),
		Deny:        req.GetDeny(),
	}
	err := dto.Validate()
	if err != nil {
		return nil, err
	}

	role, err := s.uc.CreateRole(ctx, dto)
	if err != nil {
		return nil, v1.ErrorDatabaseQuery("failed to create role, error: %s", err.Error())
	}
	return &v1.RoleReply{
		Role: s.roleReply(role),
	}, nil
}

func (s *RolesService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.RoleReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetRoleId() == 0 {
		return nil, v1.ErrorBadRequest("empty role id")
	}

	updated, err := s.uc.UpdateRole(ctx, tenantID, req.GetRoleId(), data.UpdateRoleDto{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Allow:       req.GetAllow(),
		Deny:        req.GetDeny(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.RoleReply{
		Role: s.roleReply(updated),
	}, nil
}

func (s *RolesService) DeleteRole(ctx context.Context, req *v1.RoleRequest) (*utils_v1.EmptyReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetRoleId() == 0 {
		return nil, v1.ErrorBadRequest("empty role id")
	}

	err := s.uc.DeleteRole(ctx, tenantID, req.GetRoleId())
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) GetRole(ctx context.Context, req *v1.RoleRequest) (*v1.RoleReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetRoleId() == 0 {
		return nil, v1.ErrorBadRequest("empty role id")
	}

	role, err := s.uc.GetRoleByID(ctx, tenantID, req.GetRoleId())
	if err != nil {
		return nil, err
	}
	return &v1.RoleReply{
		Role: s.roleReply(role),
	}, nil
}

func (s *RolesService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	roles, err := s.uc.GetRoles(ctx, tenantID, req.GetSearch())
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

func (s *RolesService) AddPermissionToRole(
	ctx context.Context,
	req *v1.AddPermissionToRoleRequest,
) (*utils_v1.EmptyReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetRoleId() == 0 {
		return nil, v1.ErrorBadRequest("empty role id")
	}

	if req.GetPermissionId() == "" {
		return nil, v1.ErrorBadRequest("empty permission id")
	}

	permission, err := s.pu.GetPermissionByID(ctx, req.GetPermissionId())
	if err != nil {
		return nil, err
	}

	err = s.uc.SetRolePermission(ctx, tenantID, req.GetRoleId(), permission, data.CreateRolePermissionDto{
		Fields: req.GetFields(),
		Deny:   req.GetDeny(),
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery("failed to add permission to role, error: %s", err.Error())
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) RemovePermissionFromRole(
	ctx context.Context,
	req *v1.RemovePermissionFromRoleRequest,
) (*utils_v1.EmptyReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	if req.GetRoleId() == 0 {
		return nil, v1.ErrorBadRequest("empty role id")
	}

	if req.GetPermissionId() == "" {
		return nil, v1.ErrorBadRequest("empty permission id")
	}

	permission, err := s.pu.GetPermissionByID(ctx, req.GetPermissionId())
	if err != nil {
		return nil, err
	}

	err = s.uc.RemovePermissionFromRole(ctx, tenantID, req.GetRoleId(), permission)
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *RolesService) ListRolePermissions(ctx context.Context, req *v1.RoleRequest) (*v1.RolePermissionsReply, error) {
	tenantID := auth.GetTenantIdFromContext(ctx)
	if tenantID == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	appID := auth.GetAppIdFromContext(ctx)
	if appID == "" {
		return nil, v1.ErrorEmptyAppId("empty app id")
	}

	rolePermissions, err := s.uc.ListRolePermissions(ctx, tenantID, appID, req.GetRoleId())
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
