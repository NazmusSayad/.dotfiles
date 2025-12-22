package constants

type WindowsTaskMode string

const (
	WindowsTaskModeUser  WindowsTaskMode = "User"
	WindowsTaskModeAdmin WindowsTaskMode = "Admin"
)

type WindowsTask struct {
	Mode     WindowsTaskMode
	Program  string
	TaskName string
}

var STARTUP_TASK_SLACK = WindowsTask{
	Mode:     WindowsTaskModeUser,
	Program:  "slack-startup",
	TaskName: "__SLACK STARTUP TASK___",
}

var STARTUP_TASK_WINDOWS = WindowsTask{
	Mode:     WindowsTaskModeAdmin,
	Program:  "windows-startup",
	TaskName: "__WINDOWS STARTUP TASK__",
}
