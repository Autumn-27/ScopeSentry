// utils-------------------------------------
// @file      : filetool.go
// @author    : Autumn
// @contact   : rainy-autumn@outlook.com
// @time      : 2025/5/25 15:54
// -------------------------------------------

package utils

import (
	"fmt"
	"os"
)

func EnsureDir(path string) error {
	// 检查路径状态
	info, err := os.Stat(path)

	// 如果路径不存在
	if os.IsNotExist(err) {
		// 自动创建多层目录
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		return nil
	}

	// 路径存在，但不是目录
	if err == nil && !info.IsDir() {
		return fmt.Errorf("路径已存在但不是目录: %s", path)
	}

	// 已存在且是目录
	return nil
}

func WriteFile(filePath string, content []byte) error {
	// os.Create 会清空已存在的文件内容，相当于覆盖写
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入内容
	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}
