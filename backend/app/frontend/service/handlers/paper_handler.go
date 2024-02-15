package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	v1 "paper-translation/api/paper/service/v1"
	"paper-translation/pkg/errutil"
)

type ReqCreatePaper struct {
	FileHash       string `json:"fileHash"`
	EmailTo        string `json:"emailTo"`
	TargetLanguage string `json:"targetLanguage"`
}

type PaperHandler struct {
	paperService v1.PaperService
}

func NewPaperHandler(paperService v1.PaperService) *PaperHandler {
	return &PaperHandler{paperService: paperService}
}

func (t *PaperHandler) CreatePaper(ctx *gin.Context) {
	var req ReqCreatePaper
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		errutil.ResponseError(ctx, errutil.RequestParamError, err)
		return
	}
	resp, err := t.paperService.Create(ctx, &v1.CreatePaper{
		PaperFileHash:  req.FileHash,
		EmailTo:        req.EmailTo,
		TargetLanguage: req.TargetLanguage,
	})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
	ctx.JSON(200, gin.H{
		"paperID": resp.Id,
	})
}

func (t *PaperHandler) GetPaper(ctx *gin.Context) {
	paper, err := t.paperService.Fetch(ctx, &v1.PaperID{Id: ctx.Param("id")})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
	ctx.JSON(200, gin.H{
		"paperID":    paper.Id,
		"status":     paper.Status,
		"createAt":   paper.CreateAt,
		"resultText": paper.ResultText,
		"fileHash":   paper.FileHash,
	})
}

func (t *PaperHandler) GetPapers(ctx *gin.Context) {
	fetchs, err := t.paperService.Fetchs(ctx, &v1.ReqFetchs{})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	var resp = make([]gin.H, 0)
	for _, paper := range fetchs.Papers {
		resp = append(resp, gin.H{
			"paperID":  paper.Id,
			"status":   paper.Status,
			"createAt": paper.CreateAt,
		})
	}
	ctx.JSON(200, resp)
}

func (t *PaperHandler) DeletePaper(ctx *gin.Context) {
	_, err := t.paperService.Delete(ctx, &v1.PaperID{Id: ctx.Param("id")})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
}

func (t *PaperHandler) DownloadPaperResult(ctx *gin.Context) {
	paper, err := t.paperService.Fetch(ctx, &v1.PaperID{Id: ctx.Param("id")})
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}

	localFile := fmt.Sprintf("%s.txt", uuid.NewString())
	defer os.Remove(localFile)
	err = os.WriteFile(localFile, []byte(paper.ResultText), 0644)
	if err != nil {
		errutil.ResponseError(ctx, errutil.UnknownError, err)
		return
	}
	ctx.FileAttachment(localFile, "paper.txt")
}
