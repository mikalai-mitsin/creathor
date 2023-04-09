package models

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/018bf/creathor/internal/configs"
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

type CreateModel struct {
	model *configs.ModelConfig
}

func NewCreateModel(modelConfig *configs.ModelConfig) *CreateModel {
	return &CreateModel{model: modelConfig}
}

func (m *CreateModel) params() []*ast.Field {
	var fields []*ast.Field
	for _, param := range m.model.Params {
		fields = append(fields, &ast.Field{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent(param.GetName())},
			Type:  astType(param.Type),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
			},
			Comment: nil,
		})
	}
	return fields
}

func (m *CreateModel) Sync() error {
	structure := NewStructure(m.model.FileName(), m.model.CreateTypeName(), m.params())
	if err := structure.Sync(); err != nil {
		return err
	}
	validate := NewValidate(structure.spec(), m.model.FileName())
	if err := validate.Sync(); err != nil {
		return err
	}
	mock := NewMock(structure.spec(), m.model.FileName())
	if err := mock.Sync(); err != nil {
		return err
	}
	return nil
}
