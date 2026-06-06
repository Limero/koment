package app

import (
	"fmt"

	"github.com/limero/koment/app/info"
)

func (a *App) setInfo(level info.InfoLevel, msg string) {
	a.infoLevel = level
	a.infoMsg = msg
}

// Info will show an error message and continue
func (a *App) Info(format string, v ...any) {
	a.mu.Lock()
	a.setInfo(info.InfoLevelInfo, fmt.Sprintf(format, v...))
	a.mu.Unlock()
}

// Terminate will show an info message and quit the program
func (a *App) Terminate(format string, v ...any) {
	a.mu.Lock()
	a.setInfo(info.InfoLevelTerminate, fmt.Sprintf(format, v...))
	a.mu.Unlock()
}

// Error will show an error message and continue
func (a *App) Error(format string, v ...any) {
	a.mu.Lock()
	a.setInfo(info.InfoLevelError, fmt.Sprintf("Error: "+format, v...))
	a.mu.Unlock()
}

// Fatal will show an error message and quit the program
func (a *App) Fatal(format string, v ...any) {
	a.mu.Lock()
	a.setInfo(info.InfoLevelFatal, fmt.Sprintf("Error: "+format, v...))
	a.mu.Unlock()
}
