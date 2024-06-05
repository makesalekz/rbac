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
	tenantID := int64(1)
	identityID := "identity1"
	roleID := int64(1)
	roleID2 := int64(2)
	teamID := int64(1)

	// Positive case
	dto := data.AssignRoleDto{
		IdentityId: identityID,
		RoleId:     roleID,
	}
	role := &ent.Role{
		ID:          roleID,
		TenantID:    tenantID,
		Name:        "testName",
		Description: "testDesc",
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID).Return(role, nil)
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, []data.AssignRoleDto{dto}).Return(nil)

	err = uc.AssignRole(ctx, tenantID, dto)
	require.NoError(t, err)

	// Negative case
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID2).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(ctx, tenantID, data.AssignRoleDto{
		IdentityId: identityID,
		RoleId:     roleID2,
	})
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)

	// Negative case
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID).Return(role, nil)
	teamRepo.EXPECT().GetTeam(ctx, tenantID, teamID, false).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(ctx, tenantID, data.AssignRoleDto{
		IdentityId: identityID,
		RoleId:     roleID,
		TeamId:     teamID,
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
	tenantID := int64(1)
	identityID := "identity1"
	roleID := int64(1)
	assignID := int64(1)
	assignID2 := int64(2)

	// Positive case
	assignedRole := &ent.ResourceAccess{
		ID:         assignID,
		TenantID:   tenantID,
		IdentityID: identityID,
		RoleID:     roleID,
	}
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantID, assignID).Return(assignedRole, nil)
	assignedRepo.EXPECT().UnassignRole(ctx, assignedRole).Return(nil)

	err = uc.UnassignRole(ctx, tenantID, assignID)
	require.NoError(t, err)

	// Negative case
	assignedRepo.EXPECT().GetAssignedRoleById(ctx, tenantID, assignID2).Return(nil, &ent.NotFoundError{})

	err = uc.UnassignRole(ctx, tenantID, assignID2)
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
	tenantID := int64(1)
	identityID := "identity1"

	// Positive case
	assignedRoles := []*ent.ResourceAccess{}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, data.ListRolesDto{TenantId: tenantID, IdentityIDs: []string{identityID}}).Return(assignedRoles, nil)

	roles, err := uc.ListAssignedRoles(ctx, data.ListRolesDto{TenantId: tenantID, IdentityIDs: []string{identityID}})
	require.NoError(t, err)
	require.Equal(t, assignedRoles, roles)
}
