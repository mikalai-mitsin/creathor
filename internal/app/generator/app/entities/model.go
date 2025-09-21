package entities

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func astType(t string) ast.Expr {
	if strings.HasPrefix(t, "*") {
		return &ast.StarExpr{
			X: astType(strings.TrimPrefix(t, "*")),
		}
	}
	if strings.HasPrefix(t, "[]") {
		return &ast.ArrayType{
			Elt: astType(strings.TrimPrefix(t, "[]")),
		}
	}
	if strings.HasPrefix(t, "entities.") {
		return ast.NewIdent(strings.TrimPrefix(t, "entities."))
	}
	return ast.NewIdent(t)
}

type Model struct {
	entity       *configs.Entity
	entityConfig configs.EntityConfig
}

func NewModel(entity *configs.Entity, entityConfig configs.EntityConfig) *Model {
	return &Model{
		entity:       entity,
		entityConfig: entityConfig,
	}
}

func (m *Model) params() []*ast.Field {
	fields := make([]*ast.Field, len(m.entity.Params))
	for i, param := range m.entity.Params {
		fields[i] = &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(param.GetName())},
			Type:  astType(param.Type),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
			},
		}
	}
	return fields
}

func (m *Model) Sync() error {
	err := os.MkdirAll(path.Dir(m.entityConfig.FileName()), 0777)
	if err != nil {
		return err
	}
	structure := NewStructure(m.entityConfig.FileName(), m.entity.Name, m.params(), m.entityConfig)
	if err := structure.Sync(); err != nil {
		return err
	}
	if m.entity.Validation {
		validate := NewValidate(structure.spec(), m.entityConfig)
		if err := validate.Sync(); err != nil {
			return err
		}
	}
	if m.entity.Mock {
		mock := NewMock(structure.spec(), m.entityConfig)
		if err := mock.Sync(); err != nil {
			return err
		}
	}
	ordering := NewOrdering(m.entityConfig)
	if err := ordering.Sync(); err != nil {
		return err
	}
	return nil
}
