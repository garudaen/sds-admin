package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sds-admin/internal/config"
	"sds-admin/internal/database"
	"sds-admin/internal/logger"
	"sds-admin/internal/router"
	"sds-admin/internal/server"

	"github.com/spf13/pflag"
)

// @title           SDS Admin API
// @version         1.0
// @description     SDS Admin API for domain management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

var (
	configPath string
	version    = "1.0.0"
	gitCommit  = "unknown"
	buildTime  = "unknown"
)

func init() {
	pflag.StringVarP(&configPath, "config", "c", "configs/config.yaml", "Configuration file path / 配置文件路径")
	pflag.BoolP("version", "v", false, "Show version information / 显示版本信息")
	pflag.BoolP("help", "h", false, "Show help message / 显示帮助信息")
}

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if showVersion, _ := pflag.CommandLine.GetBool("version"); showVersion {
		fmt.Printf("SDS Admin Version: %s\n", version)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	if showHelp, _ := pflag.CommandLine.GetBool("help"); showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		fmt.Println("Using default configuration...")
		cfg = config.LoadDefault()
	}

	if err := logger.Init(&cfg.Log); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Starting SDS Admin application...")
	logger.Infof("Version: %s, Git Commit: %s", version, gitCommit)

	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := router.Setup(cfg)
	server.Init(&cfg.Server, router)

	// 启动服务器
	go func() {
		if err := server.Start(); err != nil {
			logger.Errorf("Failed to start server: %v", err)
			logger.Fatal("Server failed to start, exiting")
		}
	}()

	// 等待关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Received shutdown signal, shutting down...")

	if err := server.Shutdown(context.Background(), cfg.Server.ShutdownTimeout); err != nil {
		logger.Errorf("Server shutdown error: %v", err)
	}

	logger.Info("SDS Admin stopped")
}
