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

	dtos, err := toDtos(req)
	if err != nil {
		return nil, err
	}

	err = s.uc.AssignRoles(ctx, tenantId, dtos)
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

	dto, err := toDto(req)
	if err != nil {
		return nil, err
	}

	err = s.uc.AssignRole(ctx, tenantId, dto)
	if err != nil {
		return nil, err
	}

	return &utils_v1.EmptyReply{}, nil
}

func (s *AssignsService) UnassignRole(ctx context.Context, req *v1.AssignRequest) (*utils_v1.EmptyReply, error) {
	tenantId := auth.GetTenantIdFromContext(ctx)
	if tenantId == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	err := s.uc.UnassignRole(ctx, tenantId, req.GetAssignId())
	if err != nil {
		return nil, err
	}
	return &utils_v1.EmptyReply{}, nil
}

func (s *AssignsService) ListAssigns(ctx context.Context, req *v1.ListAssignsRequest) (*v1.ListAssignsReply, error) {
	tenantId := auth.GetTenantIdFromContext(ctx)
	if tenantId == 0 {
		return nil, v1.ErrorEmptyActorId("empty tenant id")
	}

	assignedRoles, err := s.uc.ListAssignedRoles(ctx, data.ListRolesDto{
		TenantId:       tenantId,
		IdentityIDs:    req.GetIdentityIds(),
		Resources:      req.GetResources(),
		ResourceFilter: req.GetResourceTypes(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.ListAssignsReply{
		Roles: assignedRolesReply(assignedRoles),
	}, nil
}

func assignedRoleReply(assignedRole *ent.ResourceAccess) *v1.AssignedRole {
	result := v1.AssignedRole{
		AssignId:   assignedRole.ID,
		IdentityId: &assignedRole.IdentityID,
	}

	if assignedRole.ResourceType != nil && assignedRole.ResourceID != nil {
		result.Resource = &v1.Resource{
			Type: *assignedRole.ResourceType,
			Id:   *assignedRole.ResourceID,
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

func assignedRolesReply(assignedRoles []*ent.ResourceAccess) []*v1.AssignedRole {
	result := make([]*v1.AssignedRole, len(assignedRoles))
	for i, assignedRole := range assignedRoles {
		result[i] = assignedRoleReply(assignedRole)
	}
	return result
}

func toDto(req *v1.AssignRoleRequest) (data.AssignRoleDto, error) {
	roleId := req.GetRoleId()
	teamId := req.GetTeamId()
	resource := req.GetResource()
	if teamId != 0 {
		resource = &v1.Resource{
			Type: data.RESOURCE_TYPE_TEAM,
			Id:   teamId,
		}
	}
	if roleId == 0 {
		return data.AssignRoleDto{}, v1.ErrorBadRequest("empty role id")
	}
	return data.AssignRoleDto{
		IdentityId: req.GetIdentityId(),
		RoleId:     roleId,
		TeamId:     teamId,
		Resource:   resource,
	}, nil
}

func toDtos(req *v1.AssignRolesRequest) ([]data.AssignRoleDto, error) {
	dtos := make([]data.AssignRoleDto, len(req.Assigns))
	for i, assign := range req.Assigns {
		dto, err := toDto(assign)
		if err != nil {
			return nil, err
		}
		dtos[i] = dto
	}
	return dtos, nil
}
