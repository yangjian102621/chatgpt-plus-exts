package core

import (
	logger2 "chatgpt-plus-exts/logger"
	"chatgpt-plus-exts/utils"
	"chatgpt-plus-exts/vo"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"runtime/debug"
)

type AppServer struct {
	Config *Config
	Engine *gin.Engine
}

var logger = logger2.GetLogger()

func NewServer(config *Config) *AppServer {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
	engine := gin.Default()
	engine.Use(errorHandler)
	return &AppServer{
		Config: config,
		Engine: engine,
	}
}

func (s *AppServer) Run() error {
	logger.Infof("http://%s", s.Config.Listen)
	return s.Engine.Run(s.Config.Listen)
}

// 全局异常处理
func errorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Handler Panic: %v\n", r)
			debug.PrintStack()
			c.JSON(http.StatusOK, vo.BizVo{Code: vo.Failed, Message: utils.InterfaceToString(r)})
			c.Abort()
		}
	}()
	//加载完 defer recover，继续后续接口调用
	c.Next()
}
