package pdf_test

import (
	"fmt"
	"paper-translation/pkg/pdf"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
 * TestConvertPdfToImages 是对 ConvertPdfToImages 函数进行测试的测试用例。
 * 它测试将指定的PDF文件转换为图像文件，并检查是否没有错误发生。
 */
func TestConvertPdfToImages(t *testing.T) {
	// 调用 ConvertPdfToImages 函数来将测试PDF文件转换为图像文件
	images, clean, err := pdf.ConvertPdfToImages("./test.pdf")

	// 使用 testify 断言库检查是否没有错误
	assert.NoError(t, err)

	// 打印图像文件路径切片
	fmt.Println(images)

	// 调用清理函数以删除临时文件和目录
	clean()
}
