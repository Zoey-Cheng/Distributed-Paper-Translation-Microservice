package ocr

import (
	"github.com/alibabacloud-go/ocr-api-20210707/client" // 阿里云OCR SDK
	"go-micro.dev/v4/config"                             // 配置管理器
)

/**
* 从配置中获取阿里云OCR客户端
* @param config - 配置管理器
* @return OCR客户端
 */
func NewAliYunOCR(config config.Config) *client.Client {

	region := config.Get("aliyun", "ocr", "region").String("cn-hangzhou") // 获取region
	keyID := config.Get("aliyun", "ocr", "key_id").String("default")      // 获取key id
	secret := config.Get("aliyun", "ocr", "secret").String("default")     // 获取secret

	return MustGetAliYunOCR(region, keyID, secret) // 获取客户端
}
