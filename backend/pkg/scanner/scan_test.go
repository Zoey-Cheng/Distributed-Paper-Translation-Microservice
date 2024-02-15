package scanner_test

import (
	"bufio"
	"fmt"
	"io"
	"paper-translation/pkg/scanner"
	"strings"
	"testing"
	"time"
)

/**
 * TestSentenceSplitter_ScanSentences 是对 SentenceSplitter 的 ScanSentences 方法进行测试的测试用例。
 * 它测试了将文本分割成句子的功能。
 */
func TestSentenceSplitter_ScanSentences(t *testing.T) {
	// 创建一个扫描器来处理中文句子
	scan := bufio.NewScanner(strings.NewReader("你好。我是谁?"))
	scan.Split(scanner.ScanSentences)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}

	// 创建一个扫描器来处理英文句子
	scan = bufio.NewScanner(strings.NewReader("Hello.Who Am I?"))
	scan.Split(scanner.ScanSentences)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}

	// 创建一个管道，模拟一个包含多个句子的文本
	r, w := io.Pipe()
	go func() {
		str := strings.Repeat("Hello。Who Am I?", 10)
		for i := range []rune(str) {
			w.Write([]byte{str[i]})
			time.Sleep(time.Millisecond * 10)
		}
		w.Close()
	}()
	// 创建扫描器来处理多个句子的文本
	scan = bufio.NewScanner(r)
	scan.Split(scanner.ScanSentences)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}
}
