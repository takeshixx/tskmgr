package tasks

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	Directory   string    `yaml:"directory"`
	Created     time.Time `yaml:"created"`
	Modified    time.Time `yaml:"modified"`
	Finished    time.Time `yaml:"finished"`
	Progress    int       `yaml:"progress"`
	Done        bool      `yaml:"done"`
	uiInfo      string
}

func (t *Task) GetUIText() string {
	return t.uiInfo
}

func (t *Task) ClearUIText() {
	t.uiInfo = ""
}

func (t *Task) RenderUIText(path string) (err error) {
	var f *os.File
	_, err = os.Stat(path)
	if err != nil {
		return
	}
	f, err = os.Open(path)
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(f)
	var line string
	var code bool
	for scanner.Scan() {
		line = scanner.Text()
		if strings.HasPrefix(line, "#") {
			line = fmt.Sprintf("[#ff0000::b]%s[white:black:-]", line)
		}
		if strings.HasPrefix(strings.TrimSpace(line), "* ") || strings.HasPrefix(strings.TrimSpace(line), "- ") {
			var prefix string
			if strings.HasPrefix(strings.TrimSpace(line), "* ") {
				prefix = line[0:strings.Index(line, "* ")]
				line = line[strings.Index(line, "* ")+2:]
			} else if strings.HasPrefix(strings.TrimSpace(line), "- ") {
				prefix = line[0:strings.Index(line, "- ")]
				line = line[strings.Index(line, "- ")+2:]
			}
			line = fmt.Sprintf("%s\t[::b]•[::-] %s", prefix, line)
		}
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			if code {
				line = fmt.Sprintf("\t[black:yellowgreen]%s[white:black]", line)
			}
			code = !code
		}
		if code {
			line = fmt.Sprintf("\t[black:yellowgreen]%s[white:black]", line)
		}
		if line == "---" {
			line = strings.Repeat("⎯", 50)
		}
		t.uiInfo += fmt.Sprintf("%s\n", line)
	}
	err = scanner.Err()
	return
}

// NewTask creates a new task object for a given title.
func NewTask(title string) (t *Task, err error) {
	t = &Task{}
	if err = t.SetTitle(title); err != nil {
		return
	}
	t.Created = time.Now()
	t.Progress = 0
	t.Done = false
	return
}
