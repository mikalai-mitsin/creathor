package layout

import (
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"
)

type BufGenerator struct {
	project *configs.Project
}

func NewBufGenerator(project *configs.Project) *BufGenerator {
	return &BufGenerator{project: project}
}

func (d *BufGenerator) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath:      "templates/api/proto/buf.yaml.tmpl",
			DestinationPath: path.Join(destinationPath, "api", "proto", "buf.yaml"),
			Name:            "buf.yaml",
		},
		{
			SourcePath:      "templates/buf.gen.yaml.tmpl",
			DestinationPath: path.Join(destinationPath, "buf.gen.yaml"),
			Name:            "buf.gen.yaml",
		},
		{
			SourcePath:      "templates/buf.work.yaml.tmpl",
			DestinationPath: path.Join(destinationPath, "buf.work.yaml"),
			Name:            "buf.work.yaml",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(d.project); err != nil {
			return err
		}
	}
	return nil
}
