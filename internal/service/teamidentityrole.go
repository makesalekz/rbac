package service

import (
	"context"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
)

type TeamIdentityRoleService struct {
	v1.UnimplementedTeamIdentityRoleServer

	uc *biz.TeamIdentityUsecase
}

func NewTeamIdentityRoleService(uc *biz.TeamIdentityUsecase) *TeamIdentityRoleService {
	return &TeamIdentityRoleService{
		uc: uc,
	}
}

func (s *TeamIdentityRoleService) AssignRole(ctx context.Context, req *v1.AssignRoleRequest) (*v1.TeamIdentityRoleReply, error) {
	identityRole, err := s.uc.AssignRole(ctx, data.AssignRoleDto{
		RoleId:     req.RoleId,
		TeamId:     req.TeamId,
		IdentityId: req.IdentityId,
		TenantId:   req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.TeamIdentityRoleReply{
		TenantId:   identityRole.TenantID,
		IdentityId: identityRole.IdentityID,
		Team: &v1.Team{
			Id:          identityRole.Edges.Team.ID,
			Name:        identityRole.Edges.Team.Name,
			Description: identityRole.Edges.Team.Description,
			TenantId:    identityRole.TenantID,
		},
		Role: &v1.RoleReply{
			Id:          identityRole.RoleID,
			Name:        identityRole.Edges.Role.Name,
			Description: identityRole.Edges.Role.Description,
			TenantId:    *identityRole.Edges.Role.TenantID,
			IsSystem:    identityRole.Edges.Role.IsSystem,
		},
	}, nil
}
func (s *TeamIdentityRoleService) DeleteRole(ctx context.Context, req *v1.DeleteRequest) (*v1.EmptyReply, error) {
	err := s.uc.DeleteIdentityRole(ctx, data.DeleteRoleDto{
		AssignId: req.AssignId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}
func (s *TeamIdentityRoleService) ListIdentityRoles(ctx context.Context, req *v1.ListIdentityRolesRequest) (*v1.ListIdentityRolesReply, error) {
	identityRoles, err := s.uc.ListIdentityRoles(ctx, data.ListIdentityRolesDto{
		IdentityId: req.IdentityId,
		TenantId:   req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*v1.TeamIdentityRoleReply, len(identityRoles))
	for i, identityRole := range identityRoles {
		result[i] = &v1.TeamIdentityRoleReply{
			TenantId:   identityRole.TenantID,
			IdentityId: identityRole.IdentityID,
			Team: &v1.Team{
				Id:          identityRole.Edges.Team.ID,
				Name:        identityRole.Edges.Team.Name,
				Description: identityRole.Edges.Team.Description,
				TenantId:    identityRole.TenantID,
			},
			Role: &v1.RoleReply{
				Id:          identityRole.RoleID,
				Name:        identityRole.Edges.Role.Name,
				Description: identityRole.Edges.Role.Description,
				TenantId:    *identityRole.Edges.Role.TenantID,
				IsSystem:    identityRole.Edges.Role.IsSystem,
			},
		}
	}
	return &v1.ListIdentityRolesReply{
		Roles: result,
	}, nil
}
func (s *TeamIdentityRoleService) ListTeamRoles(ctx context.Context, req *v1.ListTeamRolesRequest) (*v1.ListIdentityRolesReply, error) {
	teamRoles, err := s.uc.ListTeamRoles(ctx, data.ListTeamRolesDto{
		TeamId:   req.TeamId,
		TenantId: req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*v1.TeamIdentityRoleReply, len(teamRoles))
	for i, teamRole := range teamRoles {
		result[i] = &v1.TeamIdentityRoleReply{
			TenantId:   teamRole.TenantID,
			IdentityId: teamRole.IdentityID,
			Team: &v1.Team{
				Id:          teamRole.Edges.Team.ID,
				Name:        teamRole.Edges.Team.Name,
				Description: teamRole.Edges.Team.Description,
				TenantId:    teamRole.TenantID,
			},
			Role: &v1.RoleReply{
				Id:          teamRole.RoleID,
				Name:        teamRole.Edges.Role.Name,
				Description: teamRole.Edges.Role.Description,
				TenantId:    *teamRole.Edges.Role.TenantID,
				IsSystem:    teamRole.Edges.Role.IsSystem,
			},
		}
	}
	return &v1.ListIdentityRolesReply{
		Roles: result,
	}, nil
}
