package app

import "fmt"

type InfoLevel string

const (
	InfoLevelInfo      InfoLevel = "info"
	InfoLevelTerminate InfoLevel = "terminate"
	InfoLevelError     InfoLevel = "error"
	InfoLevelFatal     InfoLevel = "fatal"
)

// Info will show an error message and continue
func (a *App) Info(format string, v ...any) {
	a.infoLevel = InfoLevelInfo
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Terminate will show an info message and quit the program
func (a *App) Terminate(format string, v ...any) {
	a.infoLevel = InfoLevelTerminate
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Error will show an error message and continue
func (a *App) Error(format string, v ...any) {
	a.infoLevel = InfoLevelError
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}

// Fatal will show an error message and quit the program
func (a *App) Fatal(format string, v ...any) {
	a.infoLevel = InfoLevelFatal
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}
