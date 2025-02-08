package info

type InfoLevel string

const (
	InfoLevelInfo      InfoLevel = "info"
	InfoLevelTerminate InfoLevel = "terminate"
	InfoLevelError     InfoLevel = "error"
	InfoLevelFatal     InfoLevel = "fatal"
)
