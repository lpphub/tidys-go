package web

import (
	"context"
	"errors"
	"fmt"
	"tidys-go/infra"
	"tidys-go/logic"
	"tidys-go/web/rest/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/goweb/ext/logx"
	"github.com/lpphub/goweb/monitor"
)

type App struct {
	engine *gin.Engine
}

func New() *App {
	engine := gin.New()
	engine.Use(gin.Recovery())

	app := &App{
		engine: engine,
	}

	app.init()

	return app
}

func (a *App) init() {
	// 1.初始化基础设施
	infra.Init()
	// 2.初始化逻辑层
	logic.Init()

	// 3.配置web路由
	a.setupRouter()
}

func (a *App) setupRouter() {
	r := a.engine

	// pprof and metrics
	//monitor.StartPprof()
	monitor.RegisterMetrics(r)

	// 全局中间件
	r.Use(logx.GinAccessLog(logx.SkipPaths("/metrics", "/health")))

	// 注册所有接口路由
	handlers.RegisterRoutes(r)
}

func (a *App) Run() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", infra.Cfg.Server.Port),
		Handler: a.engine,
	}
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
