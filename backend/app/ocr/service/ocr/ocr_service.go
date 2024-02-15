package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	v1 "paper-translation/api/ocr/service/v1"
	"paper-translation/pkg/pdf"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	ocr "github.com/alibabacloud-go/ocr-api-20210707/client"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// OCRStatus 存储OCR任务的状态
type OCRStatus struct {
	Text     string
	Finished bool
}

// UnmarshalBinary 从二进制数据中反序列化OCRStatus
func (t *OCRStatus) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

// MarshalBinary 将OCRStatus序列化为二进制数据
func (t OCRStatus) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

// OCRService 包含OCR服务的实现
type OCRService struct {
	ocrRepo     OCRRepository // OCR任务的存储库
	ocr         *ocr.Client   // 阿里云OCR客户端
	oss         *oss.Client   // 阿里云OSS客户端
	redisClient *redis.Client // Redis客户端，用于存储OCR任务状态
}

// NewOCRService 创建一个新的OCRService实例
func NewOCRService(ocrRepo OCRRepository, ocr *ocr.Client, oss *oss.Client, redisClient *redis.Client) *OCRService {
	return &OCRService{ocrRepo: ocrRepo, ocr: ocr, oss: oss, redisClient: redisClient}
}

// OCR 启动OCR任务，处理文档的OCR识别
func (t *OCRService) OCR(ctx context.Context, param *v1.OCRParam, resp *v1.OCRTaskID) error {
	// 检查是否已经存在OCR结果
	ocx, err := t.ocrRepo.Get(param.Bucket, param.ObjectKey, param.FileType)
	if err == nil {
		resp.TaskId = uuid.NewString()
		t.redisClient.Set(ctx, resp.TaskId, OCRStatus{Text: ocx.OcredText, Finished: true}, time.Hour)
		return nil
	}

	// 如果存在其他错误（不是没有找到文档），则返回错误
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	// 创建新的OCR任务
	resp.TaskId = uuid.NewString()
	// 每次进行一个新的 OCR  大任务， 都要往 redis 里面存一下，记录一下这个开始的任务，
	//存到 Redis 里 key 是 taskID， value 是一个对象，字段  text 是将文件序列化后变成字符串存进去，status 就是这个 taskID 的执行状态
	t.redisClient.Set(ctx, resp.TaskId, OCRStatus{Text: "", Finished: false}, time.Hour)
	go func() {
		err = t.StartPipeline(context.TODO(), resp.TaskId, param.Bucket, param.ObjectKey)
		if err != nil {
			log.Printf("exec ocr pipeline failed err: %+v", err)
		}
	}()
	return nil
}

// GetStatus 获取OCR任务的状态
// 就是按照 taskID 去查看这个任务状态
func (t *OCRService) GetStatus(ctx context.Context, req *v1.OCRTaskID, resp *v1.OCRText) error {
	var status OCRStatus
	err := t.redisClient.Get(ctx, req.TaskId).Scan(&status)
	if err != nil {
		return err
	}
	resp.Text = status.Text
	resp.Finished = status.Finished
	return nil
}

// OCRLocalImage 对本地图像执行OCR识别
func (t *OCRService) OCRLocalImage(bucket, filePath string) (string, error) {
	// 获取OSS存储桶
	bkt, err := t.oss.Bucket(bucket)
	if err != nil {
		return "", err
	}

	// 打开本地图像文件
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	// 生成随机的对象键，将图像上传到OSS，因为 OCR 接口只能传 url 进去
	objectKey := fmt.Sprintf("images/%s.jpg", uuid.NewString())
	err = bkt.PutObject(objectKey, f)
	if err != nil {
		return "", err
	}

	// 生成带签名的URL以下载图像
	fileURL, err := bkt.SignURL(objectKey, http.MethodGet, 120)
	if err != nil {
		return "", err
	}

	log.Printf("start ocr for image %s", fileURL)

	// 创建OCR请求并执行OCR
	req := &ocr.RecognizeEnglishRequest{Url: &fileURL}
	resp, err := t.ocr.RecognizeEnglish(req)
	if err != nil {
		return "", err
	}

	// 解析OCR响应数据
	var data map[string]any
	err = json.Unmarshal([]byte(*resp.Body.Data), &data)
	if err != nil {
		return "", err
	}

	return data["content"].(string), nil
}

// ConvertLocalImages 将本地PDF文件转换为图像
// 因为接口不支持直接输入多页 PDF ，所以用命令把 PDF 拆分成图片，一个个喂
func (t *OCRService) ConvertLocalImages(bucket, filePath string) ([]string, func(), error) {
	// 获取OSS存储桶
	bkt, err := t.oss.Bucket(bucket)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("start ocr for object: %s", filePath)

	// 获取OSS对象
	object, err := bkt.GetObject(filePath)
	if err != nil {
		log.Printf("get object %s/%s err: %+v", bucket, filePath, err)
		return nil, nil, err
	}
	defer object.Close()

	// 到这里我们已经从 oss 里面拿到了需要处理的 PDF

	// 生成本地临时PDF文件并将对象内容复制到该文件，我们在对生产环境的任何文件进行更改的时候，都要复制一下去操作副本
	localFilePath := fmt.Sprintf("%s/%s.pdf", os.TempDir(), uuid.NewString())
	file, err := os.OpenFile(localFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	_, _ = io.Copy(file, object)
	_ = file.Close()
	defer func() {
		_ = os.RemoveAll(localFilePath)
	}()

	// 然后将PDF文件转换为图像
	return pdf.ConvertPdfToImages(localFilePath)
}

// StartPipeline 启动OCR处理管道，包括图像转换和OCR识别
// 是总的流水线函数，对一个 PDF 做 OCR

func (t *OCRService) StartPipeline(ctx context.Context, taskID, bucket, filePath string) error {
	// 将本地PDF文件转换为图像
	images, clean, err := t.ConvertLocalImages(bucket, filePath)
	if err != nil {
		return err
	}
	defer clean()

	log.Printf("convert images is %+v", images)
	var wg sync.WaitGroup
	var texts = make([]string, len(images)) //这里先记录一下顺序，免得并发执行后 OCR 的文本顺序混乱
	for index, imagePath := range images {
		wg.Add(1)
		go func(index int, imagePath string) { //并发执行图片的 OCR，调接口同时进行
			defer wg.Done()
			// 对每个图像执行OCR识别
			text, err := t.OCRLocalImage(bucket, imagePath)
			if err != nil {
				log.Printf("ocr err: %+v", err)
				return
			}
			texts[index] = text
		}(index, imagePath)
	}
	wg.Wait() //等待并发任务全部结束

	// 将OCR识别的文本合并成一个文本，按照刚才记录的顺序
	var buf bytes.Buffer
	for i := range texts {
		buf.WriteString(texts[i])
	}

	// 将OCR任务的状态标记为已完成，并存储OCR结果到 Redis
	t.redisClient.Set(ctx, taskID, OCRStatus{Text: buf.String(), Finished: true}, time.Hour)
	if buf.String() != "" {
		_ = t.ocrRepo.Create(&OCR{
			ID:        taskID,
			Bucket:    bucket,
			ObjectKey: filePath,
			OcredText: buf.String(),
		})
	}
	return nil
}
