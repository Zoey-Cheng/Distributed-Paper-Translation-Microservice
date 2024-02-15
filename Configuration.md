# 在Golang中使用Protobuf的步骤

Protobuf是一个序列化结构数据的工具,可以用来在不同语言之间交换数据。要在Golang中使用Protobuf,主要需要以下步骤:

## 1. 定义.proto文件

在.proto文件中定义你需要序列化的消息结构,指定语法版本,包名等信息。

## 2. 使用protoc生成代码

使用protoc编译器根据.proto文件生成具体编程语言的代码,例如:

    ```
    protoc --go_out=. *.proto
    ```

这个命令会根据proto文件生成Golang代码。

## 3. 使用生成的代码

可以在代码中导入生成的pb.go文件,然后使用里面的类型来序列化和反序列化数据,例如:

```go
import "xxx/pb"

// 序列化
data, err := proto.Marshal(&pb.Person{Name:"John"}) 

// 反序列化
person := &pb.Person{}
proto.Unmarshal(data, person)
```

这样就可以在Golang代码中方便地使用Protobuf了。

总结一下,主要步骤是定义.proto文件,使用protoc生成目标语言代码,然后导入使用即可。



# 全局配置文件

## 邮件服务配置

```markdown
{
  "email": {
    "address": "smtp.163.com:25", 
    "username": "golangproject@163.com",
    "password": "GGFWFZWFQRDAKVYV",
    "host": "smtp.163.com"
  }
}
```

配置邮件服务器的地址、用户名、密码和主机名。

## 文件服务mongo和redis配置

```json
{
  "mongo": {
    "uri": "mongodb://mongo:I8380LHI6P0GMUVV@mongo:27017",
    "db-name": "files"
  },
  "redis": {  
    "uri": "redis://redis:6379"
  }
}
```

配置文件服务使用的mongo和redis地址。

## OCR服务配置

```json
{
  "mongo": {
    "uri": "mongodb://mongo:I8380LHI6P0GMUVV@mongo:27017",
    "db-name": "ocrs"
  },
  "redis": {
    "uri": "redis://redis:6379"
  },
  "aliyun": {
    "ocr": {
      "region": "cn-hangzhou",
      "key_id": "LTAI5tJQyGRJqHzeU9DVHSt4", 
      "secret": "PcelPuPYKsFhfDculGYdColBKzmOUA"
    },
    "oss": {
      "region": "cn-beijing",
      "key_id": "LTAI5tJQyGRJqHzeU9DVHSt4",
      "secret": "PcelPuPYKsFhfDculGYdColBKzmOUA"
    }
  }
}
```

配置OCR服务使用的mongo、redis、阿里云OCR和OSS的相关参数。

## 论文服务配置

```json
{

  "mongo": {
    "uri": "mongodb://mongo:I8380LHI6P0GMUVV@mongo:27017",
    "db-name": "papers"
  },
  "redis": {
    "uri": "redis://redis:6379" 
  }

}
```

配置论文服务使用的mongo和redis地址。


## 翻译服务配置

```json
{
  "xf": {
    "appid": "93932e12",
    "secret": "MjEyNjc3ZDA0NWU4ODQ1MjFlMjM0OGUz",
    "key": "496bfc567df22e7cf65de37f8cd2aa00"
  },
  "redis": {
    "uri": "redis://redis:6379"
  }
}
```

配置翻译服务使用的讯飞平台和redis相关参数。



# Docker Compose配置说明

本文档对docker-compose.yml文件的配置进行说明。

## 服务列表

- etcd:etcd注册中心
- mongo:MongoDB数据库
- redis:redis数据库
- file-service:文件服务
- ocr-service:OCR服务
- translation-service:翻译服务
- paper-service:论文服务
- email-service:邮件服务
- frontend:前端服务

## 各服务详细配置

### etcd

使用bitnami/etcd镜像,容器名称为etcd。配置开放2379和2380端口,环境变量设置etcd的访问地址和允许空认证。

### mongo 

使用官方mongo镜像,容器名称为mongo。开放默认27017端口,设置环境变量配置用户名和密码。

### redis

使用bitnami/redis镜像,容器名为redis。开放6379端口,允许空密码。

### 各业务服务

业务服务均使用各自编译的镜像,容器名称为服务名。暴露4000端口供内部访问。

### frontend

使用编译的前端镜像,容器名为frontend。开放80端口作为 web 访问端口。

以上对配置文件的各服务作了概述说明,目的是帮助新同事快速理解系统组成。具体参数可参考配置文件,如有疑问欢迎随时讨论。



