package ocr

import (
	"paper-translation/pkg/errutil" // 错误处理包

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client" // 阿里云OpenAPI SDK
	"github.com/alibabacloud-go/ocr-api-20210707/client"             // 阿里云OCR SDK
	"github.com/alibabacloud-go/tea/tea"                             // 阿里云Tea工具库
)

/**
* 获取阿里云OCR客户端
* @param region - 区域
* @param keyID - AccessKey ID
* @param secret - AccessKey Secret
* @return OCR客户端
 */
func MustGetAliYunOCR(region string, keyID string, secret string) *client.Client {

	config := &openapi.Config{
		AccessKeyId:     &keyID,  // AccessKey ID
		AccessKeySecret: &secret, // AccessKey Secret
	}

	config.Endpoint = tea.String("ocr-api.cn-hangzhou.aliyuncs.com") // 设置Endpoint

	cli, err := client.NewClient(config) // 创建客户端

	errutil.PanicIfErr(err) // 错误处理

	return cli
}
