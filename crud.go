package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

func CreateCRUD(data Model) error {
	files := []*Template{
		{
			SourcePath:      "templates/internal/domain/models/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", data.FileName()),
			Name:            "model",
		},
		{
			SourcePath:      "templates/internal/domain/models/crud_mock.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "models", "mock", data.FileName()),
			Name:            "model_mock",
		},
		{
			SourcePath:      "templates/internal/domain/repositories/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/domain/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/domain/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "domain", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/usecases/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.FileName()),
			Name:            "usecase",
		},
		{
			SourcePath:      "templates/internal/usecases/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "usecases", data.TestFileName()),
			Name:            "usecase test",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.FileName()),
			Name:            "interceptor",
		},
		{
			SourcePath:      "templates/internal/interceptors/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "interceptors", data.TestFileName()),
			Name:            "interceptor test",
		},
		{
			SourcePath:      "templates/internal/repositories/crud.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.FileName()),
			Name:            "repository",
		},
		{
			SourcePath:      "templates/internal/repositories/crud_test.go.tmpl",
			DestinationPath: filepath.Join(destinationPath, "internal", "repositories", data.TestFileName()),
			Name:            "repository test",
		},
		{
			SourcePath:      "templates/internal/interfaces/rest/crud.go.tmpl",
			DestinationPath: path.Join(destinationPath, "internal", "interfaces", "rest", data.FileName()),
			Name:            "rest mark",
		},
	}
	for _, tmpl := range files {
		if err := tmpl.renderToFile(data); err != nil {
			return err
		}
	}
	if err := addToDI("usecases", fmt.Sprintf("New%s", data.UseCaseTypeName())); err != nil {
		return err
	}
	if err := addToDI("interceptors", fmt.Sprintf("New%s", data.InterceptorTypeName())); err != nil {
		return err
	}
	if err := addToDI("repositories", fmt.Sprintf("New%s", data.RepositoryTypeName())); err != nil {
		return err
	}
	if err := addToDI("interfaces/rest", fmt.Sprintf("New%s", data.RESTHandlerTypeName())); err != nil {
		return err
	}
	if err := registerHandler(data.RESTHandlerVariableName(), data.RESTHandlerTypeName()); err != nil {
		return err
	}
	return nil
}

func registerHandler(variableName, typeName string) error {
	packagePath := filepath.Join(destinationPath, "internal", "interfaces", "rest")
	fileset := token.NewFileSet()
	tree, err := parser.ParseDir(fileset, packagePath, func(info fs.FileInfo) bool {
		return true
	}, parser.SkipObjectResolution)
	if err != nil {
		return err
	}
	for _, p := range tree {
		for filePath, file := range p.Files {
			for _, decl := range file.Decls {
				funcDecl, ok := decl.(*ast.FuncDecl)
				if ok {
					if funcDecl.Name.String() == "NewRouter" {
						field := &ast.Field{
							Doc: &ast.CommentGroup{
								List: nil,
							},
							Names: []*ast.Ident{
								{
									NamePos: 0,
									Name:    variableName,
									Obj:     nil,
								},
							},
							Type: &ast.StarExpr{
								Star: 0,
								X: &ast.Ident{
									Name: typeName,
								},
							},
							Tag: nil,
							Comment: &ast.CommentGroup{
								List: nil,
							},
						}
						funcDecl.Type.Params.List = append(funcDecl.Type.Params.List, field)
						registerCall := &ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										NamePos: 0,
										Name:    variableName,
										Obj:     nil,
									},
									Sel: &ast.Ident{
										NamePos: 0,
										Name:    "Register",
										Obj:     nil,
									},
								},
								Lparen: 0,
								Args: []ast.Expr{
									&ast.Ident{
										NamePos: 0,
										Name:    "router",
										Obj:     nil,
									},
								},
								Ellipsis: 0,
								Rparen:   0,
							},
						}
						le := len(funcDecl.Body.List)
						newBody := append(funcDecl.Body.List[:le-1], registerCall, funcDecl.Body.List[le-1])
						funcDecl.Body.List = newBody
						openFile, err := os.OpenFile(filePath, os.O_WRONLY, 0777)
						if err != nil {
							return err
						}
						if err := printer.Fprint(openFile, fileset, file); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
