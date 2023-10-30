package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"rbac/internal/biz"
	"rbac/internal/data"

	pb "rbac/api/permissions/v1"
)

type PermissionsService struct {
	pb.PermissionsServer
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

func (s *PermissionsService) CreatePermission(ctx context.Context, req *pb.CreatePermissionRequest) (*pb.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.CreatePermission request")
	permission, err := s.uc.CreatePermission(ctx, data.CreatePermissionDto{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		AppId:       req.AppId,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, pb.ErrorDatabaseQuery("Internal error")
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) UpdatePermission(ctx context.Context, req *pb.UpdatePermissionRequest) (*pb.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.UpdatePermission request")
	permission, err := s.uc.UpdatePermission(ctx, req.PermissionId, data.UpdatePermissionDto{
		Name:        req.Name,
		Description: req.Description,
		Fields:      req.Fields,
	})
	if err != nil {
		return nil, pb.ErrorDatabaseQuery("Internal error")
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) DeletePermission(ctx context.Context, req *pb.DeletePermissionRequest) (*pb.EmptyReply, error) {
	s.log.Debug("Received PermissionsService.DeletePermission request")
	err := s.uc.DeletePermission(ctx, req.PermissionId)
	if err != nil {
		return nil, pb.ErrorDatabaseQuery("Internal error")
	}
	return &pb.EmptyReply{}, nil
}
func (s *PermissionsService) GetPermission(ctx context.Context, req *pb.GetPermissionRequest) (*pb.PermissionReply, error) {
	s.log.Debug("Received PermissionsService.GetPermission request")
	permission, err := s.uc.GetPermissionById(ctx, req.PermissionId)
	if err != nil {
		return nil, pb.ErrorDatabaseQuery("Internal error")
	}
	return &pb.PermissionReply{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
		AppId:       permission.AppID,
		Fields:      permission.Fields,
	}, nil
}
func (s *PermissionsService) ListPermissions(ctx context.Context, req *pb.ListPermissionsRequest) (*pb.ListPermissionsReply, error) {
	s.log.Debug("Received PermissionsService.ListPermissions request")
	permissions, err := s.uc.GetPermissions(ctx, req.AppId, req.Ids)
	if err != nil {
		return nil, pb.ErrorDatabaseQuery("Internal error")
	}
	var replyPermissions []*pb.PermissionReply
	for _, permission := range permissions {
		replyPermissions = append(replyPermissions, &pb.PermissionReply{
			Id:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
			AppId:       permission.AppID,
			Fields:      permission.Fields,
		})
	}
	return &pb.ListPermissionsReply{Permissions: replyPermissions}, nil
}
