package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/markbates/pkger"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"

	"github.com/alibaba/derrick/common"
	"github.com/alibaba/derrick/core"
)

var projectPath, dockerImage string

func Init() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"ini"},
		Short:   "Detect application's platform and compile the application",
		Long:    "Detect application's platform and compile the application",
		Example: `derrick init`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return execute(projectPath, dockerImage)
		},
	}
	cmd.Flags().StringP("debug", "d", "", "debug mod")
	cmd.Flags().StringVarP(&projectPath, "project-path", "p", "", "Path of a project which is about to be detected")
	cmd.Flags().StringVarP(&dockerImage, "image", "i", "", "The image and its tag which will be built")
	return cmd
}

type SuitableRiggings struct {
	Platform       common.LanguagePlatform
	ExtensionPoint core.ExtensionPoint
}

func execute(workspace, dockerImage string) error {
	var err error
	if workspace == "" {
		workspace, err = os.Getwd()
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(workspace); err != nil {
		return err
	}
	suitableRiggings := detect(workspace)
	riggingNo := len(suitableRiggings)
	switch riggingNo {
	case 0:

		fmt.Println("Failed to detect your application's platform.\nMaybe you can upgrade Derrick to get more platforms supported.")
		return nil
	case 1:
		suitableRigging := suitableRiggings[0]
		return storeDetections(suitableRigging)
	case 2:
		var suitableRigging *SuitableRiggings
		var isSpringBootFlag = true
		for _, rig := range suitableRiggings {
			switch rig.Platform {
			case common.JavaBasic:
			case common.JavaSpringBoot:
				suitableRigging = rig
			default:
				isSpringBootFlag = false

			}
		}
		if isSpringBootFlag {
			return storeDetections(suitableRigging)
		}
		fmt.Println("More than one rigging can handle the application.")
		return nil

	default:
		// TODO(zzxwill) ask users to choose from one of them
		fmt.Println("More than one rigging can handle the application.")
		return nil
	}

}

func detect(projectPath string) []*SuitableRiggings {
	allRigging := core.LoadRiggings()
	if projectPath == "" {
		projectPath = "./"
	}
	var suitableRiggings []*SuitableRiggings
	for _, rig := range allRigging {
		success, platform := rig.Rigging.Detect(projectPath)
		if success {
			suitableRiggings = append(suitableRiggings,
				&SuitableRiggings{
					Platform:       platform,
					ExtensionPoint: core.ExtensionPoint{Rigging: rig.Rigging},
				})
		}
	}
	return suitableRiggings
}

func renderTemplates(rig common.Rigging, detectedContext map[string]string, destDir string) error {
	// TODO(zzxwill) PkgPath() returns github.com/alibaba/derrick/rigging/golang/templates
	// there might be a better solution get the direcotry of the templates
	pkgPath := strings.Join(strings.Split(reflect.TypeOf(rig).PkgPath(), "/")[3:], "/")
	absTemplateDir := filepath.Join("/", filepath.Clean(pkgPath))
	templateDir := filepath.Join(absTemplateDir, "templates")
	var templates []string
	err := pkger.Walk(absTemplateDir, func(path string, info os.FileInfo, err error) error {
		if info != nil && strings.HasSuffix(info.Name(), ".tmpl") {
			templates = append(templates, info.Name())
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, t := range templates {
		renderedTemplate, err := renderTemplate(templateDir, t, detectedContext)
		if err != nil {
			return err
		}
		renderedTemplateName := strings.Split(t, ".tmpl")
		if len(renderedTemplateName) != 2 {
			return fmt.Errorf("template %s is not in the right format", t)
		}
		if err := ioutil.WriteFile(filepath.Join(destDir, renderedTemplateName[0]), []byte(renderedTemplate), 0750); err != nil {
			return err
		}
	}
	return nil
}

func renderTemplate(templateDir, templateFile string, detectedContext map[string]string) (string, error) {
	var ctx common.TemplateRenderContext
	if err := mapstructure.Decode(detectedContext, &ctx); err != nil {
		return "", err
	}
	f, err := pkger.Open(filepath.Join(templateDir, templateFile))
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New(templateFile).Funcs(template.FuncMap(sprig.FuncMap())).Parse(string(data))
	if err != nil {
		return "", err
	}
	var wr bytes.Buffer
	err = tmpl.Execute(&wr, ctx)
	if err != nil {
		return "", err
	}
	return wr.String(), nil
}

func storeDetections(suitableRigging *SuitableRiggings) error {
	rig := suitableRigging.ExtensionPoint.Rigging
	detectedContext, err := rig.Compile(dockerImage)
	if err != nil {
		return err
	}
	if err := renderTemplates(rig, detectedContext, workspace); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	fmt.Printf("Successfully detected your platform is [%s] and compiled it successfully.\n", suitableRigging.Platform)

	// write configuration context to a file located in the application folder
	data, err := json.Marshal(detectedContext)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(workspace, common.DerrickApplicationConf), data, 0750); err != nil {
		return err
	}
	return nil
}
