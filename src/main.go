package main

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/handler"
	logger2 "chatgpt-plus-exts/logger"
	"chatgpt-plus-exts/store"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

var logger = logger2.GetLogger()

// AppLifecycle 应用程序生命周期
type AppLifecycle struct {
}

// OnStart 应用程序启动时执行
func (l *AppLifecycle) OnStart(context.Context) error {
	log.Println("AppLifecycle OnStart")
	return nil
}

// OnStop 应用程序停止时执行
func (l *AppLifecycle) OnStop(context.Context) error {
	log.Println("AppLifecycle OnStop")
	return nil
}

func main() {
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "config.toml"
	}
	logger.Info("Loading config file: ", configFile)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Panic Error:", err)
		}
	}()

	app := fx.New(
		// 初始化配置应用配置
		fx.Provide(func() *core.Config {
			appConfig, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			appConfig.Path = configFile
			return appConfig
		}),

		// 初始化数据库
		fx.Provide(store.NewLevelDB),

		// 创建应用服务
		fx.Provide(func(config *core.Config) *core.AppServer {
			return core.NewServer(config)
		}),

		fx.Provide(handler.NewMidJourneyHandler),
		// 注册路由
		fx.Invoke(func(s *core.AppServer, h *handler.MidJourneyHandler) {
			group := s.Engine.Group("/api/mj/")
			group.POST("image", h.Image)
		}),

		fx.Invoke(func(s *core.AppServer) {
			err := s.Run()
			if err != nil {
				log.Fatal(err)
			}
		}),

		// 注册生命周期回调函数
		fx.Invoke(func(lifecycle fx.Lifecycle, lc *AppLifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return lc.OnStart(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return lc.OnStop(ctx)
				},
			})
		}),
	)
	// 启动应用程序
	go func() {
		if err := app.Start(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭应用程序
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Stop(ctx); err != nil {
		log.Fatal(err)
	}

}
