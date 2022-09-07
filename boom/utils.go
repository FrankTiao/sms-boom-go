package boom

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	if strings.Contains(file, "go-build") {
		pwd, _ := os.Getwd()
		return pwd
	}

	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index]
}

// HttpGet 发起HTTP GET请求
func HttpGet(url string) (int, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}

func WriteFileByString(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writeIO := bufio.NewWriter(file)
	_, err = writeIO.WriteString(content)
	if err != nil {
		return err
	}

	err = writeIO.Flush()
	if err != nil {
		return err
	}

	return nil
}

// FileExists 判断所给路径文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
