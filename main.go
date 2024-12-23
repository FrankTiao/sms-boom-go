package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"
)

func main() {
	myApp := app.New()
	mainWin := myApp.NewWindow("短信轰炸器")
	mainWin.Resize(fyne.NewSize(700, 500))
	mainWin.CenterOnScreen()

	content := buildContent(mainWin)

	mainWin.SetContent(content)
	mainWin.ShowAndRun()
}

func buildContent(mainWin fyne.Window) fyne.CanvasObject {
	//日志区域
	logArea := widget.NewMultiLineEntry()
	logArea.Disable()
	for i := 0; i < 10; i++ {
		is := strconv.Itoa(i)
		logArea.SetText(logArea.Text + is + is + is + is + is + is + is + "\n")
		//logArea.Add(widget.NewLabel(is + is + is + is + is + is + is + "\n"))
	}

	// 日志函数
	logMessage := func(msg string) {
		log.Println(msg)
		//logArea.SetText(logArea.Text + msg + "\n")
	}

	// 进度条
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// 表单
	form := buildForm(func(phoneEntry, roundsEntry, intervalEntry *widget.Entry) func() {
		return func() {
			phones := strings.Split(phoneEntry.Text, "\n")
			var filteredPhones []string
			for _, phone := range phones {
				phone = strings.TrimSpace(phone)
				if phone != "" {
					filteredPhones = append(filteredPhones, phone)
				}
			}
			if len(filteredPhones) == 0 {
				dialog.ShowInformation("提示", "请至少输入一个手机号", mainWin)
				return
			}

			rounds, err := strconv.Atoi(roundsEntry.Text)
			if err != nil {
				dialog.ShowInformation("提示", "轮数输入错误: "+err.Error(), mainWin)
				return
			}
			log.Println("轰炸轮数:", rounds)

			interval, err := strconv.Atoi(intervalEntry.Text)
			if err != nil {
				dialog.ShowInformation("提示", "轰炸间隔输入错误: "+err.Error(), mainWin)
				return
			}
			log.Println("轰炸间隔:", interval)

			// 禁用输入框和按钮
			phoneEntry.Disable()
			roundsEntry.Disable()
			intervalEntry.Disable()

			progressBar.Show()
			progressBar.SetValue(0)

			go func() {
				total := rounds * len(filteredPhones)

				for i := 0; i < rounds; i++ {
					for index, phone := range filteredPhones {
						logMsg := fmt.Sprintf("正在轰炸第%d轮，第%d个手机号: %s\n", i+1, index+1, phone)
						logMessage(logMsg)

						// 模拟短信轰炸
						time.Sleep(time.Second)

						progress := float64((i*len(filteredPhones))+index+1) / float64(total)
						progressBar.SetValue(progress)
					}
				}

				progressBar.Hide()
				dialog.ShowInformation("提示", "轰炸完成", mainWin)

				// 恢复输入框和按钮
				phoneEntry.Enable()
				roundsEntry.Enable()
				intervalEntry.Enable()
			}()
		}
	})

	// 表单容器上下的间距
	rect := canvas.NewRectangle(&color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	rect.SetMinSize(fyne.NewSize(10, 10))

	return container.NewVSplit(
		container.NewBorder(rect, rect, nil, nil, container.NewVBox(
			form,
			progressBar,
		)),
		container.NewScroll(logArea),
	)
}

func buildForm(onSubmitFun func(phone, rounds, interval *widget.Entry) func()) fyne.CanvasObject {
	// 输入表单
	phoneEntry := widget.NewMultiLineEntry()
	phoneEntry.SetPlaceHolder("请输入要轰炸的手机号")

	// 轰炸轮数
	roundsEntry := widget.NewEntry()
	roundsEntry.SetText("1")
	roundsEntry.Validator = func(s string) error {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			return errors.New("只能输入大于1的数字")
		}
		return nil
	}

	// 每轮轰炸间隔
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText("60")
	intervalEntry.Validator = func(s string) error {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			return errors.New("只能输入大于1的数字")
		}
		return nil
	}

	//表单布局
	return &widget.Form{
		Items: []*widget.FormItem{
			{Text: "手 机 号：", Widget: phoneEntry, HintText: "多个手机号时使用一个空格分隔"},
			{Text: "轰炸轮数：", Widget: roundsEntry, HintText: "对每个手机号轰炸几轮，默认1轮"},
			{Text: "轰炸间隔：", Widget: intervalEntry, HintText: "每轮轰炸结束后休息几秒，默认60秒"},
		},
		SubmitText: "开始轰炸",
		OnSubmit:   onSubmitFun(phoneEntry, roundsEntry, intervalEntry),
		CancelText: "重置",
		OnCancel: func() {
			phoneEntry.SetText("")
			roundsEntry.SetText("1")
			intervalEntry.SetText("60")
		},
	}
}
