package layouts

import (
	"fyne.io/fyne/v2"
)

// RowRatioLayout 行比例布局
type RowRatioLayout struct {
	Interval float32 // 每个组件的间距
	Ratio    []int   // 每个组件的比例，所有组件的比例之和必须为12
}

func (r *RowRatioLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) <= 0 {
		panic("objects is empty")
	}
	if len(objects) != len(r.Ratio) {
		panic("objects and ratio length is not equal")
	}

	// 计算可用宽度
	availWidth := containerSize.Width - (float32(len(objects)-1) * r.Interval)

	// 遍历组件，计算每个组件的位置和尺寸
	var sumWidth float32
	for i, obj := range objects {
		width := (availWidth * float32(r.Ratio[i])) / 12

		obj.Resize(fyne.NewSize(width, containerSize.Height))
		obj.Move(fyne.NewPos(sumWidth, 0))

		sumWidth += width + r.Interval
	}
}

func (r *RowRatioLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	width := float32(0)
	height := float32(0)
	for _, obj := range objects {
		minSize := obj.MinSize()
		width += minSize.Width
		height = fyne.Max(height, minSize.Height)
	}
	return fyne.NewSize(width, height)

}
