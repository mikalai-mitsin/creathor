package i18n

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (c *Generator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath:      "templates/internal/pkg/i18n/translator.go.tmpl",
			DestinationPath: path.Join("internal", "pkg", "i18n", "translator.go"),
			Name:            "i18n",
		},
		{
			SourcePath:      "templates/internal/pkg/i18n/translations/en/default.po.tmpl",
			DestinationPath: path.Join("internal", "pkg", "i18n", "translations", "en", "default.po"),
			Name:            "i18n",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
