package main

import (
	"fyne.io/fyne/v2/app"
	"sms-boom-go/gui/window"
)

func main() {
	myApp := app.New()
	mainWin := window.NewMainWindow(myApp)
	mainWin.ShowAndRun()
}
