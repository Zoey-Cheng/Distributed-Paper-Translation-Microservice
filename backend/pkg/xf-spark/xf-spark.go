package xf_spark

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HostUrl = "wss://aichat.xf-yun.com/v1/chat"
)

// Message 表示与 XFSpark 通信的消息结构。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// XFSparkClient 是与 XFSpark 服务通信的客户端。
type XFSparkClient struct {
	appID     string
	apiSecret string
	apiKey    string
}

// NewXFSparkClient 创建一个新的 XFSparkClient 实例。
//
// 参数:
// - appID (string): XFSpark 应用程序 ID。
// - apiSecret (string): XFSpark API 秘钥。
// - apiKey (string): XFSpark API 密钥。
//
// 返回值:
// - *XFSparkClient: XFSparkClient 实例。
func NewXFSparkClient(appID string, apiSecret string, apiKey string) *XFSparkClient {
	return &XFSparkClient{appID: appID, apiSecret: apiSecret, apiKey: apiKey}
}

// CreateChat 启动与 XFSpark 服务的对话。
//
// 参数:
// - ctx (context.Context): 上下文。
// - prompt (string): 用户的提示信息。
// - fc (func(text string)): 处理接收到的文本回调函数。
//
// 返回值:
// - error: 错误信息，如果发生错误。
func (t *XFSparkClient) CreateChat(ctx context.Context, prompt string, fc func(text string)) error {

	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}

	authUrl, err := t.assembleAuthUrl(HostUrl, t.apiKey, t.apiSecret)
	if err != nil {
		return err
	}

	conn, resp, err := dialer.DialContext(ctx, authUrl, nil)
	if err != nil || resp.StatusCode != http.StatusSwitchingProtocols {
		return errors.New(t.readResp(resp))
	}

	err = conn.WriteJSON(t.createParams(prompt))
	if err != nil {
		return err
	}
	return t.readMessages(conn, fc)
}

// readMessages 读取并处理来自 XFSpark 服务的消息。
//
// 参数:
// - conn (*websocket.Conn): WebSocket 连接。
// - fc (func(text string)): 处理接收到的文本回调函数。
//
// 返回值:
// - error: 错误信息，如果发生错误。
func (t *XFSparkClient) readMessages(conn *websocket.Conn, fc func(text string)) error {
	for {

		_, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		var data map[string]interface{}
		err = json.Unmarshal(msg, &data)
		if err != nil {
			return err
		}

		// 解析数据
		payload := data["payload"].(map[string]interface{})
		choices := payload["choices"].(map[string]interface{})
		header := data["header"].(map[string]interface{})
		code := header["code"].(float64)

		if code != 0 {
			log.Println(data["payload"])
			break
		}

		status := choices["status"].(float64)
		text := choices["text"].([]interface{})
		content := text[0].(map[string]interface{})["content"].(string)
		if len(content) > 0 {
			fc(content)
		}
		if status == 2 {
			conn.Close()
			break
		}
	}
	return nil
}

// createParams 创建用于与 XFSpark 服务通信的参数。
//
// 参数:
// - question (string): 用户的问题或提示信息。
//
// 返回值:
// - map[string]interface{}: 用于通信的参数。
func (t *XFSparkClient) createParams(question string) map[string]interface{} {

	messages := []Message{
		{Role: "user", Content: question},
	}

	data := map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
		"header": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"app_id": t.appID, // 根据实际情况修改返回的数据结构和字段名
		},
		"parameter": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"chat": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"domain":      "general",   // 根据实际情况修改返回的数据结构和字段名
				"temperature": 0.8,         // 根据实际情况修改返回的数据结构和字段名
				"top_k":       int64(6),    // 根据实际情况修改返回的数据结构和字段名
				"max_tokens":  int64(2048), // 根据实际情况修改返回的数据结构和字段名
				"auditing":    "default",   // 根据实际情况修改返回的数据结构和字段名
			},
		},
		"payload": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
			"message": map[string]interface{}{ // 根据实际情况修改返回的数据结构和字段名
				"text": messages, // 根据实际情况修改返回的数据结构和字段名
			},
		},
	}
	return data // 根据实际情况修改返回的数据结构和字段名
}

// assembleAuthUrl 组装用于进行身份验证的 URL。
//
// 参数:
// - hostURL (string): XFSpark 主机 URL。
// - apiKey (string): API 密钥。
// - apiSecret (string): API 秘钥。
//
// 返回值:
// - string: 组装后的身份验证 URL。
func (t *XFSparkClient) assembleAuthUrl(hostURL string, apiKey, apiSecret string) (string, error) {
	ul, err := url.Parse(hostURL)
	if err != nil {
		return "", err
	}
	date := time.Now().UTC().Format(time.RFC1123)
	signString := []string{"host: " + ul.Host, "date: " + date, "GET " + ul.Path + " HTTP/1.1"}
	sign := strings.Join(signString, "\n")
	sha := t.hmacWithSha256ToBase64(sign, apiSecret)
	authUrl := fmt.Sprintf("hmac username=\"%s\", algorithm=\"%s\", headers=\"%s\", signature=\"%s\"", apiKey,
		"hmac-sha256", "host date request-line", sha)
	authorization := base64.StdEncoding.EncodeToString([]byte(authUrl))
	v := url.Values{}
	v.Add("host", ul.Host)
	v.Add("date", date)
	v.Add("authorization", authorization)
	return hostURL + "?" + v.Encode(), nil
}

// hmacWithSha256ToBase64 使用 HMAC-SHA256 对数据进行签名，并返回 Base64 编码的签名结果。
//
// 参数:
// - data (string): 要签名的数据。
// - key (string): 签名所使用的密钥。
//
// 返回值:
// - string: Base64 编码的签名结果。
func (t *XFSparkClient) hmacWithSha256ToBase64(data, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	encodeData := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(encodeData)
}

// readResp 读取响应并返回字符串表示形式的响应信息。
//
// 参数:
// - resp (*http.Response): HTTP 响应。
//
// 返回值:
// - string: 字符串表示形式的响应信息。
func (t *XFSparkClient) readResp(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("code=%d,body=%s", resp.StatusCode, string(b))
}
