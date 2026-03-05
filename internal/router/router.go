package router

import (
	"net/http"
	"strconv"
	"time"

	_ "sds-admin/docs"
	"sds-admin/internal/config"
	"sds-admin/internal/database"
	"sds-admin/internal/handler"
	"sds-admin/internal/middleware"
	"sds-admin/internal/service"
	_ "sds-admin/internal/static"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup creates and configures the Gin router with all routes.
// Setup 创建并配置包含所有路由的Gin路由器
func Setup(cfg *config.Config) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	setupHealthRoutes(router)
	setupRootRoutes(router)
	setupV1Routes(router)
	if cfg.Swagger.Enabled {
		setupSwaggerRoutes(router)
	}

	// Serve static files from statik
	setupStaticRoutes(router)

	return router
}

// setupHealthRoutes sets up health check routes.
// setupHealthRoutes 设置健康检查路由
func setupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})
}

// setupRootRoutes sets up root routes.
// setupRootRoutes 设置根路由
func setupRootRoutes(router *gin.Engine) {
	// Root route will be handled by static files (index.html)
	// 根路由将由静态文件(index.html)处理
}

// setupV1Routes sets up v1 API routes.
// setupV1Routes 设置v1 API路由
func setupV1Routes(router *gin.Engine) {
	v1 := router.Group("/api/v1")

	domainService := service.NewDomainService(database.GetDB())
	domainHandler := handler.NewDomainHandler(domainService)

	domains := v1.Group("/domains")
	{
		domains.POST("", domainHandler.CreateDomain)
		domains.GET("", domainHandler.ListDomains)
		domains.GET("/:id", domainHandler.GetDomain)
		domains.PUT("/:id", domainHandler.UpdateDomain)
		domains.PATCH("/:id/config", domainHandler.UpdateDomainConfig)
		domains.DELETE("/:id", domainHandler.DeleteDomain)
		domains.POST("/:id/disable", domainHandler.DisableDomain)
		domains.POST("/:id/enable", domainHandler.EnableDomain)
	}

	recordTypeService := service.NewRecordTypeService(database.GetDB())
	recordTypeHandler := handler.NewRecordTypeHandler(recordTypeService)

	recordTypes := v1.Group("/record-types")
	{
		recordTypes.POST("", recordTypeHandler.CreateRecordType)
		recordTypes.GET("", recordTypeHandler.ListRecordTypes)
		recordTypes.GET("/:id", recordTypeHandler.GetRecordType)
		recordTypes.PUT("/:id", recordTypeHandler.UpdateRecordType)
		recordTypes.DELETE("/:id", recordTypeHandler.DeleteRecordType)
	}

	// Record routes nested under domain ID.
	// 记录路由嵌套在域名ID下
	recordService := service.NewRecordService(database.GetDB())
	recordHandler := handler.NewRecordHandler(recordService)

	records := v1.Group("/domains/:id/records")
	{
		records.POST("", recordHandler.CreateRecord)
		records.GET("", recordHandler.ListRecords)
		records.GET("/:recordId", recordHandler.GetRecord)
		records.PUT("/:recordId", recordHandler.UpdateRecord)
		records.DELETE("/:recordId", recordHandler.DeleteRecord)
		records.POST("/:recordId/disable", recordHandler.DisableRecord)
		records.POST("/:recordId/enable", recordHandler.EnableRecord)
	}
}

// setupSwaggerRoutes sets up Swagger documentation routes.
// setupSwaggerRoutes 设置Swagger文档路由
func setupSwaggerRoutes(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// setupStaticRoutes sets up static file routes using statik.
// setupStaticRoutes 使用statik设置静态文件路由
func setupStaticRoutes(router *gin.Engine) {
	// Open statik file system
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	// Serve assets directory
	router.GET("/assets/*filepath", func(c *gin.Context) {
		// 从assets路径中提取文件路径
		filePath := c.Param("filepath")

		// 尝试多种路径格式
		pathsToTry := []string{}

		// 提取文件名
		fileName := filePath
		if len(fileName) > 0 && fileName[0] == '/' {
			fileName = fileName[1:]
		}

		// 构建各种可能的路径
		pathsToTry = append(pathsToTry,
			"assets"+filePath,       // assets/index-CU6T5UVH.css
			"assets/"+fileName,      // assets/index-CU6T5UVH.css
			fileName,                // index-CU6T5UVH.css
			"/assets"+filePath,      // /assets/index-CU6T5UVH.css
			"pub/assets"+filePath,   // pub/assets/index-CU6T5UVH.css
			"pub/assets/"+fileName,  // pub/assets/index-CU6T5UVH.css
			"dist/assets"+filePath,  // dist/assets/index-CU6T5UVH.css
			"dist/assets/"+fileName, // dist/assets/index-CU6T5UVH.css
		)

		// 尝试所有可能的路径
		var file http.File
		var err error
		for _, path := range pathsToTry {
			file, err = statikFS.Open(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer file.Close()

		// 获取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 设置Content-Type
		contentType := "application/octet-stream"
		if len(filePath) >= 5 && filePath[len(filePath)-4:] == ".css" {
			contentType = "text/css"
		} else if len(filePath) >= 4 && filePath[len(filePath)-3:] == ".js" {
			contentType = "application/javascript"
		}

		// 提供文件内容
		c.Header("Content-Type", contentType)
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.DataFromReader(http.StatusOK, fileInfo.Size(), contentType, file, nil)
	})

	// Serve vite.svg
	router.GET("/vite.svg", func(c *gin.Context) {
		file, err := statikFS.Open("vite.svg")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "image/svg+xml")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.DataFromReader(http.StatusOK, fileInfo.Size(), "image/svg+xml", file, nil)
	})

	// Handle root path
	router.GET("/", func(c *gin.Context) {
		// 尝试打开index.html文件，尝试多种路径
		pathsToTry := []string{
			"index.html",      // 直接路径
			"/index.html",     // 带/前缀
			"pub/index.html",  // 带pub前缀
			"dist/index.html", // 带dist前缀
		}

		var file http.File
		var err error

		// 尝试所有可能的路径
		for _, path := range pathsToTry {
			file, err = statikFS.Open(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/html")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.DataFromReader(http.StatusOK, fileInfo.Size(), "text/html", file, nil)
	})

	// Handle single page application routes
	router.NoRoute(func(c *gin.Context) {
		// Skip API routes, health check, and swagger
		path := c.Request.URL.Path
		if path == "/health" ||
			len(path) >= 5 && path[:5] == "/api/" ||
			len(path) >= 8 && path[:8] == "/swagger" {
			c.Next()
			return
		}

		// Serve index.html for all other routes
		pathsToTry := []string{
			"index.html",      // 直接路径
			"/index.html",     // 带/前缀
			"pub/index.html",  // 带pub前缀
			"dist/index.html", // 带dist前缀
		}

		var file http.File
		var err error

		// 尝试所有可能的路径
		for _, path := range pathsToTry {
			file, err = statikFS.Open(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.Header("Content-Type", "text/html")
		c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		c.DataFromReader(http.StatusOK, fileInfo.Size(), "text/html", file, nil)
	})
}
