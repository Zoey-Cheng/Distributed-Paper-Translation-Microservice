package handlers

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	fs "paper-translation/api/file/service/v1"
	"paper-translation/pkg/errutil"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// Bucket 是存储桶的名称。
const (
	Bucket = "forwork"
)

// ReqStartUpload 是开始上传文件的请求结构体。
type ReqStartUpload struct {
	Hash        string `json:"hash"`        // 文件哈希值
	FileName    string `json:"fileName"`    // 文件名
	ChunkNums   int64  `json:"chunkNums"`   // 总分块数
	SegmentSize int64  `json:"segmentSize"` // 分块大小
}

// ReqUploadChunk 是上传文件块的请求结构体。
type ReqUploadChunk struct {
	Chunk      *multipart.FileHeader `form:"chunk"`      // 文件块
	Hash       string                `form:"hash"`       // 文件哈希值
	ChunkIndex int64                 `form:"chunkIndex"` // 文件块索引
}

// FileHandler 是处理文件相关操作的处理器。
type FileHandler struct {
	fileService fs.FileService
	oss         *oss.Client
}

// NewFileHandler 创建一个新的 FileHandler 实例。
func NewFileHandler(fileService fs.FileService, oss *oss.Client) *FileHandler {
	return &FileHandler{fileService: fileService, oss: oss}
}

// StartUploadFile 处理开始上传文件的请求。
func (f *FileHandler) StartUploadFile(ctx *gin.Context) {
	var req ReqStartUpload
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errutil.ResponseError(ctx, errutil.RequestParamError, err)
		return
	}
	// 调用文件服务的 StartSegmentUpload 方法开始文件分块上传过程。
	_, err = f.fileService.StartSegmentUpload(ctx, &fs.SegmentUpload{
		Hash:        req.Hash,
		Filename:    req.FileName,
		ChunkNums:   req.ChunkNums,
		SegmentSize: req.SegmentSize,
		Bucket:      Bucket,
		FilePath:    fmt.Sprintf("files/%s-%s", uuid.NewString(), req.FileName),
	})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
}

// UploadChunk 处理上传文件块的请求。
func (f *FileHandler) UploadChunk(ctx *gin.Context) {

	var req ReqUploadChunk
	err := ctx.Bind(&req)
	if err != nil {
		errutil.ResponseError(ctx, errutil.RequestParamError, err)
		return
	}

	file, err := req.Chunk.Open()
	if err != nil {
		errutil.ResponseError(ctx, errutil.RequestParamError, err)
		return
	}

	bkt, err := f.oss.Bucket(Bucket)
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	key := fmt.Sprintf("chunks/%s/%d", req.Hash, req.ChunkIndex)
	// 将文件块上传到对象存储。
	err = bkt.PutObject(key, file)
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	// 标记文件块上传完成，并返回文件信息。
	fileInfo, err := f.fileService.MarkChunkOK(ctx, &fs.MarkChunk{
		Hash:       req.Hash,
		ChunkIndex: req.ChunkIndex,
	})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	// 如果文件已完全上传，将分块合并成完整文件。
	if fileInfo.Status == fs.FileStatus_Uploaded {
		result, err := bkt.ListObjectsV2(oss.Prefix(fmt.Sprintf("chunks/%s/", req.Hash)))
		if err != nil {
			errutil.ResponseError(ctx, errutil.UnknownError, err)
			return
		}

		var keys = make([]string, 0)
		for _, object := range result.Objects {
			keys = append(keys, object.Key)
		}

		// 对分块文件按索引排序，以确保正确顺序合并。
		sort.SliceStable(keys, func(i, j int) bool {
			ix, _ := strconv.ParseInt(strings.Split(keys[i], "/")[2], 10, 64)
			jx, _ := strconv.ParseInt(strings.Split(keys[j], "/")[2], 10, 64)
			return ix < jx
		})

		var appendPosition = int64(0)
		for _, objectKey := range keys {
			log.Printf("append file: %s", objectKey)
			obj, err := bkt.GetObject(objectKey)
			if err != nil {
				errutil.ResponseError(ctx, errutil.UnknownError, err)
				return
			}
			// 合并分块文件。
			nextPosition, err := bkt.AppendObject(*fileInfo.FilePath, obj, appendPosition)
			if err != nil {
				errutil.ResponseError(ctx, errutil.UnknownError, err)
				return
			}
			appendPosition = nextPosition
			obj.Close()
		}

		// 启动文件备份协程。
		go func() {
			// 文件备份
			_, _ = bkt.CopyObject(*fileInfo.FilePath, fmt.Sprintf("buckup/%s", *fileInfo.FilePath))
		}()
	}
}

// QueryFile 处理查询文件信息的请求。
func (f *FileHandler) QueryFile(ctx *gin.Context) {
	query, err := f.fileService.Query(ctx, &fs.QueryFile{Hash: ctx.Param("hash")})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
	ctx.JSON(200, query)
}

// GetFileURL 处理获取文件 URL 的请求。
func (f *FileHandler) GetFileURL(ctx *gin.Context) {
	query, err := f.fileService.Query(ctx, &fs.QueryFile{Hash: ctx.Param("hash")})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	if query.Status != fs.FileStatus_Uploaded {
		errutil.ResponseError(ctx, errutil.FileNotExistError)
		return
	}

	bkt, err := f.oss.Bucket(Bucket)
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	// 生成文件的可访问 URL 并返回给客户端。
	url, err := bkt.SignURL(*query.FilePath, http.MethodGet, 10000)
	if err != nil {
		return
	}
	ctx.JSON(200, bson.M{
		"url": url,
	})
}
