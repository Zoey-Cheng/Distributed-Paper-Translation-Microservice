package xf_spark

import (
	"bufio"
	"strings"
)

/**
 * WordCount 统计文本中的单词数量。
 *
 * 参数:
 * - text (string): 要统计的文本。
 *
 * 返回值:
 * - count (int): 单词数量。
 */
func WordCount(text string) (count int) {
	scanner := bufio.NewScanner(strings.NewReader(text)) // 创建一个文本扫描器。
	scanner.Split(bufio.ScanWords)                       // 将文本扫描器配置为按单词拆分文本。
	for scanner.Scan() {
		count++
	}
	return count
}
