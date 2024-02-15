# 网关部署

## 需要挂载nginx.conf配置文件，提前创建挂载目录

``nginx.conf``
```shell
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;


events {
    worker_connections  1024;
}


http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;
    sendfile        off;
    keepalive_timeout  65;
    gzip  on;

    upstream frontend {
        server frontend:80 weight=1;
    }

    upstream web {
        server web:80 weight=1;
    }
    
    server {
        listen 80;
        server_name localhost;
        location /v1 {
            proxy_pass http://frontend;
        }

        location / {
            proxy_pass http://web;
        }
    }
}
```

# network
因为服务那边都是没有实际暴露端口的，所以需要绑定后端和前端的docker网络

### 启动命令
```shell
docker-compose up -d
```

`docker-compose.yaml`
```shell
version: '3.1'

services:

  nginx:
    image: nginx:alpine
    restart: always
    container_name: nginx
    ports:
      - "80:80"
    expose:
      - "80"
    environment:
      - NGINX_PORT=80
    networks: 
      - paper-translation_default
      - paper-translation-web_default
    volumes:
      - /root/gateway/data/www:/usr/share/nginx/html
      - /root/gateway/data/conf/nginx.conf:/etc/nginx/nginx.conf
      - /root/gateway/data/logs:/var/log/nginx
networks:
  paper-translation_default:
    external: true
  paper-translation-web_default:
    external: true
```

