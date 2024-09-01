package grpc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type RequestIDMiddleware struct {
	project *configs.Project
}

func NewRequestIDMiddleware(project *configs.Project) *RequestIDMiddleware {
	return &RequestIDMiddleware{project: project}
}
func (m RequestIDMiddleware) filename() string {
	return path.Join("internal", "pkg", "grpc", "request_id_middleware.go")
}

func (m RequestIDMiddleware) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "grpc",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"context"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/log"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/google/uuid"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc"`,
						},
					},
				},
			},
			// struct
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "RequestIDMiddleware",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{},
						},
					},
				},
			},
			// constructor
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewRequestIDMiddleware",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "RequestIDMiddleware",
									},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "RequestIDMiddleware",
										},
									},
								},
							},
						},
					},
				},
			},
			// unary server usecase
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "m",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "RequestIDMiddleware",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "UnaryServerInterceptor",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "ctx",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "context",
									},
									Sel: &ast.Ident{
										Name: "Context",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "req",
									},
								},
								Type: ast.NewIdent("any"),
							},
							{
								Names: []*ast.Ident{
									{
										Name: "_",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "UnaryServerInfo",
										},
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "handler",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "grpc",
									},
									Sel: &ast.Ident{
										Name: "UnaryHandler",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("any"),
							},
							{
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "newCtx",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "context",
										},
										Sel: &ast.Ident{
											Name: "WithValue",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "log",
											},
											Sel: &ast.Ident{
												Name: "RequestIDKey",
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "uuid",
														},
														Sel: &ast.Ident{
															Name: "New",
														},
													},
												},
												Sel: &ast.Ident{
													Name: "String",
												},
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "handler",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "newCtx",
										},
										&ast.Ident{
											Name: "req",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m RequestIDMiddleware) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: "RequestIDMiddleware",
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{},
		},
	}
}

func (m RequestIDMiddleware) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == "RequestIDMiddleware" {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = m.astStruct()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  []ast.Spec{structure},
			Rparen: 0,
		}
		file.Decls = append(file.Decls, gd)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m RequestIDMiddleware) astConstructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewRequestIDMiddleware",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "RequestIDMiddleware",
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.Ident{
									Name: "RequestIDMiddleware",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m RequestIDMiddleware) syncConstructor() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewRequestIDMiddleware" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = m.astConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m RequestIDMiddleware) astUnaryServerInterceptorMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "m",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "RequestIDMiddleware",
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "UnaryServerInterceptor",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "req",
							},
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "_",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "grpc",
								},
								Sel: &ast.Ident{
									Name: "UnaryServerInfo",
								},
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "handler",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "grpc",
							},
							Sel: &ast.Ident{
								Name: "UnaryHandler",
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "newCtx",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "context",
								},
								Sel: &ast.Ident{
									Name: "WithValue",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "log",
									},
									Sel: &ast.Ident{
										Name: "RequestIDKey",
									},
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "uuid",
												},
												Sel: &ast.Ident{
													Name: "New",
												},
											},
										},
										Sel: &ast.Ident{
											Name: "String",
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{
								Name: "handler",
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "newCtx",
								},
								&ast.Ident{
									Name: "req",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (m RequestIDMiddleware) syncUnaryServerInterceptorMethod() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "UnaryServerInterceptor" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = m.astUnaryServerInterceptorMethod()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (m RequestIDMiddleware) Sync() error {
	if err := m.syncStruct(); err != nil {
		return err
	}
	if err := m.syncConstructor(); err != nil {
		return err
	}
	if err := m.syncUnaryServerInterceptorMethod(); err != nil {
		return err
	}
	return nil
}
