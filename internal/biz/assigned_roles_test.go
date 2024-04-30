package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAssignedRolesUsecase_AssignRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(assignedRepo, roleRepo, teamRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"
	roleId := int64(1)
	roleId2 := int64(2)
	teamId := int64(1)

	// Positive case
	dto := data.AssignRoleDto{
		IdentityId: identityId,
		RoleId:     roleId,
	}
	role := &ent.Role{
		ID:          roleId,
		TenantID:    tenantId,
		Name:        "testName",
		Description: "testDesc",
	}
	roleRepo.EXPECT().GetRoleById(ctx, tenantId, roleId).Return(role, nil)
	assignedRepo.EXPECT().AssignRoles(ctx, tenantId, []data.AssignRoleDto{dto}).Return(nil)

	err = uc.AssignRole(ctx, tenantId, dto)
	require.NoError(t, err)

	// Negative case
	roleRepo.EXPECT().GetRoleById(ctx, tenantId, roleId2).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(ctx, tenantId, data.AssignRoleDto{
		IdentityId: identityId,
		RoleId:     roleId2,
	})
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)

	// Negative case
	roleRepo.EXPECT().GetRoleById(ctx, tenantId, roleId).Return(role, nil)
	teamRepo.EXPECT().GetTeam(ctx, tenantId, teamId, false).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(ctx, tenantId, data.AssignRoleDto{
		IdentityId: identityId,
		RoleId:     roleId,
		TeamId:     teamId,
	})
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("team not found"), err)
}

func TestAssignedRolesUsecase_UnassignRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(assignedRepo, roleRepo, teamRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"
	roleId := int64(1)
	assignId := int64(1)
	assignId2 := int64(2)

	// Positive case
	assignedRole := &ent.TeamIdentityRole{
		ID:         assignId,
		TenantID:   tenantId,
		IdentityID: identityId,
		RoleID:     roleId,
	}
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantId, assignId).Return(assignedRole, nil)
	assignedRepo.EXPECT().UnassignRole(ctx, assignedRole).Return(nil)

	err = uc.UnassignRole(ctx, tenantId, assignId)
	require.NoError(t, err)

	// Negative case
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantId, assignId2).Return(nil, &ent.NotFoundError{})

	err = uc.UnassignRole(ctx, tenantId, assignId2)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("assigned role not found"), err)
}

func TestAssignedRolesUsecase_ListIdentityRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(assignedRepo, roleRepo, teamRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"
	roleId := int64(1)
	roleId2 := int64(2)

	// Positive case
	assignedRoles := []*ent.TeamIdentityRole{
		{
			ID:         1,
			TenantID:   tenantId,
			IdentityID: identityId,
			RoleID:     roleId,
		},
		{
			ID:         2,
			TenantID:   tenantId,
			IdentityID: identityId,
			RoleID:     roleId2,
		},
	}
	listRolesDto := data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: []string{identityId},
	}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, listRolesDto).Return(assignedRoles, nil)

	roles, err := uc.ListIdentityRoles(ctx, tenantId, identityId)
	require.NoError(t, err)
	require.Equal(t, assignedRoles, roles)
}

func TestAssignedRolesUsecase_ListAssignedRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(assignedRepo, roleRepo, teamRepo)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"

	// Positive case
	assignedRoles := []*ent.TeamIdentityRole{}
	listRolesDto := data.ListRolesDto{
		TenantId:    tenantId,
		IdentityIDs: []string{identityId},
	}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, listRolesDto).Return(assignedRoles, nil)

	roles, err := uc.ListAssignedRoles(ctx, listRolesDto)
	require.NoError(t, err)
	require.Equal(t, assignedRoles, roles)
}
