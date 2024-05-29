package biz

import (
	u_nats "gitlab.calendaria.team/services/utils/v1/nats"

	"github.com/google/wire"
)

var (
	QueueRoleAssign   = "role_assign"
	QueueRoleUnassign = "role_unassign"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	u_nats.NewQueueManager,
	NewPermissionsUsecase,
	NewRolesUsecase,
	NewAssignedRolesUsecase,
	NewCheckPermissionsUsecase,
	NewTeamsUsecase,
)
