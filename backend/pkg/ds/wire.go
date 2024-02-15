package ds

import (
	"github.com/redis/go-redis/v9"      // redis客户端包
	"go-micro.dev/v4/config"            // micro配置管理包
	"go.mongodb.org/mongo-driver/mongo" // MongoDB驱动包
)

/**
* 从配置中获取redis客户端
*
* @param config - 配置管理器
* @return - redis客户端实例
 */
func NewRedisClient(config config.Config) *redis.Client {

	uri := config.Get("redis", "uri").String("redis://localhost:6379") // 从配置中获取redis uri

	return MustGetRedisClient(uri) // 使用uri获取客户端
}

/**
* 从配置中获取MongoDB客户端
*
* @param config - 配置管理器
* @return - MongoDB客户端实例
 */
func NewMongoClient(config config.Config) *mongo.Client {

	uri := config.Get("mongo", "uri").String("mongodb://localhost:27017") // 从配置中获取MongoDB uri

	return MustGetMongoClient(uri) // 使用uri获取客户端
}

/**
* 从配置和客户端中获取MongoDB数据库
*
* @param config - 配置管理器
* @param client - MongoDB客户端
* @return - MongoDB数据库实例
 */
func NewMongoDatabase(config config.Config, client *mongo.Client) *mongo.Database {

	name := config.Get("mongo", "db-name").String("paper-translation") // 从配置中获取数据库名

	return client.Database(name) // 从客户端中获取指定数据库
}
