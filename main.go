package main

import (
	"bufio"
	"fmt"
	"github.com/gookit/color"
	"github.com/urfave/cli"
	"log"
	"os"
	"sms-boom-go/boom"
	"sms-boom-go/scans"
	"sms-boom-go/utils"
	"strings"
)

var app *cli.App

func init() {
	// 创建AppData目录
	_, _ = utils.DirExistsOrCreate(utils.GetAppDataConfigDir())
	_, _ = utils.DirExistsOrCreate(utils.GetAppDataLogDir())

	// 初始化log
	utils.InitLog()
}

func main() {
	welcome()
	chooseAction()

	app = cli.NewApp()
	app.Name = "sms-boom-go"
	app.Usage = "一个健壮、免费、适用于小白的短信轰炸机"

	app.Action = action

	app.Commands = []cli.Command{
		{
			Name:    "boom",
			Aliases: []string{"1"},
			Usage:   "轰炸指定手机号",
			Action: func(c *cli.Context) error {
				scanner := bufio.NewScanner(os.Stdin)

				phone := scans.ScanHandelByStringSlice(scanner, scans.ScanPhone)
				frequency := scans.ScanHandelByInt(scanner, scans.ScanFrequency)
				interval := scans.ScanHandelByInt(scanner, scans.ScanInterval)
				coroutineCount := scans.ScanHandelByInt(scanner, scans.ScanCoroutine)

				if scans.Confirm(scanner, phone, frequency, interval, coroutineCount) {
					err := boom.Start(phone, frequency, interval, coroutineCount)
					if err != nil {
						color.Error.Printf("轰炸失败: %s!\n", err.Error())
					}
				}

				return nil
			},
		},
		{
			Name:    "update",
			Aliases: []string{"2"},
			Usage:   "更新API接口",
			Action: func(c *cli.Context) error {
				color.Info.Println("正在从GitHub拉取最新接口...")
				err := boom.UpdateApi()
				if err != nil {
					color.Error.Printf("接口保存失败: %s, 请关闭所有代理软件多尝试几次!\n", err.Error())
				} else {
					color.Success.Println("API接口已更新")
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

// welcome 输出欢迎语
func welcome() {
	logo := "   _____ __  __  _____    ____                            _____       \n  / ____|  \\/  |/ ____|  |  _ \\                          / ____|      \n | (___ | \\  / | (___    | |_) | ___   ___  _ __ ___    | |  __  ___  \n  \\___ \\| |\\/| |\\___ \\   |  _ < / _ \\ / _ \\| '_ ` _ \\   | | |_ |/ _ \\ \n  ____) | |  | |____) |  | |_) | (_) | (_) | | | | | |  | |__| | (_) |\n |_____/|_|  |_|_____/   |____/ \\___/ \\___/|_| |_| |_|   \\_____|\\___/ \n                                                                      \n                                                                      "
	main := `
本项目为参考SMSBoom实现的适用于小白的简易版本
项目地址：https://github.com/FrankTiao/sms-boom-go
SMSBoom：https://github.com/OpenEthan/SMSBoom

免责声明：
1. 本程序仅供娱乐, 源码全部开源, 禁止滥用和二次贩卖盈利. 禁止用于商业用途.
2. 若使用者滥用本项目, 作者无需承担任何法律责任.
`
	color.Green.Println(logo)
	color.White.Println(main)
}

// chooseAction 可选的操作列表
func chooseAction() {
	main := `
功能列表：（输入 exit 或按 Ctrl+C 退出）
=============================================
 1. 轰炸指定手机号      2. 更新API接口
=============================================
`
	color.Green.Println(main)
}

func action(c *cli.Context) error {
	if c.NArg() != 0 {
		color.Error.Printf("\n\n未找到要执行的操作: %s，请重新选择\n\n", c.Args().Get(0))
		return nil
	}

L:
	for {
		fmt.Print("请输入数字选择要执行的操作: ")

		//　读取用户输入
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "exit":
			fmt.Println("exit...")
			break L
		}

		cmdArgs := strings.Split(input, " ")
		if len(cmdArgs) == 0 {
			continue
		}

		s := []string{app.Name}
		s = append(s, cmdArgs...)

		err := c.App.Run(s)
		if err != nil {
			color.Error.Print("\n\n应用发生错误，请重新选择\n\n")
		}

		// 显示选择项
		chooseAction()
	}

	return nil
}
