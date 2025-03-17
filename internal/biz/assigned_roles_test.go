package biz_test

import (
	"context"
	"testing"

	v1 "gitlab.calendaria.team/services/rbac/api/rbac/v1"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/biz"
	"gitlab.calendaria.team/services/rbac/internal/data"
	"gitlab.calendaria.team/services/rbac/internal/data/mock"
	u_nats_mock "gitlab.calendaria.team/services/utils/v2/nats/mock"
	"gitlab.calendaria.team/services/utils/v2/zap"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAssignedRolesUsecase_AssignRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	queue := u_nats_mock.NewMockIQueue(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())

	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, qm)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	identityID := "identity1"
	roleID := int64(1)
	roleID2 := int64(2)
	teamID := int64(1)

	// Positive case
	dto := data.AssignRoleDto{
		IdentityID: identityID,
		RoleID:     roleID,
	}
	role := &ent.Role{
		ID:          roleID,
		TenantID:    tenantID,
		Name:        "testName",
		Description: "testDesc",
	}
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID).Return(role, nil)
	assignedRepo.EXPECT().AssignRoles(ctx, tenantID, []data.AssignRoleDto{dto}).Return(nil)

	qm.EXPECT().GetLocal(biz.QueueRoleAssign).Return(queue)
	queue.EXPECT().Pub(
		biz.AssignRoleMessage{
			AssignRoleDto: dto,
			TenantID:      tenantID,
		},
	)

	err = uc.AssignRole(ctx, tenantID, dto)
	require.NoError(t, err)

	// Negative case
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID2).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(
		ctx, tenantID, data.AssignRoleDto{
			IdentityID: identityID,
			RoleID:     roleID2,
		},
	)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("role not found"), err)

	// Negative case
	roleRepo.EXPECT().GetRoleByID(ctx, tenantID, roleID).Return(role, nil)
	teamRepo.EXPECT().GetTeam(ctx, tenantID, teamID, false).Return(nil, &ent.NotFoundError{})

	err = uc.AssignRole(
		ctx, tenantID, data.AssignRoleDto{
			IdentityID: identityID,
			RoleID:     roleID,
			TeamID:     teamID,
		},
	)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("team not found"), err)
}

func TestAssignedRolesUsecase_UnassignRole(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())

	queue := u_nats_mock.NewMockIQueue(ctrl)
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, qm)
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
	assignedRepo.EXPECT().GetAssignedRoleByID(ctx, tenantID, assignID).Return(assignedRole, nil)
	assignedRepo.EXPECT().UnassignRole(ctx, assignedRole).Return(nil)

	qm.EXPECT().GetLocal(biz.QueueRoleUnassign).Return(queue)
	queue.EXPECT().Pub(
		biz.AssignRoleMessage{
			AssignRoleDto: data.AssignRoleDto{
				IdentityID: assignedRole.IdentityID,
				RoleID:     assignedRole.RoleID,
			},
			TenantID: tenantID,
		},
	)

	err = uc.UnassignRole(ctx, tenantID, assignID)
	require.NoError(t, err)

	// Negative case
	assignedRepo.EXPECT().GetAssignedRoleByID(ctx, tenantID, assignID2).Return(nil, &ent.NotFoundError{})

	err = uc.UnassignRole(ctx, tenantID, assignID2)
	require.Error(t, err)
	require.Equal(t, v1.ErrorNotFound("assigned role not found"), err)
}

func TestAssignedRolesUsecase_ListAssignedRoles(t *testing.T) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	qm := u_nats_mock.NewMockIQueueManager(ctrl)
	qm.EXPECT().AddConsumer(biz.QueueRoleAssignHandler, gomock.Any())
	qm.EXPECT().AddConsumer(biz.QueueRoleUnassignHandler, gomock.Any())
	assignedRepo := mock.NewMockAssignedRolesRepo(ctrl)
	roleRepo := mock.NewMockRoleRepo(ctrl)
	teamRepo := mock.NewMockTeamsRepo(ctrl)
	uc, err := biz.NewAssignedRolesUsecase(logger, assignedRepo, roleRepo, teamRepo, qm)
	require.NoError(t, err)

	ctx := context.Background()
	tenantID := int64(1)
	identityID := "identity1"

	// Positive case
	assignedRoles := []*ent.ResourceAccess{}
	dto := data.ListRolesDto{TenantID: tenantID, IdentityIDs: []string{identityID}}
	assignedRepo.EXPECT().ListAssignedRoles(ctx, dto).Return(assignedRoles, nil)

	roles, err := uc.ListAssignedRoles(ctx, dto)
	require.NoError(t, err)
	require.Equal(t, assignedRoles, roles)
}
