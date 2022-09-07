package boom

func UpdateApi() error {
	// 从 github 获取最新接口
	GETAPIJsonUrl := "https://hk1.monika.love/OpenEthan/SMSBoom/master/GETAPI.json"
	APIJsonUrl := "https://hk1.monika.love/OpenEthan/SMSBoom/master/api.json"

	_, body, err := HttpGet(GETAPIJsonUrl)
	if err != nil {
		return err
	}

	appPath := GetAppPath()
	err = WriteFileByString(appPath+"/"+GetAPI, body)
	if err != nil {
		return err
	}

	_, body, err = HttpGet(APIJsonUrl)
	if err != nil {
		return err
	}

	err = WriteFileByString(appPath+"/"+API, body)
	if err != nil {
		return err
	}

	return nil
}
