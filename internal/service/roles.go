package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

type RolesService struct {
	v1.UnimplementedRolesServer

	uc *biz.RolesUsecase
}

func NewRolesService(uc *biz.RolesUsecase) *RolesService {
	return &RolesService{
		uc: uc,
	}
}

func (s *RolesService) roleReply(role ent.Role) *v1.RoleReply {
	return &v1.RoleReply{
		Id:          role.ID,
		TenantId:    *role.TenantID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) (*v1.RoleReply, error) {
	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		TenantId:    req.TenantId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) UpdateRole(ctx context.Context, req *v1.UpdateRoleRequest) (*v1.RoleReply, error) {
	role, err := s.uc.UpdateRole(ctx, req.RoleId, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) DeleteRole(ctx context.Context, req *v1.DeleteRoleRequest) (*v1.EmptyReply, error) {
	err := s.uc.DeleteRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}
func (s *RolesService) GetRole(ctx context.Context, req *v1.GetRoleRequest) (*v1.RoleReply, error) {
	role, err := s.uc.GetRoleById(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) ListRoles(ctx context.Context, req *v1.ListRolesRequest) (*v1.ListRolesReply, error) {
	roles, err := s.uc.GetRoles(ctx, req.TenantId, *req.Name)
	if err != nil {
		return nil, err
	}
	result := make([]*v1.RoleReply, len(roles))
	for i, role := range roles {
		result[i] = s.roleReply(*role)
	}
	return &v1.ListRolesReply{
		Roles: result,
	}, nil
}
func (s *RolesService) AddPermissionToRole(ctx context.Context, req *v1.AddPermissionToRoleRequest) (*v1.RolePermissionReply, error) {
	rolePermission, err := s.uc.AddPermissionToRole(ctx, data.CreateRolePermissionDto{
		RoleId:       req.RoleId,
		PermissionId: req.PermissionId,
		TenantId:     req.TenantId,
		Fields:       req.Fields,
		Deny:         *req.Deny,
	})
	if err != nil {
		return nil, v1.ErrorDatabaseQuery(err.Error())
	}
	return &v1.RolePermissionReply{
		Role: s.roleReply(*rolePermission.Edges.Role),
		Permission: &v1.PermissionReply{
			Id:    rolePermission.Edges.Permission.ID,
			AppId: rolePermission.Edges.Permission.AppID,
		},
		TenantId: *rolePermission.TenantID,
		Fields:   rolePermission.Fields,
		Deny:     &rolePermission.Deny,
	}, nil
}
func (s *RolesService) RemovePermissionFromRole(ctx context.Context, req *v1.RemovePermissionFromRoleRequest) (*v1.EmptyReply, error) {
	err := s.uc.RemovePermissionFromRole(ctx, req.RoleId, req.TenantId, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}
func (s *RolesService) ListRolePermissions(ctx context.Context, req *v1.RolesPermissionsRequest) (*v1.RolesPermissionsReply, error) {
	rolePermissions, err := s.uc.ListRolePermissions(ctx, req.RoleId, req.TenantId)
	if err != nil {
		return nil, err
	}
	permissions := make([]*v1.RolePermissionReply, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = &v1.RolePermissionReply{
			Role: s.roleReply(*rp.Edges.Role),
			Permission: &v1.PermissionReply{
				Id:    rp.Edges.Permission.ID,
				AppId: rp.Edges.Permission.AppID,
			},
			TenantId: *rp.TenantID,
			Fields:   rp.Fields,
			Deny:     &rp.Deny,
		}
	}
	return &v1.RolesPermissionsReply{
		Permissions: permissions,
	}, nil
}
