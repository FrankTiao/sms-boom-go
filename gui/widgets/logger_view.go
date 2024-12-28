package widgets

import (
	"bufio"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"os"
)

type LoggerView struct {
	widget.BaseWidget

	OpenFolderHandler func()
	ClearHandler      func()

	w io.Writer
	r io.Reader
}

func NewLoggerView(openFolderHandler, clearHandler func()) *LoggerView {
	lv := &LoggerView{
		OpenFolderHandler: openFolderHandler,
		ClearHandler:      clearHandler,
	}
	lv.ExtendBaseWidget(lv)
	return lv
}

func (lv *LoggerView) CreateRenderer() fyne.WidgetRenderer {
	th := lv.Theme()
	v := fyne.CurrentApp().Settings().ThemeVariant()
	lv.ExtendBaseWidget(lv)

	box := canvas.NewRectangle(th.Color(theme.ColorNameInputBackground, v))
	box.CornerRadius = th.Size(theme.SizeNameInputRadius)
	box.StrokeWidth = 0.0

	border := canvas.NewRectangle(color.Transparent)
	border.SetMinSize(fyne.NewSize(0.5, 0.5))

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.FolderOpenIcon(), lv.OpenFolderHandler),
		widget.NewToolbarAction(theme.DeleteIcon(), lv.ClearHandler),
	)

	//text := canvas.NewText("Text Object", color.White)

	text := widget.NewLabel("aaa\nbbb\nccc")
	//tc := container.NewWithoutLayout()
	//
	//go func() {
	//	i := 0
	//	for {
	//		i++
	//		t := strconv.Itoa(i) + ": aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n"
	//		//text.SetText(t + text.Text)
	//		tc.Add(widget.NewLabel(t))
	//		time.Sleep(time.Millisecond * 300)
	//	}
	//}()

	//text := canvas.NewText("Text Object", color.White)

	//filePath := "/Users/tiaooo/projects/sms-boom-go/large_text.txt" // 替换成你的大文本文件路径

	c := container.NewBorder(
		border, border, border, border,
		container.NewStack(
			box,
			container.NewBorder(toolbar, nil, nil, nil),
			container.NewScroll(text),
		),
	)
	return widget.NewSimpleRenderer(c)
}

func createLargeTextFile(filePath string, lineCount int) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for i := 0; i < lineCount; i++ {
		_, err = writer.WriteString(fmt.Sprintf("%d==aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n", i+1))
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()
}
