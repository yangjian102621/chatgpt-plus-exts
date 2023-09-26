package sd

import (
	"chatgpt-plus-exts/store"
	"chatgpt-plus-exts/utils"
)

type StableDiffusionClient struct {
	mq *store.RedisQueue
}

func NewSdClient(mqs *store.RedisMQs) *StableDiffusionClient {
	return &StableDiffusionClient{
		mq: mqs.StableDiffusion,
	}
}

func (client *StableDiffusionClient) Txt2Img(param Txt2ImgReq) error {
	var data []interface{}
	err := utils.JsonDecode(Text2ImgParamTemplate, &data)
	if err != nil {
		return err
	}
	data[ParamKeys["task_id"]] = param.TaskId
	data[ParamKeys["prompt"]] = param.Prompt
	data[ParamKeys["negative_prompt"]] = param.NegativePrompt
	data[ParamKeys["steps"]] = param.Steps
	data[ParamKeys["sampler"]] = param.Sampler
	data[ParamKeys["face_fix"]] = param.FaceFix
	data[ParamKeys["cfg_scale"]] = param.CfgScale
	data[ParamKeys["seed"]] = param.Seed
	data[ParamKeys["height"]] = param.Height
	data[ParamKeys["width"]] = param.Width
	data[ParamKeys["hd_fix"]] = param.HdFix
	data[ParamKeys["hd_redraw_rate"]] = param.HdRedrawRate
	data[ParamKeys["hd_scale"]] = param.HdScale
	data[ParamKeys["hd_scale_alg"]] = param.HdScaleAlg
	data[ParamKeys["hd_sample_num"]] = param.HdSampleNum
	task := TaskInfo{
		TaskId:      param.TaskId,
		Data:        data,
		EventData:   nil,
		FnIndex:     494,
		SessionHash: "ycaxgzm9ah",
	}

	client.mq.RPush(task)
	return nil
}
