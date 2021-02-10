package core

import (
	"github.com/alibaba/derrick/common"
)

type ExtensionPoint struct {
	Rigging common.Rigging
}

func Register(rig common.Rigging) ExtensionPoint {
	return ExtensionPoint{
		Rigging: rig,
	}
}
