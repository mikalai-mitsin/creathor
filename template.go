package main

import (
	"errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"os"
	"path"
	"strings"
	"text/template"
)

var funcMap = template.FuncMap{
	"ToUpper": strings.ToUpper,
	"ToLower": strings.ToLower,
	"Title":   func(value string) string { return cases.Title(language.English).String(value) },
}

type Template struct {
	SourcePath      string
	DestinationPath string
	Name            string
}

func (t *Template) renderToFile(data interface{}) error {
	a := path.Base(t.SourcePath)
	tmpl, err := template.New(a).Funcs(funcMap).ParseFS(content, t.SourcePath)
	if err != nil {
		return NewBadTemplateError(err.Error())
	}
	file, err := os.Create(t.DestinationPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return NewDirectoryNotExistsError(t.DestinationPath)
		}
		if errors.Is(err, os.ErrPermission) {
			return NewPermissionError(t.DestinationPath)
		}
		return NewUnexpectedBehaviorError(err.Error())
	}
	if err := tmpl.Execute(file, data); err != nil {
		return NewUnexpectedBehaviorError(err.Error())
	}
	return nil
}
