package paper

import (
	"context"
	"errors"
	"go-micro.dev/v4/client"
	"log"
	es "paper-translation/api/email/service/v1"
	fs "paper-translation/api/file/service/v1"
	os "paper-translation/api/ocr/service/v1"
	v1 "paper-translation/api/paper/service/v1"
	ts "paper-translation/api/translation/service/v1"
	"time"

	"github.com/google/uuid"
)

type PaperService struct {
	repo             PaperRepository
	fileService      fs.FileService
	ocrService       os.OCRService
	translateService ts.TranslationService
	emailService     es.EmailService
}

func NewPaperService(
	repo PaperRepository,
	fileService fs.FileService,
	ocrService os.OCRService,
	translateService ts.TranslationService,
	emailService es.EmailService,
) *PaperService {
	return &PaperService{
		repo:             repo,
		fileService:      fileService,
		ocrService:       ocrService,
		translateService: translateService,
		emailService:     emailService,
	}
}

func (t *PaperService) Create(ctx context.Context, req *v1.CreatePaper, resp *v1.Paper) error {

	fileInfo, err := t.fileService.Query(ctx, &fs.QueryFile{Hash: req.PaperFileHash})
	if err != nil {
		return err
	}

	if fileInfo.Status != fs.FileStatus_Uploaded {
		return errors.New("file is not uploaded")
	}

	paper := Paper{
		ID:             uuid.NewString(),
		FileHash:       req.PaperFileHash,
		CreateAt:       time.Now(),
		Status:         0,
		EmailTo:        req.EmailTo,
		TargetLanguage: req.TargetLanguage,
	}

	go func() {
		err = t.StartPipeline(context.TODO(), paper.ID, *fileInfo.Bucket, *fileInfo.FilePath, req.TargetLanguage, req.EmailTo)
		if err != nil {
			log.Printf("exec paper pipeline failed err: %+v", err)
		}
	}()
	t.ConvertPaper(&paper, resp)
	return t.repo.Create(&paper)
}

func (t *PaperService) OCR(ctx context.Context, bucket, filePath string) (string, error) {
	ocrID, err := t.ocrService.OCR(
		ctx,
		&os.OCRParam{
			Bucket:    bucket,
			ObjectKey: filePath,
		},
		client.WithDialTimeout(time.Second*300),
		client.WithRequestTimeout(time.Second*300),
	)
	if err != nil {
		log.Printf("do ocr err: %+v", err)
		return "", err
	}

	for {
		status, err := t.ocrService.GetStatus(ctx, ocrID)
		if err != nil {
			return "", err
		}
		if status.Finished {
			if status.Text == "" {
				return "", errors.New("ocr failed")
			}
			return status.Text, nil
		}
		time.Sleep(time.Second)
	}
}

func (t *PaperService) Translate(ctx context.Context, text, targetLanguage string) (string, error) {
	translateID, err := t.translateService.Translate(
		ctx,
		&ts.Translation{Text: text, TargetLanguage: targetLanguage},
		client.WithDialTimeout(time.Second*300),
		client.WithRequestTimeout(time.Second*300),
	)
	if err != nil {
		log.Printf("do translate text err: %+v", err)
		return "", err
	}

	for {
		status, err := t.translateService.GetStatus(ctx, translateID)
		if err != nil {
			return "", err
		}

		if status.Finished {
			if status.Text == "" {
				return "", errors.New("translate failed")
			}
			return status.Text, nil
		}
		time.Sleep(time.Second)
	}
}

func (t *PaperService) StartPipeline(ctx context.Context, id, bucket, filePath, targetLanguage, emailTo string) (err error) {

	defer func() {
		if err != nil {
			_ = t.repo.SetStatus(id, int32(v1.Paper_failed))
		} else {
			_ = t.repo.SetStatus(id, int32(v1.Paper_finished))
		}
	}()

	text, err := t.OCR(ctx, bucket, filePath)
	if err != nil {
		return err
	}

	_ = t.repo.SetStatus(id, int32(v1.Paper_translation))
	translate, err := t.Translate(ctx, text, targetLanguage)
	if err != nil {
		return err
	}

	if emailTo != "" {
		_, _ = t.emailService.SendEmail(ctx, &es.SendEmailParam{
			EmailTo:  emailTo,
			Subject:  "你的paper翻译完成",
			Template: "{{.Text}}",
			Vars: map[string]string{
				"Text": translate,
			},
		})
	}
	return t.repo.UpdateText(id, translate)
}

func (t *PaperService) Fetch(ctx context.Context, id *v1.PaperID, resp *v1.Paper) error {
	paper, err := t.repo.Get(id.Id)
	if err != nil {
		return err
	}
	t.ConvertPaper(paper, resp)
	return nil
}

func (t *PaperService) Delete(ctx context.Context, id *v1.PaperID, re *v1.DeletePaper) error {
	return t.repo.Delete(id.Id)
}

func (t *PaperService) Fetchs(ctx context.Context, req *v1.ReqFetchs, resp *v1.RespFetchs) error {
	papers, err := t.repo.GetPapers()
	if err != nil {
		return nil
	}
	resp.Total = int32(len(papers))
	resp.Papers = func() (res []*v1.Paper) {
		for i := 0; i < len(papers); i++ {
			res = append(res, &v1.Paper{
				Id:             papers[i].ID,
				FileHash:       papers[i].FileHash,
				CreateAt:       papers[i].CreateAt.Unix(),
				Status:         v1.Paper_Status(papers[i].Status),
				TargetLanguage: papers[i].TargetLanguage,
			})
		}
		return res
	}()
	return nil
}

func (t *PaperService) ConvertPaper(paper *Paper, resp *v1.Paper) {
	resp.Id = paper.ID
	resp.Status = v1.Paper_ocr
	resp.FileHash = paper.FileHash
	resp.CreateAt = paper.CreateAt.Unix()
	resp.TargetLanguage = paper.TargetLanguage
	resp.ResultText = paper.ResultText
}
