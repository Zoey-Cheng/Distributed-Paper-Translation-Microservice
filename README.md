# Distributed paper translation MicroService

###What it do?

Users can upload a foreign language PDF file over the Internet, select a target language, and leave an e-mail address, which is then processed and sent to the user's e-mail address.

Upon receipt, the server performs a series of processing and sends the translated file to the user's e-mail address.

**Keywords**: distributed storage + microservices



### Overall Architecture

The front-end and back-end are separated and deployed separately.

The backend is split into five microservices:

1. **paper-service** for file uploading. 
2. **email-service** Electronic mail delivery service
3. **ocr-service** file recognition service
4. **translation-service** translation service
5. **file-service** file upload service

The front-end framework used is **React**.



### Technology Selection

- **Golang** language, **micro-go** framework for microservice architecture
- **mongoDB** for key information storage
- **Redis** for caching
- **OSS** object storage for PDF file storage
- **OCR** port for file recognition and translation
- **SparkDesk LLM** for file translation and summarizing files.



### Project Structure

```
api         protobuf 定义
app         service实现
build       dockerfile
pkg         通用工具
config      配置文件
docs        文档
Makefile    构建脚本
docker-compose.yml  docker-compose配置文件
```



## Files

- Overall Deployment Document
  - `configuration.md`
- Back-end Document: `backend/README.md`
  - Deploy document:  `backend/docs/deploy.md`
  - Development document: ``backend/docs/development.md``
  - Gateway Deployment: ``backend/docs/gateway.md``

- Front-end Document: `frontend/README.md`

