package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
)

const (
	_api    = "https://api.openai.com/v1/engines/chatgpt-4/completions"
	_apiKey = "OPENAI_API_KEY"
)

type Request struct {
	Prompt string `json:"prompt"`
	Length int    `json:"length"`
}

type Response struct {
	Completion string `json:"completion"`
}

/*
在这个文件中，我们定义了一个 handler 函数来处理传入的 HTTP 请求。
该函数会解码请求体，并使用 go-resty/resty/v2 包发送 POST 请求到 ChatGPT4 API。
响应数据将被反序列化为 Response 结构体，并通过 HTTP 响应返回。
*/
func handler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var reqBody Request
	err := decoder.Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := _api
	apiKey := os.Getenv(_apiKey)

	client := resty.New()
	respBody := Response{}
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+apiKey).
		SetBody(reqBody).
		SetResult(&respBody).
		Post(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}

func main() {
	http.HandleFunc("/chatgpt4", handler)
	http.ListenAndServe(":8080", nil)
}
