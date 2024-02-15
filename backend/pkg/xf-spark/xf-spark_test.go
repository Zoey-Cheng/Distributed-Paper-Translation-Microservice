package xf_spark_test

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	xf_spark "paper-translation/pkg/xf-spark"
	"testing"
)

func TestXFSparkClient_CreateChat(t *testing.T) {
	xfSparkClient := xf_spark.NewXFSparkClient(
		"93932e12",
		"MjEyNjc3ZDA0NWU4ODQ1MjFlMjM0OGUz",
		"496bfc567df22e7cf65de37f8cd2aa00",
	)

	err := xfSparkClient.CreateChat(context.TODO(), "帮我翻译下面这段文字为英语\n在Go语言中slice比数组更强大、灵活、方便，是一种轻量级的数据结构。slice是一个可变长度的序列，它存储了相似类型的元素，你不允许在同一个slice中存储不同类型的元素。", func(text string) {
		fmt.Print(text)
	})
	assert.NoError(t, err)
}
