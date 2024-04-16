package jsruntime

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dop251/goja"
)

func NewJSRuntime(name string) *JSRuntime {

	vm := goja.New()
	return &JSRuntime{
		name: name,
		vm:   vm,
	}
}

type JSRuntime struct {
	vm *goja.Runtime

	name string
}

func (j *JSRuntime) Close() error {
	return nil
}

func (j *JSRuntime) Load(modules ...string) error {
	for _, moduleFile := range modules {
		jsScript, err := JSAssetsFs.ReadFile(normalizedModulePath(moduleFile))
		if err != nil {
			return err
		}
		_, err = j.vm.RunScript(j.name, string(jsScript))
		if err != nil {
			return err
		}
		fmt.Printf("module %s is loaded into %s\n", moduleFile, j.name)
	}
	return nil
}

func normalizedModulePath(moduleFile string) string {
	if strings.HasPrefix(moduleFile, "assets") {
		return moduleFile
	}
	return filepath.Join("assets", moduleFile)
}

func (j *JSRuntime) Run(execScript string) (string, error) {
	v, err := j.vm.RunScript(j.name, execScript)
	if err != nil {
		return "", err
	}
	return v.Export().(string), nil
}
