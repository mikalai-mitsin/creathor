package persmission

import (
	"path"

	"github.com/018bf/creathor/internal/pkg/configs"
	"github.com/018bf/creathor/internal/pkg/tmpl"
)

var destinationPath = "."

type Repository struct {
	project *configs.Project
}

func NewRepository(project *configs.Project) *Repository {
	return &Repository{project: project}
}

func (c *Repository) Sync() error {
	files := []*tmpl.Template{
		{
			SourcePath: "templates/internal/user/repositories/postgres/permission.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"user",
				"repositories",
				"postgres",
				"permission.go",
			),
			Name: "permission repository implementation",
		},
		{
			SourcePath: "templates/internal/user/repositories/postgres/permission_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"app",
				"user",
				"repositories",
				"postgres",
				"permission_test.go",
			),
			Name: "test permission repository implementation",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/permissions.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000002_permissions.up.sql",
			),
			Name: "postgres permissions migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/permissions.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000002_permissions.down.sql",
			),
			Name: "postgres permissions migration down",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/groups.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000003_groups.up.sql",
			),
			Name: "postgres groups migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/groups.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000003_groups.down.sql",
			),
			Name: "postgres groups migration down",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/group_permissions.up.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000004_group_permissions.up.sql",
			),
			Name: "postgres group permissions migration up",
		},
		{
			SourcePath: "templates/internal/pkg/postgres/migrations/group_permissions.down.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"pkg",
				"postgres",
				"migrations",
				"000004_group_permissions.down.sql",
			),
			Name: "postgres group permissions migration down",
		},
	}
	for _, file := range files {
		if err := file.RenderToFile(c.project); err != nil {
			return err
		}
	}
	return nil
}
