package pdf

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

/**
 * ConvertPdfToImages 函数接受一个PDF文件的路径作为输入，并将其转换为多个图像文件。
 * 它返回一个图像文件路径的切片，以及一个清理临时文件的函数和可能的错误。
 *
 * @param inputFile - 输入的PDF文件路径
 * @return 图像文件路径的切片、清理函数和可能的错误
 */
func ConvertPdfToImages(inputFile string) ([]string, func(), error) {
	// 创建一个唯一的临时目录，用于存储转换后的图像文件
	dirPath := fmt.Sprintf("%s%s", os.TempDir(), uuid.NewString())
	_ = os.MkdirAll(dirPath, os.ModePerm)

	// 使用外部命令 "convert" 将PDF文件转换为图像文件 ，就是用 Golang 去调用 shell 脚本，或者说执行 cmd 命令，比较挫但是很方便，别学我
	cmd := exec.Command("convert", "-density", "150", inputFile, "-quality", "90", fmt.Sprintf("%s/%s", dirPath, "output.jpg"))
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return nil, func() {
			os.RemoveAll(dirPath)
		}, err
	}

	var images = make([]string, 0)
	// 遍历临时目录，获取所有图像文件的路径并存储在切片中
	_ = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			images = append(images, path)
		}
		return nil
	})

	// 返回图像文件路径切片、清理临时目录，好习惯。
	// 注意这里返回的不止是一个图像数组（数组里存的 path），还有一个函数。为了在上层函数里 defer 执行清理任务
	return images, func() {
		os.RemoveAll(dirPath)
	}, nil
}
