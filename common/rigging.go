package common

type Rigging interface {
	Detect(workspace string) (bool, LanguagePlatform)
	Compile(dockerImage string) (map[string]string, error)
}
