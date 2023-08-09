package resp

import (
	"chatgpt-plus-exts/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SUCCESS(c *gin.Context, values ...interface{}) {
	if values != nil {
		c.JSON(http.StatusOK, vo.BizVo{Code: vo.Success, Data: values[0]})
	} else {
		c.JSON(http.StatusOK, vo.BizVo{Code: vo.Success})
	}

}

func ERROR(c *gin.Context, messages ...string) {
	if messages != nil {
		c.JSON(http.StatusOK, vo.BizVo{Code: vo.Failed, Message: messages[0]})
	} else {
		c.JSON(http.StatusOK, vo.BizVo{Code: vo.Failed})
	}
}

func HACKER(c *gin.Context) {
	c.JSON(http.StatusOK, vo.BizVo{Code: vo.Failed, Message: "Hacker attempt!!!"})
}
