package core

import (
	"github.com/alibaba/derrick/common"
	"github.com/alibaba/derrick/rigging/golang"
	"github.com/alibaba/derrick/rigging/java"
	"github.com/alibaba/derrick/rigging/java/springboot"
	"github.com/alibaba/derrick/rigging/nodejs"
	"github.com/alibaba/derrick/rigging/php"
	"github.com/alibaba/derrick/rigging/python"
)

func LoadRiggings() []ExtensionPoint {
	riggings := []common.Rigging{golang.GolangRigging{}, java.JavaBasicRigging{}, springboot.SpringBootRigging{}, nodejs.NodeJSRigging{}, php.PHPRigging{}, python.PythonRigging{}}
	extensionPoints := make([]ExtensionPoint, len(riggings))
	for i, rig := range riggings {
		extensionPoints[i] = Register(rig)
	}
	return extensionPoints

	//TODO(zzxwill) Load developer's custom rigging
}
