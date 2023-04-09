package models

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/018bf/creathor/internal/configs"
)

type UpdateModel struct {
	model *configs.ModelConfig
}

func NewUpdateModel(modelConfig *configs.ModelConfig) *UpdateModel {
	return &UpdateModel{model: modelConfig}
}

func (m *UpdateModel) params() []*ast.Field {
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
	}
	for _, param := range m.model.Params {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(param.GetName())},
			Type:  &ast.StarExpr{X: astType(param.Type)},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("`json:\"%s\"`", param.Tag()),
			},
		})
	}
	return fields
}

func (m *UpdateModel) Sync() error {
	structure := NewStructure(m.model.FileName(), m.model.UpdateTypeName(), m.params())
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
