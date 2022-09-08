package utils

import (
	"github.com/mitchellh/go-homedir"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetAppRuntimePath 获取应用运行时所在目录
func GetAppRuntimePath() string {
	file, _ := exec.LookPath(os.Args[0])
	if strings.Contains(file, "go-build") {
		pwd, _ := os.Getwd()
		return pwd
	}

	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

// GetAppDataBaseDir 获取AppData路径
func GetAppDataBaseDir(child ...string) string {
	path, _ := homedir.Expand("~/sms-boom-go")
	if len(child) > 0 {
		return path + "/" + strings.Join(child, "/")
	}

	return path
}

// GetAppDataLogDir 获取日志目录路径
func GetAppDataLogDir(child ...string) string {
	return GetAppDataBaseDir(append([]string{"logs"}, child...)...)
}

// GetAppDataConfigDir 获取配置文件路径
func GetAppDataConfigDir(child ...string) string {
	return GetAppDataBaseDir(append([]string{"configs"}, child...)...)
}
