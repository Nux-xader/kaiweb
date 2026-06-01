package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	CAPTCHA_APIS = []string{
		"https://captcha.xhub.casa/api/",
		"https://api.capsolver.com/",
		"https://api.2captcha.com/",
		"https://api.anti-captcha.com/",
	}
	MAX_POLLING_TIME = 20 * time.Second
	POLLING_INTERVAL = 250 * time.Millisecond
)

type CreateTaskRequest struct {
	ClientKey    string `json:"clientKey"`
	Task         Task   `json:"task"`
	SoftId       int    `json:"softId"`
	LanguagePool string `json:"languagePool"`
}

type Task struct {
	Type      string `json:"type"`
	Body      string `json:"body"`
	MinLength int    `json:"minLength"`
	MaxLength int    `json:"maxLength"`
	Phrase    bool   `json:"phrase"`
	Case      bool   `json:"case"`
	Numeric   int    `json:"numeric"`
	Math      bool   `json:"math"`
	Comment   string `json:"commecnt"`
}

type CreateTaskResponse struct {
	ErrorID  int    `json:"errorId"`
	TaskID   int    `json:"taskId,omitempty"`
	Status   string `json:"status"`
	Solution struct {
		Text string `json:"text"`
	} `json:"solution"`
}

type GetTaskResultRequest struct {
	ApiKey string `json:"clientKey"`
	TaskID int    `json:"taskId"`
}

type GetBalanceResp struct {
	ErrorId int     `json:"errorId"`
	Balance float64 `json:"balance"`
}

func GetBalance(api, apiKey string) (balance GetBalanceResp, err error) {
	client := resty.New().SetTimeout(15 * time.Second)
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{"clientKey": apiKey}).
		Post(api)
	if err != nil {
		err = fmt.Errorf("failed call captcha API: %w", err)
		return
	}

	balance = GetBalanceResp{}
	err = json.Unmarshal(resp.Body(), &balance)
	if err != nil {
		if balance.Balance < 5 {
			WarningPrint("Saldo anda di bawah $5, berpotensi gagal membaca captcha!")
		}
	}

	return
}

func ImageOCR3RdParty(imgPath, b64Img, api, apiKey string) (result string, err error) {
	if b64Img == "" {
		imageData, err := os.ReadFile(imgPath)
		if err != nil {
			return "", fmt.Errorf("failed read captcha image: %w", err)
		}
		b64Img = base64.StdEncoding.EncodeToString(imageData)
	}

	client := resty.New().SetTimeout(15 * time.Second)

	OriPrint("Submit captcha to captcha service")
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(CreateTaskRequest{
			ClientKey: apiKey,
			Task: Task{
				Type:      "ImageToTextTask",
				Body:      b64Img,
				MinLength: 0,
				MaxLength: 0,
				Phrase:    false,
				Case:      true,
				Numeric:   0,
			},
			SoftId:       802,
			LanguagePool: "rn",
		}).
		Post(api + "createTask")

	if err != nil {
		return "", fmt.Errorf("failed call captcha API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("got HTTP status code: %d", resp.StatusCode())
	}

	var taskResp = &CreateTaskResponse{}
	if err := json.Unmarshal(resp.Body(), taskResp); err != nil {
		return "", err
	}
	// fmt.Println(string(resp.Body()))

	if taskResp.ErrorID != 0 {
		return "", errors.New(taskResp.Status)
	}

	if taskResp.ErrorID != 0 {
		return "", fmt.Errorf("createTask error: %d", taskResp.ErrorID)
	}

	if taskResp.Status == "ready" {
		return taskResp.Solution.Text, nil
	}

	taskID := taskResp.TaskID
	getURL := api + "getTaskResult"

	startTime := time.Now()
	OriPrint("Waiting result from captcha service")
	n := 0
	for time.Since(startTime) < MAX_POLLING_TIME {
		if n >= 70 {
			fmt.Println("It's look captcha slow.")
			n = 0
		} else {
			n++
		}
		time.Sleep(POLLING_INTERVAL)

		resultReq := GetTaskResultRequest{
			ApiKey: apiKey,
			TaskID: taskID,
		}

		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(resultReq).
			Post(getURL)

		if err != nil {
			fmt.Printf("Failed call API getTaskResult: %v. Trying again...\n", err)
			continue
		}

		if resp.StatusCode() != 200 {
			return "", fmt.Errorf("got HTTP status code %d from getTaskResult", resp.StatusCode())
		}

		var resultResp CreateTaskResponse
		if err := json.Unmarshal(resp.Body(), &resultResp); err != nil {
			return "", fmt.Errorf("failed parse getTaskResult resp body: %w. Body: %s", err, string(resp.Body()))
		}

		if resultResp.ErrorID != 0 {
			return "", fmt.Errorf("getTaskResult error: %d", resultResp.ErrorID)
		}

		switch resultResp.Status {
		case "ready":
			return resultResp.Solution.Text, nil
		case "processing", "idle":
			continue
		default:
			return "", fmt.Errorf("unknown task status: %s", resultResp.Status)
		}
	}

	return "", errors.New("captcha task reach timeout")
}
