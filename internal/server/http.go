package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	v1 "github.com/ch3nnn/webstack-go/api/v1"
	"github.com/ch3nnn/webstack-go/docs"
	"github.com/ch3nnn/webstack-go/internal/dal/model"
	"github.com/ch3nnn/webstack-go/internal/healthcheck"
	adminextHandler "github.com/ch3nnn/webstack-go/internal/handler/adminext"
	categoryHandler "github.com/ch3nnn/webstack-go/internal/handler/category"
	configHandler "github.com/ch3nnn/webstack-go/internal/handler/config"
	dashboardHandler "github.com/ch3nnn/webstack-go/internal/handler/dashboard"
	indexHandler "github.com/ch3nnn/webstack-go/internal/handler/index"
	publicapiHandler "github.com/ch3nnn/webstack-go/internal/handler/publicapi"
	siteHandler "github.com/ch3nnn/webstack-go/internal/handler/site"
	userHandler "github.com/ch3nnn/webstack-go/internal/handler/user"
	"github.com/ch3nnn/webstack-go/internal/middleware"
	"github.com/ch3nnn/webstack-go/pkg/jwt"
	"github.com/ch3nnn/webstack-go/pkg/log"
	httpx "github.com/ch3nnn/webstack-go/pkg/server/http"
	assets "github.com/ch3nnn/webstack-go/web"
)

