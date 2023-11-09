package service

import (
	"context"
	"rbac/ent"
	"rbac/internal/biz"
	"rbac/internal/data"

	pb "rbac/api/rbac/v1"
)

type RolesService struct {
	pb.UnimplementedRolesServer

	uc *biz.RolesUsecase
}

func NewRolesService(uc *biz.RolesUsecase) *RolesService {
	return &RolesService{
		uc: uc,
	}
}

func (s *RolesService) roleReply(role ent.Role) *pb.RoleReply {
	return &pb.RoleReply{
		Id:          role.ID,
		TenantId:    *role.TenantID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *pb.CreateRoleRequest) (*pb.RoleReply, error) {
	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		TenantId:    req.TenantId,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, pb.ErrorDatabaseQuery(err.Error())
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) UpdateRole(ctx context.Context, req *pb.UpdateRoleRequest) (*pb.RoleReply, error) {
	role, err := s.uc.UpdateRole(ctx, req.RoleId, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) DeleteRole(ctx context.Context, req *pb.DeleteRoleRequest) (*pb.EmptyReply, error) {
	err := s.uc.DeleteRole(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *RolesService) GetRole(ctx context.Context, req *pb.GetRoleRequest) (*pb.RoleReply, error) {
	role, err := s.uc.GetRoleById(ctx, req.RoleId)
	if err != nil {
		return nil, err
	}
	return s.roleReply(*role), nil
}
func (s *RolesService) ListRoles(ctx context.Context, req *pb.ListRolesRequest) (*pb.ListRolesReply, error) {
	roles, err := s.uc.GetRoles(ctx, req.TenantId, *req.Name)
	if err != nil {
		return nil, err
	}
	result := make([]*pb.RoleReply, len(roles))
	for i, role := range roles {
		result[i] = s.roleReply(*role)
	}
	return &pb.ListRolesReply{
		Roles: result,
	}, nil
}
func (s *RolesService) AddPermissionToRole(ctx context.Context, req *pb.AddPermissionToRoleRequest) (*pb.RolePermissionReply, error) {
	rolePermission, err := s.uc.AddPermissionToRole(ctx, data.CreateRolePermissionDto{
		RoleId:       req.RoleId,
		PermissionId: req.PermissionId,
		TenantId:     req.TenantId,
		Fields:       req.Fields,
		Deny:         *req.Deny,
	})
	if err != nil {
		return nil, pb.ErrorDatabaseQuery(err.Error())
	}
	return &pb.RolePermissionReply{
		Role: s.roleReply(*rolePermission.Edges.Role),
		Permission: &pb.PermissionReply{
			Id:    rolePermission.Edges.Permission.ID,
			AppId: rolePermission.Edges.Permission.AppID,
		},
		TenantId: *rolePermission.TenantID,
		Fields:   rolePermission.Fields,
		Deny:     &rolePermission.Deny,
	}, nil
}
func (s *RolesService) RemovePermissionFromRole(ctx context.Context, req *pb.RemovePermissionFromRoleRequest) (*pb.EmptyReply, error) {
	err := s.uc.RemovePermissionFromRole(ctx, req.RoleId, req.TenantId, req.PermissionId)
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *RolesService) ListRolePermissions(ctx context.Context, req *pb.RolesPermissionsRequest) (*pb.RolesPermissionsReply, error) {
	rolePermissions, err := s.uc.ListRolePermissions(ctx, req.RoleId, req.TenantId)
	if err != nil {
		return nil, err
	}
	permissions := make([]*pb.RolePermissionReply, len(rolePermissions))
	for i, rp := range rolePermissions {
		permissions[i] = &pb.RolePermissionReply{
			Role: s.roleReply(*rp.Edges.Role),
			Permission: &pb.PermissionReply{
				Id:    rp.Edges.Permission.ID,
				AppId: rp.Edges.Permission.AppID,
			},
			TenantId: *rp.TenantID,
			Fields:   rp.Fields,
			Deny:     &rp.Deny,
		}
	}
	return &pb.RolesPermissionsReply{
		Permissions: permissions,
	}, nil
}
