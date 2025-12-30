package repository

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewTenantRepository,
	NewContractRepository,
	NewRoomRepository,
	NewFeeRepository,
	NewMaintenanceRepository,
)
