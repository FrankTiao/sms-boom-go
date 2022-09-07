package boom

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"io/ioutil"
	"os"
	"sms-boom-go/configs"
	"sms-boom-go/utils"
	"strconv"
	"strings"
	"time"
)

type Api struct {
	Desc   string     `json:"desc"`
	Url    string     `json:"url"`
	Method string     `json:"method"`
	Header headerJson `json:"header,omitempty"`
	Data   dataJson   `json:"data,omitempty"`
}

// 由于json中的header和data字段类型无法确定，所以统一转换为string
type headerJson string
type dataJson string

func (h *headerJson) UnmarshalJSON(data []byte) error {
	da := headerJson(data)
	*h = da

	return nil
}
func (d *dataJson) UnmarshalJSON(data []byte) error {
	da := dataJson(data)
	*d = da

	return nil
}

// HandelApi 组装api请求头等信息
func (api *Api) HandelApi(phone string) {
	var he = make(map[string]string, 5)
	_ = json.Unmarshal([]byte(api.Header), &he)
	if val, ok := he["Referer"]; !ok || val == "" {
		he["Referer"] = api.Url
	}

	he["UserAgent"] = utils.RandomUserAgent()
	headerJson, err := json.Marshal(he)
	if err != nil {
		headerJson = []byte(nil)
	}

	api.Url = replaceData(api.Url, phone, 1)
	api.Data = dataJson(replaceData(string(api.Data), phone, 1))
	_ = json.Unmarshal([]byte(replaceData(string(headerJson), phone, 1)), &api.Header)
}

// Send 发起请求
func (api *Api) Send() (*resty.Response, error) {
	resp, err := sendRequest(api.Url, api.Method, string(api.Header), string(api.Data))
	return resp, err
}

// SendByGetApi 发起请求
func SendByGetApi(api, phone string) (*resty.Response, error) {
	url := replaceData(api, phone, 2)
	resp, err := sendRequest(url, "GET", "", "")
	return resp, err
}

// sendRequest 发起请求
func sendRequest(url, method, apiHeader, apiData string) (*resty.Response, error) {
	client := resty.New()

	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.SetContentLength(true)
	client.SetTimeout(time.Duration(3) * time.Second)

	request := client.R()

	var data map[string]string
	_ = json.Unmarshal([]byte(apiData), &data)

	if strings.ToLower(method) == "get" {
		request = request.SetQueryParams(data)
	} else {
		request = request.SetBody(data)
	}

	var he map[string]string
	_ = json.Unmarshal([]byte(apiHeader), &he)
	request = request.SetHeaders(he)

	var resp *resty.Response
	var err error
	if strings.ToLower(method) == "get" {
		resp, err = request.Get(url)
	} else {
		resp, err = request.Post(url)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// replaceData 替换占位符
func replaceData(content, phone string, apiType int) string {
	if apiType == 1 {
		if len(phone) > 0 && apiType == 1 {
			content = strings.ReplaceAll(content, "[phone]", phone)
			content = strings.ReplaceAll(content, "[timestamp]", strconv.FormatInt(time.Now().Unix(), 10))
		}
	} else {
		content = strings.ReplaceAll(content, "[phone]", phone)
		content = strings.ReplaceAll(content, "\\n", "")
		content = strings.ReplaceAll(content, "\\r", "")
	}

	return strings.ReplaceAll(content, "'", "\"")
}

// loadApi 加载 API
func loadApi() (*[]Api, error) {
	path := utils.GetAppRuntimePath() + "/" + configs.API
	if !utils.PathExists(path) {
		err := UpdateApi()
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)

	var apis []Api
	err = json.Unmarshal(content, &apis)
	if err != nil {
		return nil, err
	}

	if len(apis) <= 0 {
		return nil, errors.New("API为空，请尝试更新API接口")
	}

	return &apis, nil
}

// loadGetApi 加载GET API
func loadGetApi() (*[]string, error) {
	path := utils.GetAppRuntimePath() + "/" + configs.GetAPI
	if !utils.PathExists(path) {
		err := UpdateApi()
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)

	var apis []string
	err = json.Unmarshal(content, &apis)
	if err != nil {
		return nil, err
	}

	if len(apis) <= 0 {
		return nil, errors.New("GetAPI为空，请尝试更新API接口")
	}

	return &apis, nil
}
