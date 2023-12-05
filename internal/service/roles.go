package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/utils/v1/jwt"
)

type RolesService struct {
	v1.UnimplementedRolesServer

	jwt *jwt.JwtProcessor
	uc  *biz.RolesUsecase
	pu  *biz.PermissionsUsecase
}

func NewRolesService(jwt *jwt.JwtProcessor, uc *biz.RolesUsecase, pu *biz.PermissionsUsecase) *RolesService {
	return &RolesService{
		jwt: jwt,
		uc:  uc,
		pu:  pu,
	}
}

func (s *RolesService) roleReply(role *ent.Role) *v1.RoleReply {
	return &v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}
}

func (s *RolesService) rolePermissionReply(rolePermission *ent.RolePermission) *v1.RolePermissionReply {
	return &v1.RolePermissionReply{
		Role: s.roleReply(rolePermission.Edges.Role),
		Permission: &v1.PermissionReply{
			Id:    rolePermission.Edges.Permission.ID,
			AppId: rolePermission.Edges.Permission.AppID,
		},
		Fields: rolePermission.Fields,
		Deny:   &rolePermission.Deny,
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can create a role

	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		TenantId:    claims.GetTenantId(),
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return s.roleReply(role), nil
}

func (s *RolesService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.RoleReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can update role

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	updated, err := s.uc.UpdateRole(ctx, role, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return s.roleReply(updated), nil
}

func (s *RolesService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.EmptyReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can delete role

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	err = s.uc.DeleteRole(ctx, role)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}

func (s *RolesService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.RoleReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can view role details

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}
	return s.roleReply(role), nil
}

func (s *RolesService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can view role list

	roles, err := s.uc.GetRoles(ctx, claims.GetTenantId(), req.Search)
	if err != nil {
		return nil, err
	}

	result := make([]*v1.RoleReply, len(roles))
	for i, role := range roles {
		result[i] = s.roleReply(role)
	}
	return &v1.ListRolesReply{
		Roles: result,
	}, nil
}

func (s *RolesService) AddPermissionToRole(ctx context.Context, req *v1.AddPermissionToRoleRequest) (*v1.RolePermissionReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can add permission to role

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	permission, err := s.pu.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, err
	}

	rolePermission, err := s.uc.AddPermissionToRole(ctx, role, permission, data.CreateRolePermissionDto{
		Fields: req.Fields,
		Deny:   *req.Deny,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return s.rolePermissionReply(rolePermission), nil
}

func (s *RolesService) RemovePermissionFromRole(ctx context.Context, req *v1.RemovePermissionFromRoleRequest) (*v1.EmptyReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can remove permission from role

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
	return &v1.EmptyReply{}, nil
}

func (s *RolesService) ListRolePermissions(ctx context.Context, req *v1.RolesPermissionsRequest) (*v1.RolesPermissionsReply, error) {
	claims, ok := s.jwt.GetClaimsFromContext(ctx)
	if !ok || !claims.IsUserTenantRequest() {
		return nil, v1.ErrorUnauthorized("invalid token")
	}
	// todo checkPermissions can remove permission from role

	role, err := s.uc.GetRoleById(ctx, claims.GetTenantId(), req.RoleId)
	if err != nil {
		return nil, err
	}

	rolePermissions, err := s.uc.ListRolePermissions(ctx, role)
	if err != nil {
		return nil, err
	}

	permissions := make([]*v1.RolePermissionReply, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = s.rolePermissionReply(rp)
	}
	return &v1.RolesPermissionsReply{
		Permissions: permissions,
	}, nil
}
