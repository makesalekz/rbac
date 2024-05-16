package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/metadata"
	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	utils_v1 "gitlab.calendaria.team/services/utils/api/utils/v1"
	"gitlab.calendaria.team/services/utils/v2/auth"
)

type AssignsService struct {
	v1.UnimplementedAssignsServer

	ru *biz.RolesUsecase
	tu *biz.TeamsUsecase
	uc *biz.AssignedRolesUsecase
	sh *ServiceHelper
}

func NewAssignsService(
	ru *biz.RolesUsecase,
	tu *biz.TeamsUsecase,
	uc *biz.AssignedRolesUsecase,
	sh *ServiceHelper,
) *AssignsService {
	return &AssignsService{
		ru: ru,
		tu: tu,
		uc: uc,
		sh: sh,
	}
}

func (s *AssignsService) AssignRoles(ctx context.Context, req *v1.AssignRolesRequest) (*utils_v1.EmptyReply, error) {
	tenantId := auth.GetTenantIdFromContext(ctx)
	if tenantId == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	isAdmin := false
	if md, ok := metadata.FromServerContext(ctx); ok {
		isAdmin = md.Get("x-md-global-actor-role") == "admin"
	}

	if !isAdmin {
		_, _, err := s.sh.HasPermission(ctx, "admin.role.assign")
		if err != nil {
			return nil, err
		}
	}

	err := s.uc.AssignRoles(ctx, tenantId, toDtos(req))
	if err != nil {
		return nil, err
	}

	return &utils_v1.EmptyReply{}, nil
}

func (s *AssignsService) AssignRole(ctx context.Context, req *v1.AssignRoleRequest) (*utils_v1.EmptyReply, error) {
	tenantId := auth.GetTenantIdFromContext(ctx)
	if tenantId == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	isAdmin := false
	if md, ok := metadata.FromServerContext(ctx); ok {
		isAdmin = md.Get("x-md-global-actor-role") == "admin"
	}

	if !isAdmin {
		_, _, err := s.sh.HasPermission(ctx, "admin.role.assign")
		if err != nil {
			return nil, err
		}
	}

	err := s.uc.AssignRole(ctx, tenantId, toDto(req))
	if err != nil {
		return nil, err
	}

	return &utils_v1.EmptyReply{}, nil
}

func (s *AssignsService) UnassignRole(ctx context.Context, req *v1.AssignRequest) (*utils_v1.EmptyReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.role.assign")
	if err != nil {
		return nil, err
	}

	err = s.uc.UnassignRole(ctx, tenantId, req.GetAssignId())
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *AssignsService) ListAssigns(ctx context.Context, req *v1.ListAssignsRequest) (*v1.ListAssignsReply, error) {
	tenantId, _, err := s.sh.HasPermission(ctx, "admin.role.assign")
	if err != nil {
		return nil, err
	}

	teamsIDs := []int64{}

	if req.GetTeamId() != 0 {
		teamsIDs = []int64{req.GetTeamId()}
	}

	assignedRoles, err := s.uc.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: req.GetIdentityIds(),
		TeamsIDs:    teamsIDs,
	})
	if err != nil {
		return nil, err
	}

	return &v1.ListAssignsReply{
		Roles: assignedRolesReply(assignedRoles),
	}, nil
}

func assignedRoleReply(assignedRole *ent.TeamIdentityRole) *v1.AssignedRole {
	result := v1.AssignedRole{
		AssignId:   assignedRole.ID,
		IdentityId: &assignedRole.IdentityID,
	}

	if assignedRole.Edges.Team != nil {
		result.Team = &v1.Team{
			Id:          assignedRole.Edges.Team.ID,
			Name:        assignedRole.Edges.Team.Name,
			Description: assignedRole.Edges.Team.Description,
		}
	}

	if assignedRole.Edges.Role != nil {
		result.Role = &v1.Role{
			Id:          assignedRole.Edges.Role.ID,
			Name:        assignedRole.Edges.Role.Name,
			Description: assignedRole.Edges.Role.Description,
			IsSystem:    assignedRole.Edges.Role.IsSystem,
		}
	}

	return &result
}

func assignedRolesReply(assignedRoles []*ent.TeamIdentityRole) []*v1.AssignedRole {
	result := make([]*v1.AssignedRole, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		result[i] = assignedRoleReply(assignedRole)
	}
	return result
}

func toDtos(req *v1.AssignRolesRequest) []data.AssignRoleDto {
	dtos := make([]data.AssignRoleDto, len(req.Assigns))
	for i, assign := range req.Assigns {
		dtos[i] = data.AssignRoleDto{
			RoleId:     assign.GetRoleId(),
			TeamId:     assign.GetTeamId(),
			IdentityId: assign.GetIdentityId(),
		}
	}

	return dtos
}

func toDto(req *v1.AssignRoleRequest) data.AssignRoleDto {
	return data.AssignRoleDto{
		RoleId:     req.GetRoleId(),
		TeamId:     req.GetTeamId(),
		IdentityId: req.GetIdentityId(),
	}
}
