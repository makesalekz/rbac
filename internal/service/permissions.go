package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"rbac/api/rbac/permissions/v1"
	"rbac/internal/biz"
	"rbac/internal/data"
)

type PermissionsService struct {
	permissions_v1.PermissionsServer
	log *log.Helper
	jwt *data.JwtProcessor
	uc  *biz.PermissionsUsecase
}

func NewPermissionsService(logger log.Logger, jwt *data.JwtProcessor, uc *biz.PermissionsUsecase) *PermissionsService {
	return &PermissionsService{
		log: log.NewHelper(logger),
		jwt: jwt,
		uc:  uc,
	}
}

func (s *PermissionsService) CreatePermission(ctx context.Context, req *permissions_v1.CreatePermissionRequest) (*permissions_v1.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.CreatePermission request")
	permission, err := s.uc.CreatePermission(ctx, data.CreatePermissionDto{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		AppId:       req.AppId,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, permissions_v1.ErrorDatabaseQuery("Internal error")
	}
	return &permissions_v1.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) UpdatePermission(ctx context.Context, req *permissions_v1.UpdatePermissionRequest) (*permissions_v1.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.UpdatePermission request")
	permission, err := s.uc.UpdatePermission(ctx, req.PermissionId, data.UpdatePermissionDto{
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, permissions_v1.ErrorDatabaseQuery("Internal error")
	}
	return &permissions_v1.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) DeletePermission(ctx context.Context, req *permissions_v1.DeletePermissionRequest) (*permissions_v1.EmptyReply, error) {
	s.log.Debug("Received PermissionsService.DeletePermission request")
	err := s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, permissions_v1.ErrorDatabaseQuery("Internal error")
	}
	return &permissions_v1.EmptyReply{}, nil
}
func (s *PermissionsService) GetPermission(ctx context.Context, req *permissions_v1.GetPermissionRequest) (*permissions_v1.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.GetPermission request")
	permission, err := s.uc.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, permissions_v1.ErrorDatabaseQuery("Internal error")
	}
	return &permissions_v1.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) ListPermissions(ctx context.Context, req *permissions_v1.ListPermissionsRequest) (*permissions_v1.ListPermissionsReply, error) {
	s.log.Debug("Received PermissionsService.ListPermissions request")
	permissions, err := s.uc.GetPermissions(ctx, req.AppId, req.Ids)
	if err != nil {
		return nil, permissions_v1.ErrorDatabaseQuery("Internal error")
	}
	var replyPermissions []*permissions_v1.PermissionReply
	for _, permission := range permissions {
		replyPermissions = append(replyPermissions, &permissions_v1.PermissionReply{
			Id:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			AppId:       permission.AppID,
			Fields:      permission.Fields,
		})
	}
	return &permissions_v1.ListPermissionsReply{Permissions: replyPermissions}, nil
}
