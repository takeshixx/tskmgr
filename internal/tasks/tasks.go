package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kennygrant/sanitize"
)

// SetDirectory checks and initializes the given path.
func (t *Task) SetDirectory(path string) (err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		err = fmt.Errorf("%s isn't a directory", path)
		return
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	t.Directory = abs
	return
}

// SetTitle sets the task title and initializes a
// directory name based on the title.
func (t *Task) SetTitle(title string) (err error) {
	if len(title) < 3 || len(title) > 256 {
		err = fmt.Errorf("Title %s has an invalid length: %d", title, len(title))
		return
	}
	t.Title = title
	t.Directory = sanitize.Path(title)
	return
}

// SetDescription sets the description of the task.
func (t *Task) SetDescription(desc string) (err error) {
	// TODO: should we do any checks on the description text?
	t.Description = desc
	return
}
