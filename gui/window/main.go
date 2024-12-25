package window

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"sms-boom-go/gui/layouts"
	"strconv"
	"strings"
	"time"
)

type MainWindow struct {
	app fyne.App
	win fyne.Window
}

type FormData struct {
	Phones   []string
	Rounds   int
	Interval int
}

// NewMainWindow 创建并初始化主窗口
func NewMainWindow(app fyne.App) *MainWindow {
	win := app.NewWindow("SMS Boom Go 短信轰炸器")
	win.Resize(fyne.NewSize(700, 500))
	win.CenterOnScreen()

	mainWin := &MainWindow{
		app: app,
		win: win,
	}

	win.SetContent(mainWin.buildContent())

	return mainWin
}

func (w *MainWindow) ShowAndRun() {
	w.win.ShowAndRun()
}

// buildContent 构建窗口内容
func (w *MainWindow) buildContent() fyne.CanvasObject {
	// 表单容器上下的间距
	rect := canvas.NewRectangle(&color.NRGBA{R: 0, G: 0, B: 0, A: 0})
	rect.SetMinSize(fyne.NewSize(10, 10))

	// 日志区域
	logMsg := ""
	logArea := widget.NewMultiLineEntry()

	//logArea.Disable() // Disable 方法会影响字体颜色，视觉效果不好
	logArea.OnChanged = func(s string) { // 使用 OnChanged 事件模拟 Disable 的禁止输入效果
		if logMsg != s {
			logArea.SetText(logMsg)
		}
	}

	// 日志函数
	logMessage := func(msg string) {
		log.Println(msg)
		logMsg = msg + "\n" + logArea.Text
		logArea.SetText(logMsg)
	}
	logMessage("欢迎使用 SMS Boom Go 短信轰炸器\n")

	// 进度条
	progressBar := widget.NewProgressBar()
	stopBut := widget.NewButton("取消", func() {
		dialog.ShowConfirm("请确认", "确定要取消轰炸吗？", func(b bool) {
			log.Printf("确定要取消吗? %v\n", b)
		}, w.win)
	})

	progressBarCon := container.NewBorder(rect, nil, nil, nil, container.New(
		&layouts.RowRatioLayout{
			Interval: 8,
			Ratio:    []int{10, 2},
		},
		progressBar,
		stopBut,
	))
	progressBarCon.Hide()

	// 表单
	form := w.buildForm(func(form *widget.Form, phoneEntry, roundsEntry, intervalEntry *widget.Entry) func() {
		return func() {
			// 手机号
			phones := strings.Split(phoneEntry.Text, "\n")
			var filteredPhones []string
			for _, v := range phones {
				if len(strings.TrimSpace(v)) > 0 {
					filteredPhones = append(filteredPhones, strings.TrimSpace(v))
				}
			}
			if len(filteredPhones) == 0 {
				dialog.ShowInformation("手机号错误", "请至少输入一个手机号", w.win)
				return
			}

			// 轰炸轮数
			rounds, err := strconv.Atoi(roundsEntry.Text)
			if err != nil || rounds < 1 {
				dialog.ShowInformation("轰炸轮数错误", "轰炸轮数只能是大于1的数字", w.win)
				return
			}

			// 轰炸间隔
			interval, err := strconv.Atoi(intervalEntry.Text)
			if err != nil || interval < 1 {
				dialog.ShowInformation("轰炸间隔错误", "轰炸间隔只能是大于1的数字", w.win)
				return
			}

			formData := &FormData{
				Phones:   filteredPhones,
				Rounds:   rounds,
				Interval: interval,
			}

			// 禁用表单、开启进度条
			form.Disable()
			progressBarCon.Show()

			logMessage(fmt.Sprintf("轰炸轮数: %d\n轰炸间隔: %d秒\n轰炸手机号: %s\n", rounds, interval, strings.Join(formData.Phones, ",")))

			// 此处一定要开启协程处理，否则会阻塞其他组件的事件
			go func() {
				// 开始轰炸
				startBoom(form, formData, progressBar, logMessage)

				// 恢复表单、隐藏进度条
				progressBarCon.Hide()
				dialog.ShowInformation("提示", "轰炸完成 ^_^", w.win)
				form.Enable()
			}()
		}
	})

	return container.NewVSplit(
		container.NewBorder(rect, rect, nil, nil, container.NewVBox(
			form,
			progressBarCon,
		)),
		container.NewScroll(logArea),
	)
}

// buildForm 构建表单组件
func (w *MainWindow) buildForm(onSubmitFun func(from *widget.Form, phone, rounds, interval *widget.Entry) func()) fyne.CanvasObject {
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

	// 表单布局
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "手 机 号：", Widget: phoneEntry, HintText: "多个手机号时每行一个"},
			{Text: "轰炸轮数：", Widget: roundsEntry, HintText: "对每个手机号轰炸几轮，默认1轮"},
			{Text: "轰炸间隔：", Widget: intervalEntry, HintText: "每轮轰炸结束后休息几秒，默认60秒"},
		},
		SubmitText: "开始轰炸",
		CancelText: "重置",
		OnCancel: func() {
			phoneEntry.SetText("")
			roundsEntry.SetText("1")
			intervalEntry.SetText("60")
		},
	}

	form.OnSubmit = onSubmitFun(form, phoneEntry, roundsEntry, intervalEntry)

	return form
}

func startBoom(form *widget.Form, formData *FormData, progressBar *widget.ProgressBar, logMessage func(msg string)) {
	total := formData.Rounds * len(formData.Phones)

	progressBar.TextFormatter = func() string {
		return fmt.Sprintf("%d / %d", int(progressBar.Value), int(progressBar.Max))
	}
	progressBar.Max = float64(total)
	progressBar.SetValue(0)

	for i := 0; i < formData.Rounds; i++ {
		for index, phone := range formData.Phones {
			logMsg := fmt.Sprintf("正在轰炸第%d轮，第%d个手机号: %s\n", i+1, index+1, phone)
			logMessage(logMsg)

			// 模拟短信轰炸
			time.Sleep(time.Second)

			//progress := float64((i*len(formData.Phones))+index+1) / float64(total)
			progressBar.SetValue(float64((i * len(formData.Phones)) + index + 1))
		}
	}
}
