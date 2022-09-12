package boom

import (
	"errors"
	"regexp"
	"runtime"
	"sms-boom-go/utils"
	"strings"
	"time"
)

func UpdateProxy(force bool) error {
	proxyFilePath := utils.GetAppDataProxyConfigDir(time.Now().Format("2006-01-02") + "_open.txt")
	if !force {
		if utils.PathExists(proxyFilePath) {
			return nil
		}
	}

	_, body, err := utils.HttpGet("https://openproxy.space/list/http")
	if err != nil {
		return err
	}

	reg := regexp.MustCompile(`(?m)((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:\d+`)
	if reg == nil {
		return errors.New("regexp err")
	}

	sy := "\n"
	if runtime.GOOS == "windows" {
		sy = "\r\n"
	}

	proxy := strings.Join(reg.FindAllString(string(body), -1), sy)
	err = utils.WriteFile(proxyFilePath, []byte(proxy))
	if err != nil {
		return err
	}

	return nil
}
