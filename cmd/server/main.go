package main

import (
	"flag"
	"log"
	infraRedis "rate-limiter-api/internal/infrastructure/redis"

	"rate-limiter-api/internal/config"
	"rate-limiter-api/internal/controller"
	corerequest "rate-limiter-api/internal/core/request"
	"rate-limiter-api/internal/helpers/timehelpers"
	"rate-limiter-api/internal/infrastructure/ratelimiter"
	"rate-limiter-api/internal/logger"
	"rate-limiter-api/internal/repository"
	"rate-limiter-api/router"
)

func main() {
	cfgPath := flag.String("config-path", "config/dev/config.yaml", "path to config file")
	flag.Parse()

	cfg := config.Load(*cfgPath)

	clock := timehelpers.NewClock()
	appLogger := logger.New(cfg.Settings.Logger.Level)

	limiter := buildLimiter(cfg)

	repo := repository.NewInMemoryRequestRepo()
	requestService := corerequest.NewService(repo, clock)

	srv := controller.New(controller.Services{Request: requestService}, appLogger)

	r := router.Init(srv, appLogger)
	router.RegisterV1(r, srv, limiter)

	port := cfg.Settings.Server.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func buildLimiter(cfg *config.Config) ratelimiter.RateLimiter {
	redisClient := infraRedis.New(
		cfg.Settings.Redis.Host,
		cfg.Settings.Redis.Password,
		cfg.Settings.Redis.DB,
	)
	return ratelimiter.NewRedisLimiter(redisClient)

}
