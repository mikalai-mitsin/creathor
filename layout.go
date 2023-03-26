package main

import (
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

func createDirectories(project *configs.Project) error {
	directories := []string{
		path.Join(destinationPath, "build"),
		path.Join(destinationPath, "cmd"),
		path.Join(destinationPath, "cmd", project.Name),
		path.Join(destinationPath, "configs"),
		path.Join(destinationPath, "dist"),
		path.Join(destinationPath, "docs"),
		path.Join(destinationPath, "docs", ".chglog"),
		path.Join(destinationPath, "internal"),
		path.Join(destinationPath, "internal", "configs"),
		path.Join(destinationPath, "internal", "domain", "errs"),
		path.Join(destinationPath, "internal", "domain", "interceptors"),
		path.Join(destinationPath, "internal", "domain", "models"),
		path.Join(destinationPath, "internal", "domain", "models", "mock"),
		path.Join(destinationPath, "internal", "domain", "repositories"),
		path.Join(destinationPath, "internal", "domain", "usecases"),
		path.Join(destinationPath, "internal", "interceptors"),
		path.Join(destinationPath, "internal", "interfaces"),
		path.Join(destinationPath, "internal", "interfaces", "postgres"),
		path.Join(destinationPath, "internal", "interfaces", "postgres", "migrations"),
		path.Join(destinationPath, "internal", "usecases"),
		path.Join(destinationPath, "internal", "repositories"),
		path.Join(destinationPath, "internal", "repositories", "postgres"),
		path.Join(destinationPath, "internal", "repositories", "jwt"),
		path.Join(destinationPath, "pkg"),
		path.Join(destinationPath, "pkg", "clock"),
		path.Join(destinationPath, "pkg", "log"),
		path.Join(destinationPath, "pkg", "utils"),
		path.Join(destinationPath, "pkg", "postgresql"),
	}
	if project.RESTEnabled {
		directories = append(
			directories,
			path.Join(destinationPath, "internal", "interfaces", "rest"),
		)
	}
	if project.GRPCEnabled {
		directories = append(
			directories,
			path.Join(destinationPath, "internal", "interfaces", "grpc"),
			path.Join(destinationPath, "api", "proto", project.ProtoPackage(), "v1"),
		)
	}
	if project.GRPCEnabled && project.GatewayEnabled {
		directories = append(
			directories,
			path.Join(destinationPath, "internal", "interfaces", "gateway"),
		)
	}
	if project.UptraceEnabled {
		directories = append(
			directories,
			path.Join(destinationPath, "internal", "interfaces", "uptrace"),
		)
	}
	for _, directory := range directories {
		if err := os.MkdirAll(directory, 0777); err != nil {
			return NewUnexpectedBehaviorError(err.Error())
		}
	}
	return nil
}

func CreateLayout(project *configs.Project) error {
	if err := createDirectories(project); err != nil {
		return err
	}
	files := []*Template{
		{
			SourcePath:      "templates/cmd/service/main.go.tmpl",
			DestinationPath: path.Join(destinationPath, "cmd", project.Name, "main.go"),
			Name:            "service main",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "config.toml"),
			Name:            "main config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "ci.toml"),
			Name:            "ci config",
		},
		{
			SourcePath:      "templates/configs/config.toml.tmpl",
			DestinationPath: path.Join(destinationPath, "configs", "test.toml"),
			Name:            "test config",
		},
		{
			SourcePath:      "templates/internal/configs/config.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "config.go"),
			Name:            "config struct",
		},
		{
			SourcePath:      "templates/internal/configs/testing.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "testing.go"),
			Name:            "config testing",
		},
		{
			SourcePath:      "templates/internal/configs/config_test.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "configs", "config_test.go"),
			Name:            "config tests",
		},
		{
			SourcePath: "templates/internal/domain/errs/errors_test.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"domain",
				"errs",
				"errors_test.go",
			),
			Name: "domain errors tests",
		},
		{
			SourcePath:      "templates/pkg/clock/clock.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "clock", "clock.go"),
			Name:            "clock",
		},
		{
			SourcePath:      "templates/pkg/postgresql/search.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "postgresql", "search.go"),
			Name:            "search",
		},
		{
			SourcePath:      "templates/pkg/log/logger.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "log", "logger.go"),
			Name:            "logger interface",
		},
		{
			SourcePath:      "templates/pkg/utils/pointer.go.tmpl",
			DestinationPath: path.Join(destinationPath, "pkg", "utils", "pointer.go"),
			Name:            "utils pointer",
		},
		{
			SourcePath:      "templates/go.mod.tmpl",
			DestinationPath: path.Join(destinationPath, "go.mod"),
			Name:            "go.mod",
		},
		{
			SourcePath:      "templates/version.go.tmpl",
			DestinationPath: path.Join(destinationPath, "version.go"),
			Name:            "version",
		},
		{
			SourcePath:      "templates/docs/README.md.tmpl",
			DestinationPath: path.Join(destinationPath, "README.md"),
			Name:            "README.md",
		},
		{
			SourcePath:      "templates/docs/chglog/CHANGELOG.tpl.md.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", ".chglog", "CHANGELOG.tpl.md"),
			Name:            ".chglog/CHANGELOG.tpl.md",
		},
		{
			SourcePath:      "templates/docs/chglog/config.yml.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", ".chglog", "config.yml"),
			Name:            ".chglog/config.yml",
		},
		{
			SourcePath:      "templates/docs/CHANGELOG.md.tmpl",
			DestinationPath: path.Join(destinationPath, "docs", "CHANGELOG.md"),
			Name:            "CHANGELOG.md",
		},
		{
			SourcePath: "templates/internal/interfaces/postgres/postgres.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"interfaces",
				"postgres",
				"postgres.go",
			),
			Name: "postgres",
		},
		{
			SourcePath: "templates/internal/interfaces/postgres/testing.go.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"interfaces",
				"postgres",
				"testing.go",
			),
			Name: "postgres testing",
		},
		{
			SourcePath: "templates/internal/interfaces/postgres/migrations/init.sql.tmpl",
			DestinationPath: path.Join(
				destinationPath,
				"internal",
				"interfaces",
				"postgres",
				"migrations",
				"000001_init.up.sql",
			),
			Name: "postgres init migration",
		},

		{
			SourcePath:      "templates/internal/domain/models/types.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "domain", "models", "types.go"),
			Name:            "model types",
		},
	}
	if project.Auth {
		files = append(
			files,
			&Template{
				SourcePath: "templates/internal/domain/usecases/auth.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"usecases",
					"auth.go",
				),
				Name: "auth usecase",
			},
			&Template{
				SourcePath: "templates/internal/domain/models/auth.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"models",
					"auth.go",
				),
				Name: "auth models",
			},
			&Template{
				SourcePath: "templates/internal/domain/models/auth_mock.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"models",
					"mock",
					"auth.go",
				),
				Name: "auth mock models",
			},
			&Template{
				SourcePath: "templates/internal/domain/models/permission.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"models",
					"permission.go",
				),
				Name: "auth permission",
			},
			&Template{
				SourcePath: "templates/internal/domain/models/user.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"models",
					"user.go",
				),
				Name: "user model",
			},
			&Template{
				SourcePath: "templates/internal/domain/models/user_mock.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"models",
					"mock",
					"user.go",
				),
				Name: "user mock model",
			},
			&Template{
				SourcePath: "templates/internal/domain/repositories/user.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"repositories",
					"user.go",
				),
				Name: "user repository",
			},
			&Template{
				SourcePath: "templates/internal/domain/usecases/user.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"usecases",
					"user.go",
				),
				Name: "user usecase",
			},
			&Template{
				SourcePath: "templates/internal/domain/interceptors/user.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"interceptors",
					"user.go",
				),
				Name: "user interceptor",
			},
			&Template{
				SourcePath: "templates/internal/domain/interceptors/auth.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"interceptors",
					"auth.go",
				),
				Name: "auth interceptor",
			},
			&Template{
				SourcePath: "templates/internal/domain/repositories/permission.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"repositories",
					"permission.go",
				),
				Name: "auth permission repository",
			},
			&Template{
				SourcePath: "templates/internal/domain/repositories/auth.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"domain",
					"repositories",
					"auth.go",
				),
				Name: "auth repository",
			},
			&Template{
				SourcePath:      "templates/internal/usecases/auth_test.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "usecases", "auth_test.go"),
				Name:            "test auth usecase implementation",
			},
			&Template{
				SourcePath:      "templates/internal/usecases/user_test.go.tmpl",
				DestinationPath: path.Join(destinationPath, "internal", "usecases", "user_test.go"),
				Name:            "test user usecase implementation",
			},
			&Template{
				SourcePath: "templates/internal/interceptors/user_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interceptors",
					"user_test.go",
				),
				Name: "test user interceptor implementation",
			},
			&Template{
				SourcePath: "templates/internal/interceptors/auth_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interceptors",
					"auth_test.go",
				),
				Name: "test auth interceptor implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/postgres/user.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"postgres",
					"user.go",
				),
				Name: "user repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/postgres/user_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"postgres",
					"user_test.go",
				),
				Name: "test user repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/postgres/permission.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"postgres",
					"permission.go",
				),
				Name: "user repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/postgres/permission_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"postgres",
					"permission_test.go",
				),
				Name: "test permission repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/jwt/auth.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"jwt",
					"auth.go",
				),
				Name: "user repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/repositories/jwt/auth_test.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"repositories",
					"jwt",
					"auth_test.go",
				),
				Name: "test auth repository implementation",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/permissions.up.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000002_permissions.up.sql",
				),
				Name: "postgres permissions migration up",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/permissions.down.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000002_permissions.down.sql",
				),
				Name: "postgres permissions migration down",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/groups.up.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000003_groups.up.sql",
				),
				Name: "postgres groups migration up",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/groups.down.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000003_groups.down.sql",
				),
				Name: "postgres groups migration down",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/group_permissions.up.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000004_group_permissions.up.sql",
				),
				Name: "postgres group permissions migration up",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/group_permissions.down.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000004_group_permissions.down.sql",
				),
				Name: "postgres group permissions migration down",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/users.up.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000005_users.up.sql",
				),
				Name: "postgres users migration up",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/postgres/migrations/users.down.sql.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"postgres",
					"migrations",
					"000005_users.down.sql",
				),
				Name: "postgres users migration down",
			},
		)
	}
	if project.RESTEnabled {
		files = append(
			files,
			&Template{
				SourcePath: "templates/internal/interfaces/rest/middleware.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"rest",
					"middleware.go",
				),
				Name: "rest middlewares",
			},
			&Template{
				SourcePath: "templates/internal/interfaces/rest/server.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"rest",
					"server.go",
				),
				Name: "rest server",
			},
		)
		if project.Auth {
			files = append(
				files,
				&Template{
					SourcePath: "templates/internal/interfaces/rest/auth.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"rest",
						"auth.go",
					),
					Name: "rest auth handler",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/rest/user.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"rest",
						"user.go",
					),
					Name: "rest user handler",
				},
			)
		}
	}
	if project.GRPCEnabled {
		files = append(
			files,
			&Template{
				SourcePath:      "templates/api/proto/buf.yaml.tmpl",
				DestinationPath: path.Join(destinationPath, "api", "proto", "buf.yaml"),
				Name:            "buf.yaml",
			},
			&Template{
				SourcePath:      "templates/buf.gen.yaml.tmpl",
				DestinationPath: path.Join(destinationPath, "buf.gen.yaml"),
				Name:            "buf.gen.yaml",
			},
			&Template{
				SourcePath:      "templates/buf.work.yaml.tmpl",
				DestinationPath: path.Join(destinationPath, "buf.work.yaml"),
				Name:            "buf.work.yaml",
			},
		)
		if project.Auth {
			files = append(files,
				&Template{
					SourcePath: "templates/api/proto/service/v1/auth.proto.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"api",
						"proto",
						project.ProtoPackage(),
						"v1",
						"auth.proto",
					),
					Name: "auth.proto",
				},
				&Template{
					SourcePath: "templates/api/proto/service/v1/user.proto.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"api",
						"proto",
						project.ProtoPackage(),
						"v1",
						"user.proto",
					),
					Name: "user.proto",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/grpc/auth.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"grpc",
						"auth.go",
					),
					Name: "grpc auth",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/grpc/auth_test.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"grpc",
						"auth_test.go",
					),
					Name: "grpc auth test",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/grpc/user.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"grpc",
						"user.go",
					),
					Name: "grpc user",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/grpc/user_test.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"grpc",
						"user_test.go",
					),
					Name: "grpc user test",
				},
				&Template{
					SourcePath: "templates/internal/interfaces/grpc/auth_middleware_test.go.tmpl",
					DestinationPath: path.Join(
						destinationPath,
						"internal",
						"interfaces",
						"grpc",
						"auth_middleware_test.go",
					),
					Name: "grpc middleware test",
				},
			)
		}
	}
	if project.GRPCEnabled && project.GatewayEnabled {
		files = append(
			files,
			&Template{
				SourcePath: "templates/internal/interfaces/gateway/server.go.tmpl",
				DestinationPath: path.Join(
					destinationPath,
					"internal",
					"interfaces",
					"gateway",
					"server.go",
				),
				Name: "gateway server",
			},
		)
	}
	for _, file := range files {
		if err := file.renderToFile(project); err != nil {
			return err
		}
	}
	return nil
}
