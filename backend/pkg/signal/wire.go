package signal

import "go-micro.dev/v4/config"

/**
 * NewSignalFactory 创建并返回一个信号工厂。
 *
 * 参数:
 * - config (config.Config): 配置对象，用于获取 Redis 连接信息。
 *
 * 返回值:
 * - SignalFactory: 信号工厂实例。
 */
func NewSignalFactory(config config.Config) SignalFactory {
	// 从配置中获取 Redis 连接信息的 URI，默认为 "redis://localhost:6379"。
	redisURI := config.Get("redis", "uri").String("redis://localhost:6379")

	// 使用获取的 Redis URI 创建一个新的 RedisSignalFactory 实例。
	factory, err := NewRedisSignalFactory(redisURI)
	if err != nil {
		panic(err)
	}

	return factory
}
