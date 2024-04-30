package biz

import (
	"github.com/google/wire"
	"gitlab.calendaria.team/services/utils/v1/nats"
)

var QueueRoleAssign = "role_assign"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	nats.NewQueueManager,
	NewPermissionsUsecase,
	NewRolesUsecase,
	NewAssignedRolesUsecase,
	NewCheckPermissionsUsecase,
	NewTeamsUsecase,
)
