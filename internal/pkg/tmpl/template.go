package tmpl

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/018bf/creathor/internal/pkg/errs"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates/*
var content embed.FS

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
	"Title":   func(value string) string { return cases.Title(language.English).String(value) },
	"inc": func(i int) int {
		return i + 1
	},
	"add": func(i, j int) int {
		return i + j
	},
}

type Template struct {
	SourcePath      string
	DestinationPath string
	Name            string
}

func (t *Template) hasRewrite() bool {
	_, err := os.Stat(t.DestinationPath)
	if err != nil && os.IsNotExist(err) {
		return true
	}
	return false
}

func (t *Template) RenderToFile(data interface{}) error {
	if err := os.MkdirAll(path.Dir(t.DestinationPath), 0777); err != nil {
		return err
	}
	if !t.hasRewrite() {
		fmt.Printf("%s already exists.\n", t.Name)
		return nil
	}
	a := path.Base(t.SourcePath)
	tmpl, err := template.New(a).Funcs(funcMap).ParseFS(content, t.SourcePath)
	if err != nil {
		return errs.NewBadTemplateError(err.Error())
	}
	file, err := os.Create(t.DestinationPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errs.NewDirectoryNotExistsError(t.DestinationPath)
		}
		if errors.Is(err, os.ErrPermission) {
			return errs.NewPermissionError(t.DestinationPath)
		}
		return errs.NewUnexpectedBehaviorError(err.Error())
	}
	if err := tmpl.Execute(file, data); err != nil {
		return errs.NewUnexpectedBehaviorError(err.Error())
	}
	return nil
}
