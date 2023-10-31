package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"rbac/api/rbac/permissions/v1"
	"rbac/api/rbac/roles/v1"
	"rbac/internal/biz"
	"rbac/internal/data"
)

type RolesService struct {
	roles_v1.UnimplementedRolesServer

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

func (s *RolesService) CreateRole(ctx context.Context, req *roles_v1.CreateRoleRequest) (*roles_v1.RoleReply, error) {
	role, err := s.uc.CreateRole(ctx, data.CreateRoleDto{
		Name:        req.Name,
		Description: req.Description,
		TeamId:      req.TeamId,
	})
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	return &roles_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      *role.TeamID,
	}, nil
}
func (s *RolesService) UpdateRole(ctx context.Context, req *roles_v1.UpdateRoleRequest) (*roles_v1.RoleReply, error) {
	role, err := s.uc.UpdateRole(ctx, req.RoleId, data.UpdateRoleDto{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}

	return &roles_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      role.TeamID,
	}, nil
}
func (s *RolesService) DeleteRole(ctx context.Context, req *roles_v1.DeleteRoleRequest) (*roles_v1.EmptyReply, error) {
	err := s.uc.DeleteRole(ctx, req.RoleId)
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	return &roles_v1.EmptyReply{}, nil
}
func (s *RolesService) GetRole(ctx context.Context, req *roles_v1.GetRoleRequest) (*roles_v1.RoleReply, error) {
	role, err := s.uc.GetRoleById(ctx, req.RoleId)
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	return &roles_v1.RoleReply{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		TeamId:      *role.TeamID,
		CreatedAt:   role.CreatedAt.String(),
		UpdatedAt:   role.UpdatedAt.String(),
		DeletedAt:   role.DeletedAt.String(),
	}, nil
}
func (s *RolesService) ListRoles(ctx context.Context, req *roles_v1.ListRolesRequest) (*roles_v1.ListRolesReply, error) {
	roles, err := s.uc.GetRoles(ctx, req.TeamId, req.GetName())
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	replyRoles := make([]*roles_v1.RoleReply, 0, len(roles))
	for _, role := range roles {
		replyRoles = append(replyRoles, &roles_v1.RoleReply{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			TeamId:      *role.TeamID,
		})
	}
	return &roles_v1.ListRolesReply{
		Roles: replyRoles,
	}, nil
}
func (s *RolesService) AddPermissionToRole(ctx context.Context, req *roles_v1.AddPermissionToRoleRequest) (*roles_v1.AddPermissionToRoleReply, error) {
	role, err := s.uc.AddPermissionToRole(ctx, req.RoleId, req.PermissionId, req.Fields)
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	s.log.Debug("role", role)
	return &roles_v1.AddPermissionToRoleReply{}, nil
}
func (s *RolesService) RemovePermissionFromRole(ctx context.Context, req *roles_v1.RemovePermissionFromRoleRequest) (*roles_v1.EmptyReply, error) {
	err := s.uc.RemovePermissionFromRole(ctx, req.RoleId, req.PermissionId)
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
	}
	return &roles_v1.EmptyReply{}, nil
}
func (s *RolesService) ListRolePermissions(ctx context.Context, req *roles_v1.RolesPermissionsRequest) (*roles_v1.RolesPermissionsReply, error) {
	permissions, err := s.uc.ListRolePermissions(ctx, req.RoleId)
	if err != nil {
		return nil, roles_v1.ErrorDatabaseQuery(err.Error())
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

	return &roles_v1.RolesPermissionsReply{
		Permissions: permissionsReply,
	}, nil
}
