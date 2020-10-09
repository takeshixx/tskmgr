package manager

import (
	"os/user"

	"github.com/takeshixx/tskmgr/internal/tasks"
)

// TaskConfigFileName defines the filename
// of a task's configuration file.
var TaskConfigFileName = ".task.yml"

// TaskNotesFileName defines the filename
// of a task's notes file.
var TaskNotesFileName = "task_notes.md"

// Manager defines the main object
// that handles everything related
// to Tasks.
type Manager struct {
	Tasks  []*tasks.Task
	Config *Config
}

// Config defines how the configuration
// file .tskmgr.yml should look like.
type Config struct {
	TasksPath string `yaml:"tasksPath"`
}

// NewManager creates a new manager and parses
// the configuration file.
func NewManager() (m *Manager, err error) {
	m = &Manager{}
	usr, err := user.Current()
	if err != nil {
		return
	}
	if err = m.LoadConfig(usr.HomeDir + "/.tskmgr.yml"); err != nil {
		return
	}
	if err = m.LoadTasks(); err != nil {
		return
	}
	return
}
