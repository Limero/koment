package app

import "fmt"

// Info will show an error message and continue
func (a *App) Info(format string, v ...any) {
	a.infoLevel = "info"
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Terminate will show an info message and quit the program
func (a *App) Terminate(format string, v ...any) {
	a.infoLevel = "terminate"
	a.infoMsg = fmt.Sprintf(format, v...)
}

// Error will show an error message and continue
func (a *App) Error(format string, v ...any) {
	a.infoLevel = "error"
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}

// Fatal will show an error message and quit the program
func (a *App) Fatal(format string, v ...any) {
	a.infoLevel = "fatal"
	a.infoMsg = fmt.Sprintf("Error: "+format, v...)
}
