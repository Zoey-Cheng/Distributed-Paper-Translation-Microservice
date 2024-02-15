# 开发注意

## `proto` 生成

### 安装protoc

#### github

``https://github.com/protocolbuffers/protobuf/releases``

#### mac

```shell
brew install protoc
```

#### windows

```shell
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.2/protoc-24.2-win64.zip
unzip protoc-24.2-win64.zip
## 添加bin目录到环境变量
```

#### linux

```shell
wget https://github.com/protocolbuffers/protobuf/releases/download/v24.2/protoc-24.2-linux-x86_64.zip
unzip protoc-24.2-linux-x86_64.zip
## 添加bin目录到环境变量
```

#### apt

```shell
sudo apt update
sudo apt install libprotobuf-dev protobuf-compiler
```

### 安装 `protoc-gen-go`、`protoc-gen-micro`、`wire`

```shell
make dependent
```

## ``wire``生成

```shell
make wire
```

如果添加新的需要wire管理依赖的服务，需要在``makefile->wire``
添加``wire ./app/${example}/service``
