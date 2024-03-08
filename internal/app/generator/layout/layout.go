package layout

import (
	"github.com/018bf/creathor/internal/app/generator"
	"github.com/018bf/creathor/internal/pkg/configs"
)

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{project: project}
}

func (g *Generator) Sync() error {
	generators := []generator.Generator{
		NewCmdGenerator(g.project),
		NewDocsGenerator(g.project),
		NewBuilderGenerator(g.project),
		NewCIGenerator(g.project),
		NewDeploymentGenerator(g.project),
	}
	if g.project.GRPCEnabled {
		generators = append(generators, NewBufGenerator(g.project))
	}
	for _, g := range generators {
		if err := g.Sync(); err != nil {
			return err
		}
	}
	return nil
}
