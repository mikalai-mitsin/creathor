package models

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/018bf/creathor/internal/mods"
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
	return ast.NewIdent(t)
}

type Model struct {
	model    *mods.Model
	filename string
}

func NewModel(model *mods.Model, filename string) *Model {
	return &Model{model: model, filename: filename}
}

func (m *Model) params() []*ast.Field {
	fields := make([]*ast.Field, len(m.model.Params))
	for i, param := range m.model.Params {
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
	permissions := NewPerm(m.model.Name, m.filename)
	if err := permissions.Sync(); err != nil {
		return err
	}
	structure := NewStructure(m.filename, m.model.Name, m.params())
	if err := structure.Sync(); err != nil {
		return err
	}
	if m.model.Validation {
		validate := NewValidate(structure.spec(), m.filename)
		if err := validate.Sync(); err != nil {
			return err
		}
	}
	if m.model.Mock {
		mock := NewMock(structure.spec(), m.filename)
		if err := mock.Sync(); err != nil {
			return err
		}
	}
	return nil
}
