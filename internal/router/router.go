package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"yuxialuozi_graduation_design_backend/internal/config"
	"yuxialuozi_graduation_design_backend/internal/handler"
	"yuxialuozi_graduation_design_backend/internal/middleware"
)

var ProviderSet = wire.NewSet(NewRouter)

type Router struct {
	engine             *gin.Engine
	config             *config.Config
	authHandler        *handler.AuthHandler
	tenantHandler      *handler.TenantHandler
	contractHandler    *handler.ContractHandler
	roomHandler        *handler.RoomHandler
	feeHandler         *handler.FeeHandler
	maintenanceHandler *handler.MaintenanceHandler
	reportHandler      *handler.ReportHandler
}

func NewRouter(
	config *config.Config,
	authHandler *handler.AuthHandler,
	tenantHandler *handler.TenantHandler,
	contractHandler *handler.ContractHandler,
	roomHandler *handler.RoomHandler,
	feeHandler *handler.FeeHandler,
	maintenanceHandler *handler.MaintenanceHandler,
	reportHandler *handler.ReportHandler,
) *Router {
	if config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	r := &Router{
		engine:             engine,
		config:             config,
		authHandler:        authHandler,
		tenantHandler:      tenantHandler,
		contractHandler:    contractHandler,
		roomHandler:        roomHandler,
		feeHandler:         feeHandler,
		maintenanceHandler: maintenanceHandler,
		reportHandler:      reportHandler,
	}

	r.setupMiddlewares()
	r.setupRoutes()

	return r
}

func (r *Router) setupMiddlewares() {
	r.engine.Use(middleware.Recovery())
	r.engine.Use(middleware.Logger())
	r.engine.Use(middleware.CORS())
}

func (r *Router) setupRoutes() {
	// Swagger 文档路由
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.engine.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(r.config))
		{
			// Auth (protected)
			protected.GET("/auth/me", r.authHandler.GetCurrentUser)
			protected.POST("/auth/logout", r.authHandler.Logout)

			// Tenants
			tenants := protected.Group("/tenants")
			{
				tenants.GET("", r.tenantHandler.List)
				tenants.GET("/:id", r.tenantHandler.GetByID)
				tenants.POST("", r.tenantHandler.Create)
				tenants.PUT("/:id", r.tenantHandler.Update)
				tenants.DELETE("/:id", r.tenantHandler.Delete)
			}

			// Contracts
			contracts := protected.Group("/contracts")
			{
				contracts.GET("", r.contractHandler.List)
				contracts.GET("/:id", r.contractHandler.GetByID)
				contracts.POST("", r.contractHandler.Create)
				contracts.PUT("/:id", r.contractHandler.Update)
				contracts.DELETE("/:id", r.contractHandler.Delete)
			}

			// Rooms
			rooms := protected.Group("/rooms")
			{
				rooms.GET("", r.roomHandler.List)
				rooms.GET("/:id", r.roomHandler.GetByID)
				rooms.POST("", r.roomHandler.Create)
				rooms.PUT("/:id", r.roomHandler.Update)
				rooms.DELETE("/:id", r.roomHandler.Delete)
				rooms.POST("/:id/assign", r.roomHandler.AssignTenant)
			}

			// Fees
			fees := protected.Group("/fees")
			{
				fees.GET("", r.feeHandler.List)
				fees.GET("/:id", r.feeHandler.GetByID)
				fees.POST("", r.feeHandler.Create)
				fees.PUT("/:id", r.feeHandler.Update)
				fees.DELETE("/:id", r.feeHandler.Delete)
				fees.POST("/:id/pay", r.feeHandler.Pay)
			}

			// Maintenance
			maintenance := protected.Group("/maintenance")
			{
				maintenance.GET("", r.maintenanceHandler.List)
				maintenance.GET("/:id", r.maintenanceHandler.GetByID)
				maintenance.POST("", r.maintenanceHandler.Create)
				maintenance.PUT("/:id", r.maintenanceHandler.Update)
				maintenance.DELETE("/:id", r.maintenanceHandler.Delete)
				maintenance.POST("/:id/assign", r.maintenanceHandler.Assign)
				maintenance.POST("/:id/complete", r.maintenanceHandler.Complete)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/income", r.reportHandler.GetIncome)
				reports.GET("/occupancy", r.reportHandler.GetOccupancy)
				reports.GET("/fees/composition", r.reportHandler.GetFeeComposition)
				reports.GET("/maintenance/stats", r.reportHandler.GetMaintenanceStats)
				reports.GET("/tenants/ranking", r.reportHandler.GetTenantRanking)
				reports.GET("/dashboard", r.reportHandler.GetDashboard)
			}
		}
	}
}

func (r *Router) Run() error {
	return r.engine.Run(":8080")
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
