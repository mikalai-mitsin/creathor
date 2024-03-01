package domain

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/018bf/creathor/internal/pkg/configs"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"golang.org/x/exp/slices"
)

type Domain struct {
	Config      *configs.DomainConfig
	Name        string
	Module      string
	ProtoModule string
	Models      []*Model
	UseCase     *Layer
	Repository  *Layer
	Interceptor *Layer
	GRPCHandler *Layer
	Auth        bool
}

func (m *Domain) SnakeName() string {
	return strcase.ToSnake(m.Name)
}
func (m *Domain) FileName() string {
	return fmt.Sprintf("%s.go", m.SnakeName())
}
func (m *Domain) TestFileName() string {
	return fmt.Sprintf("%s_test.go", m.SnakeName())
}

func (m *Domain) MigrationUpFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.up.sql", last+1, m.TableName())
}

func (m *Domain) MigrationDownFileName() string {
	last, err := lastMigration()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%06d_%s.down.sql", last+1, m.TableName())
}

func lastMigration() (int, error) {
	dir, err := os.ReadDir(path.Join("internal", "interfaces", "postgres", "migrations"))
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

func (m *Domain) CamelName() string {
	return strcase.ToCamel(m.Name)
}

func (m *Domain) LowerCamelName() string {
	return strcase.ToLowerCamel(m.Name)
}

func (m *Domain) DirName() string {
	return m.SnakeName()
}

func (m *Domain) ModelsImportPath() string {
	return fmt.Sprintf(`"%s/internal/app/%s/models"`, m.Module, m.DirName())
}

func (m *Domain) GetMainModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeMain },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) TableName() string {
	return strcase.ToSnake(inflection.Plural(m.Name))
}

func (m *Domain) SearchEnabled() bool {
	return slices.ContainsFunc(
		m.GetMainModel().Params,
		func(param *configs.Param) bool { return param.Search },
	)
}

func (m *Domain) GetCreateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeCreate },
	)
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) GetUpdateModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeUpdate },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) GetFilterModel() *Model {
	index := slices.IndexFunc(
		m.Models,
		func(model *Model) bool { return model.Type == ModelTypeFilter },
	)
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Domain) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.CamelName()))
}

func (m *Domain) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", m.CamelName())
}

func (m *Domain) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", m.CamelName())
}

func (m *Domain) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", m.CamelName())
}

func (m *Domain) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", m.CamelName())
}
