package ds

import (
	"context"
	"paper-translation/pkg/errutil"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/**
* 获取MongoDB客户端实例
*
* @param mongoUri - MongoDB连接URI
* @return - MongoDB客户端实例指针
 */
func MustGetMongoClient(mongoUri string) *mongo.Client {

	clientOptions := options.Client().ApplyURI(mongoUri) // 创建客户端配置选项

	client, err := mongo.Connect(context.TODO(), clientOptions) // 通过配置选项连接MongoDB

	errutil.PanicIfErr(err) // 如果连接错误,panic终止程序

	err = client.Ping(context.TODO(), nil) // 测试连接是否正常

	errutil.PanicIfErr(err) // 如果Ping错误,panic终止程序

	return client // 返回客户端实例
}
