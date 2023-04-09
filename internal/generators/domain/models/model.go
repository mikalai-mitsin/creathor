package models

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/018bf/creathor/internal/configs"
)

type MainModel struct {
	model *configs.ModelConfig
}

func NewMainModel(modelConfig *configs.ModelConfig) *MainModel {
	return &MainModel{model: modelConfig}
}

func (m *MainModel) params() []*ast.Field {
	fields := []*ast.Field{
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("ID")},
			Type:  ast.NewIdent("UUID"),
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"id\"`",
			},
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("UpdatedAt")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("time"),
				Sel: ast.NewIdent("Time"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"updated_at\"`",
			},
			Comment: nil,
		},
		{
			Doc:   nil,
			Names: []*ast.Ident{ast.NewIdent("CreatedAt")},
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("time"),
				Sel: ast.NewIdent("Time"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"created_at\"`",
			},
			Comment: nil,
		},
	}
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

func (m *MainModel) Sync() error {
	if m.model.Auth {
		permissions := NewPerm(m.model.ModelName(), m.model.FileName())
		if err := permissions.Sync(); err != nil {
			return err
		}
	}
	structure := NewStructure(m.model.FileName(), m.model.ModelName(), m.params())
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
