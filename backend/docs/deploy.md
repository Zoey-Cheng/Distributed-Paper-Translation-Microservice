# 部署

## 数据库

### 启动``etcd``
```shell
docker-compose up -d etcd
```

### 启动 ``mongo``
```shell
docker-compose up -d mongo
```

### 启动 ``redis``
```shell
docker-compose up -d redis
```

## 配置
### 初始化服务配置配置
```shell
make put-config
```

## 构建镜像
```shell
make images
```

## 启动服务
```shell
docker-compose up -d
```