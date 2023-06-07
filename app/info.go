package app

import "fmt"

func (a *App) Info(format string, v ...any) {
	a.infoLevel = "info"
	a.infoMsg = fmt.Sprintf(format, v...)
}

func (a *App) Error(format string, v ...any) {
	a.infoLevel = "error"
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}

func (a *App) Fatal(format string, v ...any) {
	a.infoLevel = "fatal"
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}
