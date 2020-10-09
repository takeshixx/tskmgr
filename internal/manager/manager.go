package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// LoadConfig parses the configuration file.
func (m *Manager) LoadConfig(cfgPath string) (err error) {
	if _, err = os.Stat(cfgPath); err != nil {
		return
	}
	cfgFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return
	}
	var config *Config
	err = yaml.Unmarshal(cfgFile, &config)
	if err != nil {
		return
	}
	// Check mandatory config params
	if config.TasksPath == "" {
		err = fmt.Errorf("TasksPath configuration is required")
		return
	}
	m.Config = config
	return
}

func (m *Manager) getNewestTaskIndex() (ret int, err error) {
	files, err := ioutil.ReadDir(m.Config.TasksPath)
	if err != nil {
		return
	}
	var pathInt int
	ret = 0
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if f.Name() == m.Config.TasksPath {
			continue
		}
		pathParts := strings.SplitN(f.Name(), "_", 2)
		pathInt, err = strconv.Atoi(pathParts[0])
		if err != nil {
			// Ignore directories without int prefixes
			err = nil
			continue
		}
		if pathInt > ret {
			ret = pathInt
		}
	}
	return
}
