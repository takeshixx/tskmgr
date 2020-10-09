package manager

import (
	"os"
	"os/exec"

	"github.com/takeshixx/tskmgr/internal/tasks"
)

// OpenTaskInfo opens task_notes.md of a given task in
// the default editor.
func (m *Manager) OpenTaskInfo(t *tasks.Task) (err error) {
	path := m.Config.TasksPath + string(os.PathSeparator) + t.Directory + string(os.PathSeparator) + "task_notes.md"
	return m.openEditor(path)
}

// OpenTaskConfig opens .task.yml of a given task in
// the default editor.
func (m *Manager) OpenTaskConfig(t *tasks.Task) (err error) {
	path := m.Config.TasksPath + string(os.PathSeparator) + t.Directory + string(os.PathSeparator) + ".task.yml"
	return m.openEditor(path)
}

// OpenTaskVSCode opens the task directory of a given
// task in Visual Studio Code.
func (m *Manager) OpenTaskVSCode(t *tasks.Task) (err error) {
	path := m.Config.TasksPath + string(os.PathSeparator) + t.Directory
	return m.openExternal("code", path)
}

// OpenTaskFiles opens the directory of a given task
// in the defaulkt filemanager.
func (m *Manager) OpenTaskFiles(t *tasks.Task) (err error) {
	path := m.Config.TasksPath + string(os.PathSeparator) + t.Directory
	return m.openFileManager(path)
}

func (m *Manager) openEditor(path string) (err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	return m.openExternal(editor, path)
}

func (m *Manager) openFileManager(path string) (err error) {
	fm := os.Getenv("FILEMANAGER")
	if fm == "" {
		fm = "ranger"
	}
	return m.openExternal(fm, path)
}

func (m *Manager) openExternal(command, path string) (err error) {
	exe, err := exec.LookPath(command)
	if err != nil {
		return
	}
	cmd := exec.Command(exe, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return
}
