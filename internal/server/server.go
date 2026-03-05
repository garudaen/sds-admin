package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"sds-admin/internal/config"
	"sds-admin/internal/logger"

	"github.com/gin-gonic/gin"
)

var httpServer *http.Server

// Init initializes the HTTP server with configuration.
// Init 使用配置初始化HTTP服务器
func Init(cfg *config.ServerConfig, router *gin.Engine) {
	gin.SetMode(cfg.Mode)

	httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      router,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

// Start starts the HTTP server.
// Start 启动HTTP服务器
func Start() error {
	if httpServer == nil {
		return fmt.Errorf("server not initialized")
	}

	addr := httpServer.Addr
	logger.Infof("Starting HTTP server on %s", addr)
	
	// 先尝试绑定端口，检查是否被占用
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// 同时在控制台输出错误信息
		fmt.Printf("ERROR: Failed to start server on %s: %v\n", addr, err)
		logger.Errorf("Failed to start server on %s: %v", addr, err)
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server with timeout.
// Shutdown 使用超时时间优雅关闭服务器
func Shutdown(ctx context.Context, timeout time.Duration) error {
	if httpServer == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	logger.Info("Shutting down HTTP server...")
	if err := httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Info("HTTP server stopped gracefully")
	return nil
}
