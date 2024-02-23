package jsruntime

import (
	"fmt"
	"path/filepath"
	"strings"

	v8 "rogchap.com/v8go"
)

func NewJSRuntime(name string) *JSRuntime {

	env := v8.NewIsolate()
	return &JSRuntime{
		name: name,
		env:  env,
		ctx:  v8.NewContext(env),
	}
}

type JSRuntime struct {
	env *v8.Isolate
	ctx *v8.Context

	name string
}

func (j *JSRuntime) Close() error {
	j.env.Dispose()
	return nil
}

func (j *JSRuntime) Load(modules ...string) error {
	for _, moduleFile := range modules {
		jsScript, err := JSAssetsFs.ReadFile(normalizedModulePath(moduleFile))
		if err != nil {
			return err
		}
		_, err = j.ctx.RunScript(string(jsScript), j.name)
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
	v, err := j.ctx.RunScript(execScript, j.name)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}
