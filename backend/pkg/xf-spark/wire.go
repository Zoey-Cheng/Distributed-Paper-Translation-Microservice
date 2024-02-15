package xf_spark

import "go-micro.dev/v4/config"

/**
 * NewXFSpark 根据配置创建并返回一个 XFSparkClient 实例。
 *
 * 参数:
 * - config (config.Config): 配置对象，包含了用于创建 XFSparkClient 的配置信息。
 *
 * 返回值:
 * - *XFSparkClient: XFSparkClient 实例。
 */
func NewXFSpark(config config.Config) *XFSparkClient {
	return NewXFSparkClient(
		config.Get("xf", "appid").String("c1d6b18e"),                          // 获取配置中的 App ID，默认为 "c1d6b18e"。
		config.Get("xf", "secret").String("ODFlNTBkMDI1NmU8ZGM2YmM5NzI8N2Q4"), // 获取配置中的 Secret，默认为 "ODFlNTBkMDI1NmU8ZGM2YmM5NzI8N2Q4"。
		config.Get("xf", "key").String("38a6e6344f781f9e927b30c62975737c"),    // 获取配置中的 Key，默认为 "38a6e6344f781f9e927b30c62975737c"。
	)
}
