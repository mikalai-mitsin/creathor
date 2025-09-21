package configs

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
)

func (m *EntityConfig) ImportPathPostgresRepositories() string {
	return fmt.Sprintf(`"%s/internal/app/%s/repositories/postgres/%s"`, m.Module, m.AppName(), m.DirName())
}

func (m *EntityConfig) ImportAliasPostgresRepositories() string {
	return fmt.Sprintf("%sPostgresRepositories", m.LowerCamelName())
}

func (m *EntityConfig) RepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) RepositoryVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) PostgresDTOTypeName() string {
	return fmt.Sprintf("%sDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) PostgresDTOListTypeName() string {
	return fmt.Sprintf("%sListDTO", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryPrivateVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryPublicVariableName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryTypeName() string {
	return fmt.Sprintf("%sRepository", strcase.ToCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryInterfaceName() string {
	return fmt.Sprintf("%sRepository", strcase.ToLowerCamel(m.Name))
}

func (m *EntityConfig) GetRepositoryConstructorName() string {
	return fmt.Sprintf("New%s", m.GetRepositoryTypeName())
}

func (m *EntityConfig) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *EntityConfig) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *Param) bool { return param.Search },
	)
}

func (m *EntityConfig) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *EntityConfig) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func (m *EntityConfig) SearchVector() string {
	var params []string
	for _, param := range m.Params {
		if param.Search {
			params = append(params, param.Tag())
		}
	}
	vector := fmt.Sprintf("to_tsvector('english', %s)", strings.Join(params, " || "))
	return vector
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "pkg", "postgres", "migrations"))
	if err != nil {
		return 0, err
	}
	var files []string
	for _, entry := range dir {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	last := files[len(files)-1]
	n, _, _ := strings.Cut(strings.Trim(last, "0"), "_")
	index, err := strconv.Atoi(n)
	if err != nil {
		return 0, err
	}
	return index, nil
}
