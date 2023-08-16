package main

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/handler"
	logger2 "chatgpt-plus-exts/logger"
	"chatgpt-plus-exts/modules/mj"
	wexin "chatgpt-plus-exts/modules/weixin"
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
		// initialize app configs
		fx.Provide(func() *core.Config {
			config, err := core.LoadConfig(configFile)
			if err != nil {
				log.Fatal(err)
			}
			config.Path = configFile
			if config.Debug {
				_ = core.SaveConfig(config)
			}
			return config
		}),

		// initialize redis queue
		fx.Provide(store.NewRedisClient),
		fx.Provide(store.NewRedisMQs),

		// create app server
		fx.Provide(func(config *core.Config) *core.AppServer {
			return core.NewServer(config)
		}),

		// creating bots
		fx.Provide(func(config *core.Config, mqs *store.RedisMQs) *wexin.WeChatBot {
			if config.WeChatConfig.Enabled {
				return wexin.NewWeChatBot(config, mqs)
			}
			return nil
		}),
		fx.Invoke(func(config *core.Config, bot *wexin.WeChatBot) {
			if config.WeChatConfig.Enabled {
				go func() {
					err := bot.Run()
					if err != nil {
						logger.Error("微信登录失败：", err)
					}
				}()
				go func() {
					bot.ConsumeMessages()
				}()
			}
		}),

		fx.Provide(func(config *core.Config) *mj.MidJourneyClient {
			if config.MidJourneyConfig.Enabled {
				return mj.NewMjClient(config)
			}
			return nil
		}),
		fx.Provide(func(config *core.Config, mqs *store.RedisMQs) (*mj.MidJourneyBot, error) {
			if config.MidJourneyConfig.Enabled {
				return mj.NewMidJourneyBot(config, mqs)
			}
			return nil, nil
		}),
		fx.Invoke(func(config *core.Config, bot *mj.MidJourneyBot) {
			if config.MidJourneyConfig.Enabled {
				go func() {
					bot.Run()
				}()
				go func() {
					bot.ConsumeMessages()
				}()
			}
		}),

		// creating controller
		fx.Provide(handler.NewMidJourneyHandler),
		// register router
		fx.Invoke(func(s *core.AppServer, h *handler.MidJourneyHandler) {
			group := s.Engine.Group("/api/mj/")
			group.POST("image", h.Image)
			group.POST("upscale", h.Upscale)
			group.POST("variation", h.Variation)
		}),

		fx.Invoke(func(s *core.AppServer) {
			err := s.Run()
			if err != nil {
				log.Fatal(err)
			}
		}),

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
	// start web server
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
