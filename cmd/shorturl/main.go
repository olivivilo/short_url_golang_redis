package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// TODO: Set up structured logging (consider using slog, zap, or zerolog)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Short URL Service...")

	// TODO: Load configuration
	// cfg, err := config.Load()
	// if err != nil {
	//     log.Fatalf("Failed to load configuration: %v", err)
	// }
	// log.Printf("Configuration loaded: server_port=%s, redis_addr=%s", cfg.Server.Port, cfg.Redis.Addr)

	// TODO: Initialize Redis client
	// redisClient := redis.NewClient(&redis.Options{
	//     Addr:     cfg.Redis.Addr,
	//     Password: cfg.Redis.Password,
	//     DB:       cfg.Redis.DB,
	//     PoolSize: cfg.Redis.PoolSize,
	// })
	// defer redisClient.Close()

	// TODO: Test Redis connection
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := redisClient.Ping(ctx).Err(); err != nil {
	//     log.Fatalf("Failed to connect to Redis: %v", err)
	// }
	// log.Println("Connected to Redis successfully")

	// TODO: Initialize dependencies
	// idGenerator := id.NewGenerator(redisClient, cfg.App.MinCodeLen)
	// urlRepo := redisrepo.NewURLRepository(redisClient)
	// urlService := service.NewURLService(urlRepo, idGenerator, cfg.App.BaseURL)

	// TODO: Initialize handlers
	// urlHandler := handler.NewURLHandler(urlService)
	// healthHandler := handler.NewHealthHandler(redisClient)

	// TODO: Set up router with middleware
	// router := httphandler.NewRouter(urlHandler, healthHandler)

	// TODO: Create HTTP server
	// server := &http.Server{
	//     Addr:         ":" + cfg.Server.Port,
	//     Handler:      router,
	//     ReadTimeout:  15 * time.Second,
	//     WriteTimeout: 15 * time.Second,
	//     IdleTimeout:  60 * time.Second,
	// }

	// TODO: Start server in a goroutine
	// go func() {
	//     log.Printf("Server listening on port %s", cfg.Server.Port)
	//     if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//         log.Fatalf("Server failed to start: %v", err)
	//     }
	// }()

	// TODO: Set up graceful shutdown
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutting down server...")

	// TODO: Graceful shutdown with timeout
	// shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.ShutdownTimeout)*time.Second)
	// defer shutdownCancel()

	// if err := server.Shutdown(shutdownCtx); err != nil {
	//     log.Fatalf("Server forced to shutdown: %v", err)
	// }

	// log.Println("Server exited gracefully")

	// Temporary placeholder
	log.Println("TODO: Implement main function")
	log.Println("Press Ctrl+C to exit")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Exiting...")
}
