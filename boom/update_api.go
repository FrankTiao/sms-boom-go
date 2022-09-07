package boom

import (
	"sms-boom-go/configs"
	"sms-boom-go/utils"
)

func UpdateApi() error {
	// 从 github 获取最新接口
	GETAPIJsonUrl := "https://hk1.monika.love/OpenEthan/SMSBoom/master/GETAPI.json"
	APIJsonUrl := "https://hk1.monika.love/OpenEthan/SMSBoom/master/api.json"

	_, body, err := utils.HttpGet(GETAPIJsonUrl)
	if err != nil {
		return err
	}

	err = utils.WriteFile(utils.GetAppDataConfigDir(configs.GetAPI), body)
	if err != nil {
		return err
	}

	_, body, err = utils.HttpGet(APIJsonUrl)
	if err != nil {
		return err
	}

	err = utils.WriteFile(utils.GetAppDataConfigDir(configs.API), body)
	if err != nil {
		return err
	}

	return nil
}
