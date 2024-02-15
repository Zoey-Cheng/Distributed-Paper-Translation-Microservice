package file

import (
	context "context"
	v1 "paper-translation/api/file/service/v1"
	"paper-translation/pkg/lock"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

// FileService 提供文件相关的服务。
type FileService struct {
	repo        FileRepository
	redisClient *redis.Client
}

// NewFileService 创建一个新的 FileService 实例。
//
// 参数:
// - repo FileRepository: 文件仓库实例。
// - redisClient *redis.Client: Redis 客户端实例。
//
// 返回值:
// - *FileService: 创建的 FileService 实例。
func NewFileService(repo FileRepository, redisClient *redis.Client) *FileService {
	return &FileService{repo: repo, redisClient: redisClient}
}

// Query 根据文件哈希查询文件信息。
//
// 参数:
// - ctx context.Context: 上下文。
// - file *v1.QueryFile: 查询文件信息的请求。
// - info *v1.FileInfo: 查询到的文件信息。
//
// 返回值:
// - error: 查询过程中的错误，如果查询成功则为 nil。
func (t *FileService) Query(ctx context.Context, file *v1.QueryFile, info *v1.FileInfo) error {
	f, err := t.repo.Query(file.Hash)
	if err != nil {
		return err
	}

	info.Hash = f.Hash
	info.Status = v1.FileStatus(f.Status)
	info.ChunkNums = f.ChunkNums
	info.CurrentIndex = f.CurrentIndex
	info.SegmentSize = f.SegmentSize
	if f.Status == int32(v1.FileStatus_Uploaded) {
		info.Bucket = &f.Bucket
		info.FilePath = &f.FilePath
	}
	return nil
}

// MarkChunkOK 标记分块上传完成。
//
// 参数:
// - ctx context.Context: 上下文。
// - chunk *v1.MarkChunk: 标记分块上传完成的请求。
// - info *v1.FileInfo: 文件信息。
//
// 返回值:
// - error: 操作过程中的错误，如果操作成功则为 nil。
func (t *FileService) MarkChunkOK(ctx context.Context, chunk *v1.MarkChunk, info *v1.FileInfo) error {
	// 创建 Redis 锁实例。
	locker := lock.NewRedisLocker(t.redisClient, chunk.Hash, time.Second*10)
	err := locker.Lock(ctx, time.Second*60)
	if err != nil {
		return err
	}
	defer locker.UnLock(ctx)

	f, err := t.repo.Query(chunk.Hash)
	if err != nil {
		return err
	}

	if f.Status == int32(v1.FileStatus_Pending) {
		f.Status = int32(v1.FileStatus_Uploading)
	}

	f.Chunks = append(f.Chunks, Chunk{
		ChunkIndex: chunk.ChunkIndex,
		ChunkOK:    true,
	})

	// 排序，寻找最近上次成功的index
	sort.SliceStable(f.Chunks, func(i, j int) bool {
		return f.Chunks[i].ChunkIndex < f.Chunks[j].ChunkIndex
	})

	var lastIndex = int64(-1)
	for _, chk := range f.Chunks {
		if lastIndex+1 == chk.ChunkIndex {
			lastIndex = chk.ChunkIndex
		} else {
			break
		}
	}

	if lastIndex >= f.ChunkNums-1 {
		f.Status = int32(v1.FileStatus_Uploaded)
	}

	info.Hash = f.Hash
	info.Status = v1.FileStatus(f.Status)
	info.ChunkNums = f.ChunkNums
	info.CurrentIndex = f.CurrentIndex
	info.SegmentSize = f.SegmentSize
	if f.Status == int32(v1.FileStatus_Uploaded) {
		info.Bucket = &f.Bucket
		info.FilePath = &f.FilePath
	}
	return t.repo.Update(chunk.Hash, map[string]any{
		"Chunks":       f.Chunks,
		"Status":       f.Status,
		"CurrentIndex": lastIndex,
	})
}

// StartSegmentUpload 启动分段上传过程。
//
// 参数:
// - ctx context.Context: 上下文。
// - upload *v1.SegmentUpload: 启动分段上传的请求。
// - info *v1.FileInfo: 文件信息。
//
// 返回值:
// - error: 操作过程中的错误，如果操作成功则为 nil。
func (t *FileService) StartSegmentUpload(ctx context.Context, upload *v1.SegmentUpload, info *v1.FileInfo) error {
	defer func() {
		info.Hash = upload.Hash
		info.Status = v1.FileStatus_Pending
		info.ChunkNums = upload.ChunkNums
		info.CurrentIndex = -1
		info.SegmentSize = upload.SegmentSize
	}()
	return t.repo.Create(&File{
		Hash:         upload.Hash,
		Status:       int32(v1.FileStatus_Pending),
		ChunkNums:    upload.ChunkNums,
		CurrentIndex: -1,
		SegmentSize:  upload.SegmentSize,
		Bucket:       upload.Bucket,
		FilePath:     upload.FilePath,
		Chunks:       nil,
	})
}
