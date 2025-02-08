package app

import (
	"fmt"

	"github.com/limero/koment/app/info"
)

// Info will show an error message and continue
func (a *App) Info(format string, v ...any) {
	a.infoLevel = info.InfoLevelInfo
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Terminate will show an info message and quit the program
func (a *App) Terminate(format string, v ...any) {
	a.infoLevel = info.InfoLevelTerminate
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Error will show an error message and continue
func (a *App) Error(format string, v ...any) {
	a.infoLevel = info.InfoLevelError
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}

// Fatal will show an error message and quit the program
func (a *App) Fatal(format string, v ...any) {
	a.infoLevel = info.InfoLevelFatal
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}
