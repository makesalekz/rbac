package biz

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/nats-io/nats.go/jetstream"
	billing_v1 "gitlab.calendaria.team/services/finance/billing/messages"
	"gitlab.calendaria.team/services/rbac/ent"
	"gitlab.calendaria.team/services/rbac/internal/data"
	u_nats "gitlab.calendaria.team/services/utils/v2/nats"
)

const (
	QalaiBasicRoleID              = 7
	QalaiTrialRoleID              = 8
	QalaiPremiumSubscriptionQueue = "qalai-premium-features"
	QalaiExpiredSubscriptionQueue = "qalai-premium-features_revoke"
)

type ъPaidContent struct {
	log           *log.Helper
	assignedRoles data.AssignedRolesRepo
	qm            u_nats.IQueueManager
}

func NewPaidContent(
	qm u_nats.IQueueManager,
	assignedRoles data.AssignedRolesRepo,
	logger log.Logger,
) *PaidContent {
	pc := &PaidContent{
		log:           log.NewHelper(log.With(logger, "module", "biz/rbac")),
		qm:            qm,
		assignedRoles: assignedRoles,
	}

	qm.AddRemoteConsumer("finance-billing", QalaiPremiumSubscriptionQueue, pc.qalaiPremiumSubscriptionHandler)
	qm.AddRemoteConsumer("finance-billing", QalaiExpiredSubscriptionQueue, pc.qalaiExpiredSubscriptionHandler)

	return pc
}

func (uc *PaidContent) qalaiPremiumSubscriptionHandler(ctx context.Context, m jetstream.Msg) bool {
	var invoiceMessage billing_v1.RefreshItems

	err := json.Unmarshal(m.Data(), &invoiceMessage)
	if err != nil {
		uc.log.Errorf("failed to unmarshal invoice message: %v", err)

		return true
	}

	nonPremiumIDs := []int64{QalaiBasicRoleID, QalaiTrialRoleID}
	if invoiceMessage.IsTrial {
		nonPremiumIDs = []int64{QalaiBasicRoleID}
	}

	nonPremiumAssigns, err := uc.assignedRoles.GetAssignedRolesByRoleIDs(ctx, invoiceMessage.TenantID, nonPremiumIDs)
	if err != nil && !ent.IsNotFound(err) {
		uc.log.Errorf("failed to get non premium assign: %v", err)

		return false
	}

	for _, assign := range nonPremiumAssigns {
		err = uc.assignedRoles.UnassignRole(ctx, assign)

		if err != nil {
			uc.log.Errorf("failed to unassign role: %v", err)

			return false
		}
	}

	return true
}

func (uc *PaidContent) qalaiExpiredSubscriptionHandler(ctx context.Context, m jetstream.Msg) bool {
	var invoiceMessage billing_v1.RefreshItems

	err := json.Unmarshal(m.Data(), &invoiceMessage)
	if err != nil {
		uc.log.Errorf("failed to unmarshal invoice message: %v", err)

		return true
	}

	dtos := []data.AssignRoleDto{
		{
			RoleID: QalaiBasicRoleID,
		},
		{
			RoleID: QalaiTrialRoleID,
		},
	}

	err = uc.assignedRoles.AssignRoles(ctx, invoiceMessage.TenantID, dtos)
	if err != nil {
		uc.log.Errorf("failed to assign role: %v", err)

		return false
	}

	return true
}
