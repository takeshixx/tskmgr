package manager

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/takeshixx/tskmgr/internal/tasks"
	"gopkg.in/yaml.v2"
)

// CreateTask creates a new tasks
func (m *Manager) CreateTask(title string) (t *tasks.Task, err error) {
	t, err = tasks.NewTask(title)
	if err != nil {
		return
	}
	newestIndex, err := m.getNewestTaskIndex()
	if err != nil {
		return
	}
	t.Directory = fmt.Sprintf("%d_%s", newestIndex+1, t.Directory)
	log.Printf("using directory: %s", t.Directory)
	m.Tasks = append(m.Tasks, t)
	return
}

// SaveTask safes a given task to the file system
// and returns the created path.
func (m *Manager) SaveTask(t *tasks.Task) (path string, err error) {
	if t.Directory == "" {
		err = fmt.Errorf("Task Directory for \"%s\" not set", t.Title)
		return
	}
	if !strings.HasSuffix(m.Config.TasksPath, string(os.PathSeparator)) {
		m.Config.TasksPath = m.Config.TasksPath + string(os.PathSeparator)
	}
	path = m.Config.TasksPath + t.Directory
	if _, err = os.Stat(path); err != os.ErrNotExist {
		if err = os.Mkdir(path, 0755); err != nil {
			return
		}
	}
	err = nil
	taskDescriptionPath := path + string(os.PathSeparator) + TaskNotesFileName
	if err = m.safeTaskDescription(t, taskDescriptionPath); err != nil {
		return
	}
	taskConfigPath := path + string(os.PathSeparator) + TaskConfigFileName
	if err = m.safeTaskConfig(t, taskConfigPath); err != nil {
		return
	}
	return
}

// safeTaskDescription formats and writes the task title and
// description to the task description file.
func (m *Manager) safeTaskDescription(t *tasks.Task, path string) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return
	}
	defer f.Close()
	fileData := t.Title
	fileData += "\n\n"
	fileData += t.Description
	_, err = f.WriteString(fileData)
	return
}

// safeTaskConfig safes the task object in YAML in the task
// configuration file.
func (m *Manager) safeTaskConfig(t *tasks.Task, path string) (err error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return
	}
	defer f.Close()
	yamlData, err := yaml.Marshal(t)
	if err != nil {
		return
	}
	_, err = f.Write(yamlData)
	return
}

// DeleteTask removes the file system directory and removes
// the given task from the global Tasks list.
func (m *Manager) DeleteTask(t *tasks.Task) (err error) {
	err = os.RemoveAll(m.Config.TasksPath + string(os.PathSeparator) + t.Directory)
	if err != nil {
		return
	}
	var newTaskList []*tasks.Task
	for _, task := range m.Tasks {
		if t == task {
			continue
		}
		newTaskList = append(newTaskList, task)
	}
	m.Tasks = newTaskList
	return
}

// LoadTask loads a single task.
func (m *Manager) LoadTask(taskDir string) (t *tasks.Task, err error) {
	taskPath := m.Config.TasksPath + string(os.PathSeparator) + taskDir

	taskConfigFile, err := os.Open(taskPath + string(os.PathSeparator) + TaskConfigFileName)
	if err != nil {
		return
	}
	defer taskConfigFile.Close()
	taskConfigData, err := ioutil.ReadAll(taskConfigFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(taskConfigData, &t)
	if err != nil {
		return
	}
	return
}

// LoadTasks reads all existing tasks from the
// file system.
func (m *Manager) LoadTasks() (err error) {
	m.Tasks = nil
	files, err := ioutil.ReadDir(m.Config.TasksPath)
	if err != nil {
		return
	}
	var task *tasks.Task
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if f.Name() == m.Config.TasksPath {
			continue
		}
		task, err = m.LoadTask(f.Name())
		if err != nil {
			err = nil
			continue
		}
		m.Tasks = append(m.Tasks, task)
	}
	return
}

// ReloadTask reloads a given task into the global Tasks list and returns it.
func (m *Manager) ReloadTask(t *tasks.Task) (tt *tasks.Task, err error) {
	tt, err = m.LoadTask(t.Directory)
	if err != nil {
		return
	}
	for i, ct := range m.Tasks {
		if ct == t {
			m.Tasks[i] = tt
		}
	}
	return
}
