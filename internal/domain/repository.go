package domain

import (
	"fmt"
	"go/ast"

	"github.com/018bf/creathor/internal/configs"
)

func NewRepository(m *configs.DomainConfig) *Layer {
	layer := &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.RepositoryTypeName(),
		Variable: m.RepositoryVariableName(),
		Methods: []*Method{
			{
				Name: "Create",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent(m.Variable()),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "List",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("filter"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("[]*models.%s", m.ModelName())),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Count",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("filter"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("uint64"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Get",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Update",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent(m.Variable()),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Delete",
				Args: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("ctx"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("uuid"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
	if m.Model == "user" {
		layer.Methods = append(layer.Methods, &Method{
			Name: "GetByEmail",
			Args: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("ctx"),
					},
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("context"),
						Sel: ast.NewIdent("Context"),
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent("email"),
					},
					Type: ast.NewIdent("string"),
				},
			},
			Return: []*ast.Field{
				{
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent(m.ModelName()),
						},
					},
				},
				{
					Type: ast.NewIdent("error"),
				},
			},
		})
	}
	return layer
}
