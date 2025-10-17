package http

import (
	"net/http"

	"github.com/M1r0-dev/Subscription-Aggregator/config"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/handler"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/mapper"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/middleware"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/controller/http/parser"
	"github.com/M1r0-dev/Subscription-Aggregator/internal/usecase"
	"github.com/M1r0-dev/Subscription-Aggregator/pkg/logger"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// Swagger spec:
// @title       Subscription Aggregator
// @description -
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, u usecase.SubscriptionUsecase, l logger.Interface) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	//Metrics
	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("Subscription-Aggregator")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	//Swagger
	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	//K8s
	app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })
	app.Get("/readyz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	// parser and mapper
	subscriptionParser := parser.New(l)
	subscriptionMapper := mapper.New()
	
	subscriptionHandler := handler.New(u, l, subscriptionParser, subscriptionMapper)

	// API routes
	api := app.Group("/v1")
	{
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.Post("/", subscriptionHandler.Store)
		}
	}
}
