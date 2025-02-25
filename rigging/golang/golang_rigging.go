package golang

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/derrick/common"
	"github.com/alibaba/derrick/detectors/general"
	image "github.com/alibaba/derrick/detectors/image/golang"
	platform "github.com/alibaba/derrick/detectors/platform/golang"
)

const (
	Platform = "Golang"
)

type GolangRigging struct {
}

func (rig GolangRigging) Detect(workspace string) (bool, string) {
	var detected bool
	err := filepath.Walk(workspace, func(workspace string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".go") {
			detected = true
			return io.EOF
		}
		return nil
	})
	if err == io.EOF && detected {
		return true, Platform
	}
	return false, ""
}

func (rig GolangRigging) Compile(dockerImage string) (map[string]string, error) {
	dr := &common.DetectorReport{
		Nodes: map[string]common.DetectorReport{},
		Store: map[string]string{},
	}

	if err := dr.RegisterDetector(general.ImageRepoDetector{DockerImage: dockerImage}, common.Meta); err != nil {
		return nil, err
	}
	if err := dr.RegisterDetector(image.GolangVersionDetector{}, common.Dockerfile); err != nil {
		return nil, err
	}
	if err := dr.RegisterDetector(platform.PackageNameDetector{}, common.Dockerfile); err != nil {
		return nil, err
	}

	//if err := dr.RegisterDetector(general.ImageRepoDetector{}, Jenkinsfile); err != nil {
	//	return nil, err
	//}

	//if err := dr.RegisterDetector(general.ImageRepoDetector{}, common.DockerCompose); err != nil {
	//	return nil, err
	//}
	//if err := dr.RegisterDetector(general.ImageRepoDetector{}, common.KubernetesDeployment); err != nil {
	//	return nil, err
	//}
	if err := dr.RegisterDetector(general.DerrickDetector{}, common.KubernetesDeployment); err != nil {
		return nil, err
	}
	return dr.GenerateReport(), nil
}
