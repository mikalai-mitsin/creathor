package models

import (
	"go/ast"
	"go/token"

	"github.com/018bf/creathor/internal/configs"
)

type FilterModel struct {
	model *configs.ModelConfig
}

func NewFilterModel(modelConfig *configs.ModelConfig) *FilterModel {
	return &FilterModel{model: modelConfig}
}

func (m *FilterModel) params() []*ast.Field {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{ast.NewIdent("IDs")},
			Type: &ast.ArrayType{
				Elt: ast.NewIdent("UUID"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"ids\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("PageSize")},
			Type:  &ast.StarExpr{X: ast.NewIdent("uint64")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"page_size\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("PageNumber")},
			Type:  &ast.StarExpr{X: ast.NewIdent("uint64")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"page_number\"`",
			},
		},
		{
			Names: []*ast.Ident{ast.NewIdent("OrderBy")},
			Type: &ast.ArrayType{
				Elt: ast.NewIdent("string"),
			},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"order_by\"`",
			},
		},
	}
	if m.model.SearchEnabled() {
		fields = append(fields, &ast.Field{
			Names: []*ast.Ident{ast.NewIdent("Search")},
			Type:  &ast.StarExpr{X: ast.NewIdent("string")},
			Tag: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "`json:\"search\"`",
			},
		})
	}
	return fields
}

func (m *FilterModel) Sync() error {
	structure := NewStructure(m.model.FileName(), m.model.FilterTypeName(), m.params())
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
