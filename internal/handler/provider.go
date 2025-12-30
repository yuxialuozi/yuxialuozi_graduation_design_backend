package handler

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewAuthHandler,
	NewTenantHandler,
	NewContractHandler,
	NewRoomHandler,
	NewFeeHandler,
	NewMaintenanceHandler,
	NewReportHandler,
)
