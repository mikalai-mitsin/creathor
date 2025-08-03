package layout

import (
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type GitGenerator struct {
	project *configs.Project
}

func NewGitGenerator(project *configs.Project) *GitGenerator {
	return &GitGenerator{project: project}
}

func (b *GitGenerator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath:      "templates/git/gitignore.tmpl",
			DestinationPath: filepath.Join(destinationPath, ".gitignore"),
			Name:            ".gitignore",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(b.project); err != nil {
			return err
		}
	}
	return nil
}
