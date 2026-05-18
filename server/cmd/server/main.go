package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/freelanceplatform/server/internal/config"
	"github.com/freelanceplatform/server/internal/handler"
	"github.com/freelanceplatform/server/internal/middleware"
	"github.com/freelanceplatform/server/internal/pkg/logger"
	"github.com/freelanceplatform/server/internal/repository"
	"github.com/freelanceplatform/server/internal/service"
	"github.com/freelanceplatform/server/internal/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(cfg.Server.Mode)
	defer logger.Sync()
	log := logger.L

	log.Info("============================================")
	log.Info("  Freelance Platform Server Starting")
	log.Info("============================================")
	log.Infof("Mode: %s", cfg.Server.Mode)
	log.Infof("Port: %s", cfg.Server.Port)
	log.Infof("Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Infof("JWT Access Expiry: %dm, Refresh Expiry: %dd", cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)
	log.Infof("Storage: type=%s, path=%s", cfg.Storage.Type, cfg.Storage.Path)

	// Connect database
	ctx := context.Background()
	db, err := repository.NewDB(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()
	log.Info("Database connected successfully")

	// Run migrations
	if err := runMigrations(cfg.Database); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Info("Database migrations completed")

	// Set Gin mode
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		if err := db.Pool.Ping(c.Request.Context()); err != nil {
			handler.Error(c, http.StatusServiceUnavailable, 50300, "database unavailable")
			return
		}
		handler.Success(c, gin.H{"status": "ok"})
	})

	// Initialize repositories
	log.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	developerRepo := repository.NewDeveloperRepository(db)
	clientRepo := repository.NewClientRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	bidRepo := repository.NewBidRepository(db)
	contractRepo := repository.NewContractRepository(db)
	milestoneRepo := repository.NewMilestoneRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	conversationRepo := repository.NewConversationRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	notificationSettingsRepo := repository.NewNotificationSettingsRepository(db)
	fileRepo := repository.NewFileRepository(db)
	bookmarkRepo := repository.NewBookmarkRepository(db)
	disputeRepo := repository.NewDisputeRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	log.Info("Repositories initialized")

	// Initialize services
	log.Info("Initializing services...")
	authService := service.NewAuthService(cfg, userRepo, refreshTokenRepo, developerRepo, clientRepo)
	userService := service.NewUserService(userRepo, developerRepo, clientRepo)
	developerService := service.NewDeveloperService(developerRepo)
	clientService := service.NewClientService(clientRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	projectService := service.NewProjectService(projectRepo)
	bidService := service.NewBidService(bidRepo, projectRepo, contractRepo)
	contractService := service.NewContractService(contractRepo)
	milestoneService := service.NewMilestoneService(milestoneRepo, paymentRepo, walletRepo, contractRepo)
	paymentService := service.NewPaymentService(paymentRepo, walletRepo, contractRepo)
	walletService := service.NewWalletService(walletRepo)
	chatService := service.NewChatService(conversationRepo, messageRepo)
	reviewService := service.NewReviewService(reviewRepo, projectRepo, contractRepo)
	notificationService := service.NewNotificationService(notificationRepo, notificationSettingsRepo)

	// File service with local storage
	storageProvider := service.NewLocalStorageProvider(cfg.Storage.Path, "")
	fileService := service.NewFileService(fileRepo, storageProvider)

	bookmarkService := service.NewBookmarkService(bookmarkRepo)
	adminService := service.NewAdminService(adminRepo, userRepo, projectRepo, contractRepo, paymentRepo, disputeRepo)
	log.Info("Services initialized")

	// Initialize handlers
	log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService, reviewService)
	developerHandler := handler.NewDeveloperHandler(developerService, userService, fileService)
	clientHandler := handler.NewClientHandler(clientService, fileService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	projectHandler := handler.NewProjectHandler(projectService, fileService, bookmarkService, reviewService)
	bidHandler := handler.NewBidHandler(bidService)
	contractHandler := handler.NewContractHandler(contractService)
	milestoneHandler := handler.NewMilestoneHandler(milestoneService)
	paymentHandler := handler.NewPaymentHandler(paymentService, walletService)
	chatHandler := handler.NewChatHandler(chatService, fileService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	fileHandler := handler.NewFileHandler(fileService)
	adminHandler := handler.NewAdminHandler(adminService)
	log.Info("Handlers initialized")

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run()
	log.Info("WebSocket hub started")

	// API routes
	v1 := r.Group("/api/v1")

	// Public routes
	authHandler.RegisterRoutes(v1)
	categoryHandler.RegisterPublicRoutes(v1)
	projectHandler.RegisterPublicRoutes(v1)
	developerHandler.RegisterPublicRoutes(v1)
	reviewHandler.RegisterPublicRoutes(v1)

	// Protected routes
	auth := v1.Group("", middleware.JWTAuth(cfg.JWT.Secret))
	authHandler.RegisterRoutesProtected(auth)
	userHandler.RegisterRoutes(auth)
	developerHandler.RegisterRoutes(auth)
	clientHandler.RegisterRoutes(auth)
	projectHandler.RegisterRoutes(auth)
	bidHandler.RegisterRoutes(auth)
	contractHandler.RegisterRoutes(auth)
	milestoneHandler.RegisterRoutes(auth)
	paymentHandler.RegisterRoutes(auth)
	chatHandler.RegisterRoutes(auth)
	reviewHandler.RegisterRoutes(auth)
	notificationHandler.RegisterRoutes(auth)
	fileHandler.RegisterRoutes(auth)

	// Admin routes
	admin := v1.Group("/admin", middleware.JWTAuth(cfg.JWT.Secret), middleware.RequireAdmin(adminRepo))
	adminHandler.RegisterRoutes(admin)
	categoryHandler.RegisterRoutes(admin)

	// WebSocket
	r.GET("/api/v1/ws", func(c *gin.Context) {
		token := c.Query("token")
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 40100, "message": "invalid token"})
			return
		}
		ws.ServeWS(hub, c, claims.UserID)
	})

	log.Info("Routes registered")

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Infof("Server listening on %s", addr)
		log.Info("============================================")
		log.Info("  Server is ready to accept connections")
		log.Info("============================================")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Infof("Received signal %v, shutting down...", sig)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	db.Close()
	log.Info("Server exited gracefully")
}

func runMigrations(cfg config.DatabaseConfig) error {
	db, err := repository.NewDB(context.Background(), cfg)
	if err != nil {
		return fmt.Errorf("connect for migration: %w", err)
	}
	defer db.Close()

	// Check if migrations have already been applied
	var tableExists bool
	err = db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("check migration status: %w", err)
	}
	if tableExists {
		return nil // Migrations already applied
	}

	migrationSQL, err := os.ReadFile("migrations/000001_init_schema.up.sql")
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}

	return nil
}
