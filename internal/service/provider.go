package service

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewAuthService,
	NewTenantService,
	NewContractService,
	NewRoomService,
	NewFeeService,
	NewMaintenanceService,
	NewReportService,
)
