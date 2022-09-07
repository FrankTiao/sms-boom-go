package utils

import (
	"os"
)

// PathExists 文件/文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// WriteFile 写入文件
func WriteFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0711)
}

// DirExistsOrCreate 目录是否存在，不存在则创建
func DirExistsOrCreate(path string) (bool, error) {
	if exists := PathExists(path); exists {
		return true, nil
	}

	err := os.MkdirAll(path, 0711)
	if err != nil {
		return false, err
	}

	return true, nil
}
