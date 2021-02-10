package php

import (
	"os"
	"path/filepath"

	"github.com/alibaba/derrick/common"
)

type PHPRigging struct {
}

func (rig PHPRigging) Detect(workspace string) (bool, common.LanguagePlatform) {
	composer := filepath.Join(workspace, "composer.json")
	if _, err := os.Stat(composer); err == nil {
		return true, common.PHP
	}
	return false, ""
}

func (rig PHPRigging) Compile(dockerImage string) (map[string]string, error) {
	return nil, nil
}
