package springboot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/derrick/common"
	"github.com/alibaba/derrick/detectors/general"
	image "github.com/alibaba/derrick/detectors/image/java"
)

type SpringBootRigging struct {
}

func (rig SpringBootRigging) Detect(workspace string) (bool, common.LanguagePlatform) {
	pom := filepath.Join(workspace, "pom.xml")
	if _, err := os.Stat(pom); err == nil {
		data, err := ioutil.ReadFile(pom)
		if err != nil {
			return false, ""
		}
		if strings.Contains(string(data), "org.springframework.boot") {
			return true, common.JavaSpringBoot
		}
	}
	return false, ""
}

func (rig SpringBootRigging) Compile(dockerImage string) (map[string]string, error) {
	dr := &common.DetectorReport{
		Nodes: map[string]common.DetectorReport{},
		Store: map[string]string{},
	}
	if err := dr.RegisterDetector(general.ImageRepoDetector{DockerImage: dockerImage}, common.Meta); err != nil {
		return nil, err
	}
	if err := dr.RegisterDetector(image.JavaVersionDetector{}, common.Dockerfile); err != nil {
		return nil, err
	}
	//if err := dr.RegisterDetector(platform.PackageNameDetector{}, common.Dockerfile); err != nil {
	//	return nil, err
	//}
	if err := dr.RegisterDetector(general.DerrickDetector{}, common.KubernetesDeployment); err != nil {
		return nil, err
	}
	return dr.GenerateReport(), nil
}
