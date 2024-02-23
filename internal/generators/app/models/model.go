package models

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/018bf/creathor/internal/domain"
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
	if strings.HasPrefix(t, "models.") {
		return ast.NewIdent(strings.TrimPrefix(t, "models."))
	}
	return ast.NewIdent(t)
}

type Model struct {
	model  *domain.Model
	domain *domain.Domain
}

func NewModel(model *domain.Model, domain *domain.Domain) *Model {
	return &Model{
		model:  model,
		domain: domain,
	}
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
	err := os.MkdirAll(path.Dir(m.domain.FileName()), 0777)
	if err != nil {
		return err
	}
	if m.domain.Auth {
		permissions := NewPerm(m.model.Name, m.domain.FileName(), m.domain)
		if err := permissions.Sync(); err != nil {
			return err
		}
	}
	structure := NewStructure(m.domain.FileName(), m.model.Name, m.params(), m.domain)
	if err := structure.Sync(); err != nil {
		return err
	}
	if m.model.Validation {
		validate := NewValidate(structure.spec(), m.domain.FileName(), m.domain)
		if err := validate.Sync(); err != nil {
			return err
		}
	}
	if m.model.Name == "User" {
		password := NewPassword(m.domain)
		if err := password.Sync(); err != nil {
			return err
		}
	}
	if m.model.Mock {
		mock := NewMock(structure.spec(), m.domain.FileName(), m.domain)
		if err := mock.Sync(); err != nil {
			return err
		}
	}
	return nil
}
