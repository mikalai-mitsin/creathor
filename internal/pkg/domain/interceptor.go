package domain

import (
	"fmt"
	"go/ast"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

func NewInterceptor(m *configs.DomainConfig) *Layer {
	interceptor := &Layer{
		Auth:     m.Auth,
		Events:   m.KafkaEnabled,
		Name:     m.InterceptorTypeName(),
		Variable: m.InterceptorVariableName(),
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
							ast.NewIdent("create"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(m.CreateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
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
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("[]*entities.%s", m.ModelName())),
					},
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
								X:   ast.NewIdent("entities"),
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
							ast.NewIdent("update"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent(m.UpdateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
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
	return interceptor
}
