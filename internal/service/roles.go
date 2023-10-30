package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"rbac/internal/biz"
	"rbac/internal/data"

	permissions_v1 "rbac/api/permissions/v1"
	role_v1 "rbac/api/roles/v1"
)

type RolesService struct {
	role_v1.UnimplementedRolesServer

	log *log.Helper
	jwt *data.JwtProcessor
	uc  *biz.RolesUsecase
}

func NewRolesService(logger log.Logger, jwt *data.JwtProcessor, uc *biz.RolesUsecase) *RolesService {
	return &RolesService{
		log: log.NewHelper(logger),
		jwt: jwt,
		uc:  uc,
	}
}

func (s *RolesService) CreateRole(ctx context.Context, req *role_v1.CreateRoleRequest) (*role_v1.RoleReply, error) {
	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		Name:        req.Name,
		Description: req.Description,
		TeamId:      req.TeamId,
	})
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	return &role_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      role.TeamID,
	}, nil
}
func (s *RolesService) UpdateRole(ctx context.Context, req *role_v1.UpdateRoleRequest) (*role_v1.RoleReply, error) {
	role, err := s.uc.UpdateRole(ctx, req.RoleId, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}

	return &role_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      role.TeamID,
	}, nil
}
func (s *RolesService) DeleteRole(ctx context.Context, req *role_v1.DeleteRoleRequest) (*role_v1.EmptyReply, error) {
	err := s.uc.DeleteRole(ctx, req.RoleId)
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	return &role_v1.EmptyReply{}, nil
}
func (s *RolesService) GetRole(ctx context.Context, req *role_v1.GetRoleRequest) (*role_v1.RoleReply, error) {
	role, err := s.uc.GetRoleById(ctx, req.RoleId)
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	return &role_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      role.TeamID,
	}, nil
}
func (s *RolesService) ListRoles(ctx context.Context, req *role_v1.ListRolesRequest) (*role_v1.ListRolesReply, error) {
	roles, err := s.uc.GetRoles(ctx, req.TeamId, req.GetName(), int64(req.Paginate.Page), req.Paginate.GetPageSize())
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	replyRoles := make([]*role_v1.RoleReply, 0, len(roles))
	for _, role := range roles {
		replyRoles = append(replyRoles, &role_v1.RoleReply{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			TeamId:      role.TeamID,
		})
	}
	return &role_v1.ListRolesReply{
		Roles: replyRoles,
	}, nil
}
func (s *RolesService) AddPermissionToRole(ctx context.Context, req *role_v1.AddPermissionToRoleRequest) (*role_v1.AddPermissionToRoleReply, error) {
	role, err := s.uc.AddPermissionToRole(ctx, req.RoleId, req.PermissionId, req.Fields)
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	s.log.Debug("role", role)
	return &role_v1.AddPermissionToRoleReply{}, nil
}
func (s *RolesService) RemovePermissionFromRole(ctx context.Context, req *role_v1.RemovePermissionFromRoleRequest) (*role_v1.EmptyReply, error) {
	err := s.uc.RemovePermissionFromRole(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	return &role_v1.EmptyReply{}, nil
}
func (s *RolesService) ListRolePermissions(ctx context.Context, req *role_v1.RolesPermissionsRequest) (*role_v1.RolesPermissionsReply, error) {
	permissions, err := s.uc.ListRolePermissions(ctx, req.RoleId, int64(req.Paginate.GetPage()), req.Paginate.GetPageSize())
	if err != nil {
		return nil, role_v1.ErrorDatabaseQuery(err.Error())
	}
	s.log.Debug("roles", permissions)
	permissionsReply := make([]*permissions_v1.PermissionReply, 0, len(permissions))
	for _, permission := range permissions {
		permissionsReply = append(permissionsReply, &permissions_v1.PermissionReply{
			Id:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			Fields:      permission.Fields,
		})
	}

	return &role_v1.RolesPermissionsReply{
		Permissions: permissionsReply,
	}, nil
}