func NewHTTPServer(
	engine *gin.Engine,
	logger *log.Logger,
	conf *viper.Viper,
	jwt *jwt.JWT,
	db *gorm.DB,
	dashboardHandler *dashboardHandler.Handler,
	indexHandler *indexHandler.Handler,
	userHandler *userHandler.Handler,
	siteHandler *siteHandler.Handler,
	categoryHandler *categoryHandler.Handler,
	configHandler *configHandler.Handler,
) *httpx.Server {
	gin.SetMode(gin.DebugMode)
	s := httpx.NewServer(
		engine,
		logger,
		httpx.WithServerHost(conf.GetString("http.host")),
		httpx.WithServerPort(conf.GetInt("http.port")),
	)
	s.MaxMultipartMemory = 32 << 20

	s.StaticFS("/assets", http.FS(assets.Static))
	s.SetHTMLTemplate(template.Must(template.New("").ParseFS(assets.Templates, "templates/**/*")))

	// Swagger Doc
	docs.SwaggerInfo.BasePath = "/v1"
	s.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerfiles.Handler,
		ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.GetInt("app.http.port"))),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true),
	))

	// Middleware
	s.Use(
		middleware.CORSMiddleware(),
		middleware.RequestLogMiddleware(logger),
		middleware.ResponseLogMiddleware(logger),
	)

	// New handlers (not Wire-managed)
	publicH := publicapiHandler.NewHandler(db, logger)
	adminExtH := adminextHandler.NewHandler(db, logger)

	// 404
	s.NoRoute(v1.ErrHandler404)

	// ========== Public Pages ==========
	s.GET("/", indexHandler.Index)
	s.GET("/about", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "about.html", nil)
	})
	s.GET("/site/:id", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "site_detail.html", nil)
	})
	s.GET("/submit", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "submit.html", nil)
	})
	s.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "admin_login.html", nil)
	})

	// ========== Admin HTML ==========
	render := s.Group("/admin").Use(middleware.NoStrictAuth(jwt, logger))
	{
		render.GET("", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_index.html", nil)
		})
		render.GET("dashboard", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "dashboard.html", nil)
		})
		render.GET("modify_password", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_modify_password.html", nil)
		})
		render.GET("category", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "category_view.html", nil)
		})
		render.GET("site", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "site_list.html", nil)
		})
		render.GET("site/add", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "site_add.html", nil)
		})
		render.GET("config", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "conf_web.html", nil)
		})
		render.GET("reviews", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_reviews.html", nil)
		})
		render.GET("users", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_users.html", nil)
		})
		render.GET("models", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_models.html", nil)
		})
		render.GET("tags", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_tags.html", nil)
		})
		render.GET("payments", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_payments.html", nil)
		})
		render.GET("submissions", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_submissions.html", nil)
		})
		render.GET("settings", func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "admin_settings.html", nil)
		})
	}

	// ========== API Routes ==========
	api := s.Group("/api")
	{
		// Public (no auth)
		pub := api.Group("")
		{
			pub.POST("/login", userHandler.Login)
			pub.GET("/about", indexHandler.About)
			pub.GET("/sites", publicH.SiteList)
			pub.GET("/sites/:id", publicH.SiteDetail)
			pub.GET("/sites/:id/reviews", publicH.ReviewList)
			pub.GET("/models", publicH.ModelList)
			pub.GET("/tags", publicH.TagList)
			pub.GET("/site-config", publicH.GetSiteConfig)
			pub.POST("/submissions", publicH.SubmitSite)
			pub.POST("/sites/:id/rating", publicH.SubmitRating)
			pub.POST("/sites/:id/review", publicH.SubmitReview)
			// Ads (public)
			pub.GET("/ads", publicH.AdList)
		}

		// Admin auth required
		strictAuthRouter := api.Group("/admin").Use(middleware.StrictAuth(jwt, logger))
		{
			// Dashboard
			strictAuthRouter.GET("/dashboard", dashboardHandler.Dashboard)
			// User
			strictAuthRouter.GET("/info", userHandler.Info)
			strictAuthRouter.POST("/logout", userHandler.Logout)
			strictAuthRouter.PATCH("/modify_password", userHandler.UpdatePassword)
			// Category
			strictAuthRouter.GET("/category", categoryHandler.List)
			strictAuthRouter.POST("/category", categoryHandler.Create)
			strictAuthRouter.PUT("/category/update", categoryHandler.Update)
			strictAuthRouter.GET("/category/:id", categoryHandler.Detail)
			strictAuthRouter.DELETE("/category/:id", categoryHandler.Delete)
			// Site
			strictAuthRouter.GET("/site", siteHandler.List)
			strictAuthRouter.GET("/site/sync/:id", siteHandler.SyncSite)
			strictAuthRouter.POST("/site", siteHandler.Create)
			strictAuthRouter.DELETE("/site/:id", siteHandler.Delete)
			strictAuthRouter.PUT("/site/:id", siteHandler.Update)
			strictAuthRouter.GET("/site/export", siteHandler.Export)
			// Config
			strictAuthRouter.GET("/config", configHandler.Config)
			strictAuthRouter.PUT("/config", configHandler.Update)

			// === Extended admin endpoints ===
			// Review moderation
			strictAuthRouter.GET("/reviews", adminExtH.ReviewList)
			strictAuthRouter.PUT("/reviews/:id", adminExtH.ReviewUpdate)
			// Rating management
			strictAuthRouter.GET("/ratings", adminExtH.RatingList)
			strictAuthRouter.DELETE("/ratings/:id", adminExtH.RatingDelete)
			strictAuthRouter.PUT("/ratings/:id", adminExtH.RatingUpdate)
			// User management
			strictAuthRouter.GET("/users", adminExtH.UserList)
			strictAuthRouter.PUT("/users/:id", adminExtH.UserUpdate)
			// Model dictionary
			strictAuthRouter.GET("/models", adminExtH.ModelList)
			strictAuthRouter.POST("/models", adminExtH.ModelCreate)
			strictAuthRouter.PUT("/models/:id", adminExtH.ModelUpdate)
			strictAuthRouter.DELETE("/models/:id", adminExtH.ModelDelete)
			// Tag dictionary
			strictAuthRouter.GET("/tags", adminExtH.TagList)
			strictAuthRouter.POST("/tags", adminExtH.TagCreate)
			strictAuthRouter.PUT("/tags/:id", adminExtH.TagUpdate)
			strictAuthRouter.DELETE("/tags/:id", adminExtH.TagDelete)
			// Payment method dictionary
			strictAuthRouter.GET("/payment-methods", adminExtH.PaymentMethodList)
			strictAuthRouter.POST("/payment-methods", adminExtH.PaymentMethodCreate)
			strictAuthRouter.PUT("/payment-methods/:id", adminExtH.PaymentMethodUpdate)
			strictAuthRouter.DELETE("/payment-methods/:id", adminExtH.PaymentMethodDelete)
			// Submission management
			strictAuthRouter.GET("/submissions", adminExtH.SubmissionList)
			strictAuthRouter.PUT("/submissions/:id", adminExtH.SubmissionUpdate)
			// Site settings (lang, anonymous, submission toggle)
			strictAuthRouter.PUT("/site-settings", adminExtH.UpdateSiteSettings)
			// Health check config
			strictAuthRouter.GET("/health-config", adminExtH.GetHealthConfig)
			strictAuthRouter.PUT("/health-config", adminExtH.UpdateHealthConfig)
			strictAuthRouter.POST("/health-trigger", adminExtH.TriggerHealthCheck)
			// Homepage layout config
			strictAuthRouter.GET("/homepage-config", adminExtH.GetHomepageConfig)
			strictAuthRouter.PUT("/homepage-config", adminExtH.UpdateHomepageConfig)
			// SEO config
			strictAuthRouter.GET("/seo-config", adminExtH.GetSEOConfig)
			strictAuthRouter.POST("/seo-config", adminExtH.UpdateSEOConfig)
			strictAuthRouter.PUT("/seo-config", adminExtH.UpdateSEOConfig)
			// Friend links
			strictAuthRouter.POST("/friend-links", adminExtH.UpdateFriendLinks)
			strictAuthRouter.PUT("/friend-links", adminExtH.UpdateFriendLinks)
			// Ads management
			strictAuthRouter.GET("/ads", adminExtH.AdList)
			strictAuthRouter.POST("/ads", adminExtH.AdCreate)
			strictAuthRouter.PUT("/ads/:id", adminExtH.AdUpdate)
			strictAuthRouter.DELETE("/ads/:id", adminExtH.AdDelete)
		}
	}

	// Start health checker
	var config model.SysConfig
	db.First(&config)
	interval := 30
	maxFails := 10
	if config.HealthCheckInterval > 0 {
		interval = config.HealthCheckInterval
	}
	if config.HealthCheckMaxFails > 0 {
		maxFails = config.HealthCheckMaxFails
	}
	checker := healthcheck.NewChecker(db, logger, interval, maxFails)
	if config.HealthCheckEnabled {
		checker.Start()
	}
	adminExtH.SetChecker(checker)

	return s
}