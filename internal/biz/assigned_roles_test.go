package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	"gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAssignedRolesUsecase_AssignRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, nil)
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
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"
	roleId := int64(1)
	assignId := int64(1)
	assignId2 := int64(2)

	// Positive case
	assignedRole := &ent.ResourceAccess{
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

func TestAssignedRolesUsecase_ListAssignedRoles(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, nil)
	require.NoError(t, err)

	ctx := context.Background()
	tenantId := int64(1)
	identityId := "identity1"

	// Positive case
	assignedRoles := []*ent.ResourceAccess{}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, data.ListRolesDto{TenantId: tenantId, IdentityIDs: []string{identityId}}).Return(assignedRoles, nil)

	roles, err := uc.ListAssignedRoles(ctx, data.ListRolesDto{TenantId: tenantId, IdentityIDs: []string{identityId}})
	require.NoError(t, err)
	require.Equal(t, assignedRoles, roles)
}
