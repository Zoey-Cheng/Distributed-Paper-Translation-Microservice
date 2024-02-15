package translation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	v1 "paper-translation/api/translation/service/v1"
	"paper-translation/pkg/signal"
	xfspark "paper-translation/pkg/xf-spark"
	"strings"
	"time"
)

const (
	Prompt = "帮我翻译下面这段文字为%s\n%s"
)

type TranslationStatus struct {
	TranslatedText string
	Finished       bool
}

func (t *TranslationStatus) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, t)
}

func (t TranslationStatus) MarshalBinary() (data []byte, err error) {
	return json.Marshal(t)
}

type TranslationService struct {
	xfSparkClient *xfspark.XFSparkClient
	signalFactory signal.SignalFactory
	redisClient   *redis.Client
}

func NewTranslationService(xfSparkClient *xfspark.XFSparkClient, signalFactory signal.SignalFactory, redisClient *redis.Client) *TranslationService {
	return &TranslationService{xfSparkClient: xfSparkClient, signalFactory: signalFactory, redisClient: redisClient}
}

func (t *TranslationService) Translate(ctx context.Context, req *v1.Translation, resp *v1.TranslationID) error {

	//分句
	sentences := strings.FieldsFunc(req.Text, func(r rune) bool {
		for _, rx := range []rune{'.', '。', '?'} {
			if rx == r {
				return true
			}
		}
		return false
	})

	// 按照token估算分段
	var buf bytes.Buffer
	var segments = make([]string, 0)
	for i := range sentences {
		buf.WriteString(sentences[i])
		if xfspark.WordCount(buf.String()) > 2000 { //最多8000token
			segments = append(segments, buf.String())
			buf.Reset()
		}
	}

	// 添加剩余部分
	if buf.Len() > 0 {
		segments = append(segments, buf.String())
	}

	resp.TaskId = uuid.NewString()
	t.redisClient.Set(ctx, resp.TaskId, TranslationStatus{TranslatedText: "", Finished: false}, time.Hour)
	go func() {
		err := t.StartPipeline(context.TODO(), resp.TaskId, segments, req.TargetLanguage)
		if err != nil {
			log.Printf("exec translate pipeline err: %+v", err)
			return
		}
	}()
	return nil
}

func (t *TranslationService) GetStatus(ctx context.Context, req *v1.TranslationID, resp *v1.TranslatedText) error {

	var status TranslationStatus
	err := t.redisClient.Get(ctx, req.TaskId).Scan(&status)
	if err != nil {
		return err
	}
	resp.Text = status.TranslatedText
	resp.Finished = status.Finished
	return nil
}

func (t *TranslationService) StartPipeline(ctx context.Context, taskID string, segments []string, language string) error {
	semaphore := t.signalFactory.Semaphore("xf-spark", 2)
	ticker := time.NewTicker(time.Millisecond * 500)
	timer := time.NewTimer(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			acquire, err := semaphore.Acquire()
			if err != nil {
				return err
			}
			if acquire {
				goto Translate
			}
		case <-timer.C:
			return errors.New("wait to translate timeout")
		}
	}

Translate:
	defer func() {
		err := semaphore.Release()
		if err != nil {
			log.Printf("release semaphore err: %v", err)
		}
	}()
	log.Printf("begin translate text: %+v", segments)
	//讯飞只给2并发
	var translatedText bytes.Buffer
	defer func() {
		log.Printf("translate result: %s", translatedText.String())
		t.redisClient.Set(ctx, taskID, TranslationStatus{TranslatedText: translatedText.String(), Finished: true}, time.Hour)
	}()
	for _, segment := range segments {
		err := t.xfSparkClient.CreateChat(ctx, fmt.Sprintf(Prompt, language, segment), func(text string) {
			translatedText.WriteString(text)
		})
		if err != nil {
			return err
		}
	}
	return nil
}
