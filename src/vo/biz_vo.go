package vo

// BizVo 业务返回 VO
type BizVo struct {
	Code    BizCode     `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type BizCode int

const (
	Success = BizCode(0)
	Failed  = BizCode(1)

	InvalidArgs = "非法参数或参数解析失败"
	HackAttempt = "Hacking attempt!!!"
)
