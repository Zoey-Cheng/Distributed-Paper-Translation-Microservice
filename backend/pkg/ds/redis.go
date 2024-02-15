package ds

import (
	"paper-translation/pkg/errutil" // 导入错误处理包

	"github.com/redis/go-redis/v9" // 导入redis客户端包
)

/**
* 获取redis客户端实例
*
* @param redisUri - redis连接URL
* @return - redis客户端实例指针
 */
func MustGetRedisClient(redisUri string) *redis.Client {

	opts, err := redis.ParseURL(redisUri) // 解析连接URL

	errutil.PanicIfErr(err) // 如果解析错误,panic

	return redis.NewClient(opts) // 创建并返回客户端
}
