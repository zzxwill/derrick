package python

import (
	"os"
	"path/filepath"

	"github.com/alibaba/derrick/common"
)

type PythonRigging struct {
}

func (rig PythonRigging) Detect(workspace string) (bool, common.LanguagePlatform) {
	requirementsTxt := filepath.Join(workspace, "requirements.txt")
	setupPy := filepath.Join(workspace, "setup.py")
	if _, err := os.Stat(requirementsTxt); err == nil {
		if _, err := os.Stat(setupPy); err == nil {
			return true, common.Python
		}
	}
	return false, ""
}

func (rig PythonRigging) Compile(dockerImage string) (map[string]string, error) {
	return nil, nil
}
