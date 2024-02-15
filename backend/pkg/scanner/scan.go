package scanner

import (
	"bufio"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	ScanLines     = bufio.ScanLines
	ScanWords     = bufio.ScanWords
	ScanRunes     = bufio.ScanRunes
	ScanSentences = scanSentences
)

const abbreviationLength = 4

// isEndOfSentence 根据一些启发式方法来判断是否是句子的结尾。
func isEndOfSentence(lastSymbol rune, prev, next []byte) bool {
	if len(prev) == 0 {
		return false
	}

	switch lastSymbol {
	case '\r', '\n', '!', '?', '。':
		return true
	case '.':
		// 不幸的是，句点并不是句子结尾的明确标志。
		// 它可以用于姓名或其他缩写。

		// 检查最后的符号。
		// 如果它们是 \w{abbreviationLength} – 看起来像 Mrs.、Dr. 或其他缩写。
		// 如果前后的符号是数字 – 它是一个浮点数。
		// 否则，它看起来像是句子的结尾。
		prevRune, width := utf8.DecodeLastRune(prev)
		if unicode.IsLetter(prevRune) {
			nextAfterCurrent := prevRune // 需要检查第一个空格之后的字母大写。
			margin := width

			// 让我们找到空格。
			for i := 0; i < abbreviationLength; i++ {
				if len(prev)-margin <= 0 {
					return false
				}

				currentRune, currentWidth := utf8.DecodeLastRune(prev[:len(prev)-margin])

				// 我们找到了空格。让我们检查空格之后的下一个字符是否是大写字母。
				if unicode.IsSpace(currentRune) {
					return !unicode.IsUpper(nextAfterCurrent)
				}

				// 符号是否在某个组中？
				if strings.ContainsAny(string(currentRune), `[({"'`) {
					return false
				}

				// 如果不是字母，它不是缩写。
				if !unicode.IsLetter(currentRune) {
					return true
				}

				nextAfterCurrent = currentRune
				margin += currentWidth
			}

			// 如果最后的 n 个字符是字母，它很可能是字符串的结尾。
			return true
		}

		if unicode.IsDigit(prevRune) {
			if len(next) == 0 {
				return true
			}

			nextRune, _ := utf8.DecodeRune(next)
			// 看起来是浮点数。
			if unicode.IsDigit(nextRune) {
				return false
			}
		}

		return true
	}

	return false
}

// scanSentences 函数将数据拆分为句子。
func scanSentences(data []byte, atEOF bool) (advance int, token []byte, err error) {
	start := 0

	// 跳过前导空格。
	for pos, symbol := range string(data) {
		if !unicode.IsSpace(symbol) {
			break
		}

		start = pos
	}

	// 扫描直到EOF、EOL或 .!? 符号。
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])

		if isEndOfSentence(r, data[start:i], data[i:]) {
			return i + width, data[start : i+width], nil
		}
	}

	// 如果我们在EOF处，有一个最终的、非空的、非终结的句子。返回它。
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}

	// 请求更多的数据。
	return start, nil, nil
}
