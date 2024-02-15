package main

import (
	fs "paper-translation/api/file/service/v1"
	v1 "paper-translation/api/paper/service/v1"
	"paper-translation/app/frontend/service/handlers"
	"paper-translation/pkg/service"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

// NewService 创建并返回一个新的 Web 服务实例。
//
// 参数:
// - registry registry.Registry: 微服务注册表。
// - config config.Config: 配置信息。
// - engine *gin.Engine: Gin 引擎实例。
//
// 返回值:
// - web.Service: 创建的 Web 服务实例。
func NewService(registry registry.Registry, config config.Config, engine *gin.Engine) web.Service {
	svc := web.NewService(
		web.Name(service.FrontendName),                          // 设置服务名称
		web.Address(config.Get("server", "addr").String(":80")), // 设置服务地址
		web.Registry(registry),                                  // 设置微服务注册表
	)
	svc.Init(web.Handler(engine)) // 初始化服务并指定 Gin 引擎作为处理程序
	return svc                    // 返回创建的 Web 服务实例
}

// NewFileService 创建并返回一个新的文件服务实例。
//
// 参数:
// - registry registry.Registry: 微服务注册表。
//
// 返回值:
// - fs.FileService: 创建的文件服务实例。
func NewFileService(registry registry.Registry) fs.FileService {
	cli := client.NewClient(
		client.Registry(registry), // 使用微服务注册表创建客户端
	)
	return fs.NewFileService(service.FileServiceName, cli) // 创建并返回文件服务实例
}

// NewPaperService 创建并返回一个新的论文服务实例。
//
// 参数:
// - registry registry.Registry: 微服务注册表。
//
// 返回值:
// - v1.PaperService: 创建的论文服务实例。
func NewPaperService(registry registry.Registry) v1.PaperService {
	cli := client.NewClient(
		client.Registry(registry), // 使用微服务注册表创建客户端
	)
	return v1.NewPaperService(service.PaperServiceName, cli) // 创建并返回论文服务实例
}

// NewRoute 创建并返回一个新的 Gin 引擎路由。
//
// 参数:
// - fileService fs.FileService: 文件服务实例。
// - paperService v1.PaperService: 论文服务实例。
// - oss *oss.Client: 阿里云 OSS 客户端。
//
// 返回值:
// - *gin.Engine: 创建的 Gin 引擎路由。
func NewRoute(fileService fs.FileService, paperService v1.PaperService, oss *oss.Client) *gin.Engine {
	r := gin.Default()                                       // 创建默认的 Gin 引擎
	r.Use(cors.Default())                                    // 使用默认的 CORS 中间件
	fileHandler := handlers.NewFileHandler(fileService, oss) // 创建文件处理器
	files := r.Group("/v1/files")                            // 创建文件路由组
	files.POST("/start", fileHandler.StartUploadFile)        // 处理文件上传请求
	files.POST("/chunk", fileHandler.UploadChunk)            // 处理文件块上传请求
	files.GET("/:hash", fileHandler.QueryFile)               // 处理文件查询请求
	files.GET("/:hash/public_url", fileHandler.GetFileURL)   // 处理获取文件公共链接请求

	paperHandler := handlers.NewPaperHandler(paperService)            // 创建论文处理器
	papers := r.Group("/v1/papers")                                   // 创建论文路由组
	papers.POST("/", paperHandler.CreatePaper)                        // 处理创建论文请求
	papers.GET("/", paperHandler.GetPapers)                           // 处理获取论文列表请求
	papers.GET("/:id", paperHandler.GetPaper)                         // 处理获取单个论文请求
	papers.DELETE("/:id", paperHandler.DeletePaper)                   // 处理删除论文请求
	papers.GET("/:id/download_txt", paperHandler.DownloadPaperResult) // 处理下载论文文本结果请求
	return r                                                          // 返回创建的 Gin 引擎路由
}
