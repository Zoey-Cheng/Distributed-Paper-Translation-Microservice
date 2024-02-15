# 从alpine 3.18.3镜像作为基础镜像
FROM alpine:3.18.3

# 定义服务名称参数,后面复制可执行文件时使用
ARG service_name=''

# 将服务可执行文件复制到容器指定目录
COPY $service_name /usr/local/bin/$service_name  

# 设置工作目录
WORKDIR /usr/local/bin

# 设置etcd地址环境变量
ENV ETCD_ADDR=etcd:2379

# 设置配置文件前缀环境变量
ENV CONFIG_PREFIX=/configs/$service_name

# 定义服务启动命令环境变量
ENV Command="/usr/local/bin/$service_name"

# 容器启动时执行配置的启动命令
ENTRYPOINT $Command