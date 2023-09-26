package sd

import logger2 "chatgpt-plus-exts/logger"

var logger = logger2.GetLogger()

type Txt2ImgReq struct {
	TaskId         string  `json:"task_id"`
	Prompt         string  `json:"prompt"`
	NegativePrompt string  `json:"negative_prompt"`
	Steps          int     `json:"steps"`
	Sampler        string  `json:"sampler"`
	FaceFix        bool    `json:"face_fix"`
	CfgScale       float32 `json:"cfg_scale"`
	Seed           int64   `json:"seed"`
	Height         int     `json:"height"`
	Width          int     `json:"width"`
	HdFix          bool    `json:"hd_fix"`
	HdRedrawRate   float32 `json:"hd_redraw_rate"`
	HdScale        int     `json:"hd_scale"`
	HdScaleAlg     string  `json:"hd_scale_alg"`
	HdSampleNum    int     `json:"hd_sample_num"`
}

type TaskInfo struct {
	TaskId      string      `json:"task_id"`
	Data        interface{} `json:"data"`
	EventData   interface{} `json:"event_data"`
	FnIndex     int         `json:"fn_index"`
	SessionHash string      `json:"session_hash"`
}

type CBReq struct {
	TaskId    string
	ImageName string
	ImageData string
	Progress  int
	Seed      string
	Success   bool
	Message   string
}

var ParamKeys = map[string]int{
	"task_id":         0,
	"prompt":          1,
	"negative_prompt": 2,
	"steps":           4,
	"sampler":         5,
	"face_fix":        6,
	"cfg_scale":       10,
	"seed":            11,
	"height":          17,
	"width":           18,
	"hd_fix":          19,
	"hd_redraw_rate":  20, //高清修复重绘幅度
	"hd_scale":        21, // 高清修复放大倍数
	"hd_scale_alg":    22, // 高清修复放大算法
	"hd_sample_num":   23, // 高清修复采样次数
}

const Text2ImgParamTemplate = `[
"",
"",
"",
[],
30,
"DPM++ SDE Karras",
false,
false,
1,
1,
7.5,
-1,
-1,
0,
0,
0,
false,
512,
512,
true,
0.7,
2,
"Latent",
10,
0,
0,
"Use same sampler",
"",
"",
[],
"None",
false,
"MultiDiffusion",
false,
true,
1024,
1024,
96,
96,
48,
4,
"None",
2,
false,
10,
1,
1,
64,
false,
false,
false,
false,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
0.4,
0.4,
0.2,
0.2,
"",
"",
"Background",
0.2,
-1,
false,
3072,
192,
true,
true,
true,
false,
null,
null,
null,
false,
"",
0.5,
true,
false,
"",
"Lerp",
false,
"🔄",
false,
false,
false,
false,
false,
false,
false,
false,
false,
"positive",
"comma",
0,
false,
false,
"",
"Seed",
"",
[],
"Nothing",
"",
[],
"Nothing",
"",
[],
true,
false,
false,
false,
0,
null,
null,
false,
null,
null,
false,
null,
null,
false,
50
]`
