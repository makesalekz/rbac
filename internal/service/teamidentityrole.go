package service

import (
	"context"
	"rbac/internal/biz"
	"rbac/internal/data"

	pb "rbac/api/rbac/v1"
)

type TeamIdentityRoleService struct {
	pb.UnimplementedTeamIdentityRoleServer

	uc *biz.TeamIdentityUsecase
}

func NewTeamIdentityRoleService(uc *biz.TeamIdentityUsecase) *TeamIdentityRoleService {
	return &TeamIdentityRoleService{
		uc: uc,
	}
}

func (s *TeamIdentityRoleService) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.TeamIdentityRoleReply, error) {
	identityRole, err := s.uc.AssignRole(ctx, data.AssignRoleDto{
		RoleId:     req.RoleId,
		TeamId:     req.TeamId,
		IdentityId: req.IdentityId,
		TenantId:   req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.TeamIdentityRoleReply{
		TenantId:   identityRole.TenantID,
		IdentityId: identityRole.IdentityID,
		Team: &pb.Team{
			Id:          identityRole.Edges.Team.ID,
			Name:        identityRole.Edges.Team.Name,
			Description: identityRole.Edges.Team.Description,
			TenantId:    identityRole.TenantID,
		},
		Role: &pb.RoleReply{
			Id:          identityRole.RoleID,
			Name:        identityRole.Edges.Role.Name,
			Description: identityRole.Edges.Role.Description,
			TenantId:    *identityRole.Edges.Role.TenantID,
			IsSystem:    identityRole.Edges.Role.IsSystem,
		},
	}, nil
}
func (s *TeamIdentityRoleService) DeleteRole(ctx context.Context, req *pb.DeleteRequest) (*pb.EmptyReply, error) {
	err := s.uc.DeleteIdentityRole(ctx, data.DeleteRoleDto{
		AssignId: req.AssignId,
	})
	if err != nil {
		return nil, err
	}
	return &pb.EmptyReply{}, nil
}
func (s *TeamIdentityRoleService) ListIdentityRoles(ctx context.Context, req *pb.ListIdentityRolesRequest) (*pb.ListIdentityRolesReply, error) {
	identityRoles, err := s.uc.ListIdentityRoles(ctx, data.ListIdentityRolesDto{
		IdentityId: req.IdentityId,
		TenantId:   req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*pb.TeamIdentityRoleReply, len(identityRoles))
	for i, identityRole := range identityRoles {
		result[i] = &pb.TeamIdentityRoleReply{
			TenantId:   identityRole.TenantID,
			IdentityId: identityRole.IdentityID,
			Team: &pb.Team{
				Id:          identityRole.Edges.Team.ID,
				Name:        identityRole.Edges.Team.Name,
				Description: identityRole.Edges.Team.Description,
				TenantId:    identityRole.TenantID,
			},
			Role: &pb.RoleReply{
				Id:          identityRole.RoleID,
				Name:        identityRole.Edges.Role.Name,
				Description: identityRole.Edges.Role.Description,
				TenantId:    *identityRole.Edges.Role.TenantID,
				IsSystem:    identityRole.Edges.Role.IsSystem,
			},
		}
	}
	return &pb.ListIdentityRolesReply{
		Roles: result,
	}, nil
}
func (s *TeamIdentityRoleService) ListTeamRoles(ctx context.Context, req *pb.ListTeamRolesRequest) (*pb.ListIdentityRolesReply, error) {
	teamRoles, err := s.uc.ListTeamRoles(ctx, data.ListTeamRolesDto{
		TeamId:   req.TeamId,
		TenantId: req.TenantId,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*pb.TeamIdentityRoleReply, len(teamRoles))
	for i, teamRole := range teamRoles {
		result[i] = &pb.TeamIdentityRoleReply{
			TenantId:   teamRole.TenantID,
			IdentityId: teamRole.IdentityID,
			Team: &pb.Team{
				Id:          teamRole.Edges.Team.ID,
				Name:        teamRole.Edges.Team.Name,
				Description: teamRole.Edges.Team.Description,
				TenantId:    teamRole.TenantID,
			},
			Role: &pb.RoleReply{
				Id:          teamRole.RoleID,
				Name:        teamRole.Edges.Role.Name,
				Description: teamRole.Edges.Role.Description,
				TenantId:    *teamRole.Edges.Role.TenantID,
				IsSystem:    teamRole.Edges.Role.IsSystem,
			},
		}
	}
	return &pb.ListIdentityRolesReply{
		Roles: result,
	}, nil
}
