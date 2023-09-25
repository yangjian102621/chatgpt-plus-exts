package sd

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/store"
	"chatgpt-plus-exts/utils"
	"fmt"
	"github.com/imroc/req/v3"
	"io"
	"time"
)

type TaskConsumer struct {
	mq     *store.RedisQueue
	config *core.SDConfig
}

func (tc *TaskConsumer) Run() {
	logger.Info("Starting consume Stable Diffusion Drawing task...")
	client := req.C()
	for {
		var task TaskInfo
		err := tc.mq.LPop(&task)
		if err != nil {
			logger.Errorf("taking message with error: %v", err)
			continue
		}
		go func() {
			tc.consumeTask(task, client)
		}()
	}
}

func (tc *TaskConsumer) consumeTask(taskInfo TaskInfo, client *req.Client) {
	body := map[string]any{
		"data":         taskInfo.Data,
		"event_data":   taskInfo.EventData,
		"fn_index":     taskInfo.FnIndex,
		"session_hash": taskInfo.SessionHash,
	}

	var result = make(chan CBReq)
	go func() {
		var res struct {
			Data            []interface{} `json:"data"`
			IsGenerating    bool          `json:"is_generating"`
			Duration        float64       `json:"duration"`
			AverageDuration float64       `json:"average_duration"`
		}
		var cbReq = CBReq{TaskId: taskInfo.TaskId}
		response, err := client.R().SetBody(body).SetSuccessResult(&res).Post(tc.config.ApiURL + "/run/predict")
		if err != nil {
			cbReq.Message = err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		if response.IsErrorState() {
			bytes, _ := io.ReadAll(response.Body)
			cbReq.Message = string(bytes)
			cbReq.Success = false
			result <- cbReq
			return
		}

		var images []struct {
			Name   string      `json:"name"`
			Data   interface{} `json:"data"`
			IsFile bool        `json:"is_file"`
		}
		err = utils.ForceCovert(res.Data[0], &images)
		if err != nil {
			cbReq.Message = err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		var info map[string]any
		err = utils.JsonDecode(utils.InterfaceToString(res.Data[1]), &info)
		if err != nil {
			cbReq.Message = err.Error()
			cbReq.Success = false
			result <- cbReq
			return
		}

		//for k, v := range info {
		//	fmt.Println(k, " => ", v)
		//}
		cbReq.ImageName = images[0].Name
		cbReq.Seed = utils.InterfaceToString(info["seed"])
		cbReq.Success = true
		cbReq.Progress = 100
		result <- cbReq
		close(result)

	}()

	for {
		select {
		case value := <-result:
			logger.Info(value)
			// TODO: 回调 API 推送失败消息
			return
		default:
			var progressReq = map[string]any{
				"id_task":         taskInfo.TaskId,
				"id_live_preview": 1,
			}

			var progressRes struct {
				Active        bool        `json:"active"`
				Queued        bool        `json:"queued"`
				Completed     bool        `json:"completed"`
				Progress      float64     `json:"progress"`
				Eta           float64     `json:"eta"`
				LivePreview   string      `json:"live_preview"`
				IDLivePreview int         `json:"id_live_preview"`
				TextInfo      interface{} `json:"textinfo"`
			}
			response, err := client.R().SetBody(progressReq).SetSuccessResult(&progressRes).Post(tc.config.ApiURL + "/internal/progress")
			var cbReq = CBReq{TaskId: taskInfo.TaskId, Success: true}
			if err != nil {
				logger.Error(err)
				// TODO: 这里可以考虑设置失败重试次数
				return
			}

			if response.IsErrorState() {
				bytes, _ := io.ReadAll(response.Body)
				logger.Error(string(bytes))
				return
			}

			cbReq.ImageData = progressRes.LivePreview
			cbReq.Progress = progressRes.Progress
			fmt.Println("Progress: ", progressRes.Progress)
			fmt.Println("Image: ", progressRes.LivePreview)
			time.Sleep(time.Second)
		}
	}
}
