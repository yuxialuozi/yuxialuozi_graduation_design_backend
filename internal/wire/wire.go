//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"

	"yuxialuozi_graduation_design_backend/internal/config"
	"yuxialuozi_graduation_design_backend/internal/database"
	"yuxialuozi_graduation_design_backend/internal/handler"
	"yuxialuozi_graduation_design_backend/internal/repository"
	"yuxialuozi_graduation_design_backend/internal/router"
	"yuxialuozi_graduation_design_backend/internal/service"
)

func InitializeApp() (*router.Router, func(), error) {
	wire.Build(
		config.ProviderSet,
		database.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		handler.ProviderSet,
		router.ProviderSet,
	)
	return nil, nil, nil
}
