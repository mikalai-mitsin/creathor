package implementations

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
)

type UseCaseUser struct {
	project *configs.Project
}

func NewUseCaseUser(project *configs.Project) *UseCaseUser {
	return &UseCaseUser{project: project}
}

func (u UseCaseUser) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "usecases", "user.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = u.file()
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

func (u UseCaseUser) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "usecases",
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
							Value: `"strings"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"time"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/models"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/repositories"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/domain/usecases"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/clock"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/pkg/log"`, u.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UserUseCase",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "userRepository",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "repositories",
											},
											Sel: &ast.Ident{
												Name: "UserRepository",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "clock",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "clock",
											},
											Sel: &ast.Ident{
												Name: "Clock",
											},
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "logger",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "log",
											},
											Sel: &ast.Ident{
												Name: "Logger",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewUserUseCase",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "userRepository",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "repositories",
									},
									Sel: &ast.Ident{
										Name: "UserRepository",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "clock",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "clock",
									},
									Sel: &ast.Ident{
										Name: "Clock",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "logger",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "log",
									},
									Sel: &ast.Ident{
										Name: "Logger",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "usecases",
									},
									Sel: &ast.Ident{
										Name: "UserUseCase",
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
											Name: "UserUseCase",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "userRepository",
												},
												Value: &ast.Ident{
													Name: "userRepository",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "clock",
												},
												Value: &ast.Ident{
													Name: "clock",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "logger",
												},
												Value: &ast.Ident{
													Name: "logger",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Get",
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
										Name: "id",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "models",
									},
									Sel: &ast.Ident{
										Name: "UUID",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "User",
										},
									},
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
									Name: "user",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "userRepository",
											},
										},
										Sel: &ast.Ident{
											Name: "Get",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "id",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "GetByEmail",
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
										Name: "email",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "User",
										},
									},
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
									Name: "user",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "userRepository",
											},
										},
										Sel: &ast.Ident{
											Name: "GetByEmail",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "strings",
												},
												Sel: &ast.Ident{
													Name: "ToLower",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "email",
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "List",
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
										Name: "filter",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "UserFilter",
										},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.ArrayType{
									Elt: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
									},
								},
							},
							{
								Type: &ast.Ident{
									Name: "uint64",
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
									Name: "users",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "userRepository",
											},
										},
										Sel: &ast.Ident{
											Name: "List",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "filter",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "count",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "userRepository",
											},
										},
										Sel: &ast.Ident{
											Name: "Count",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.Ident{
											Name: "filter",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "users",
								},
								&ast.Ident{
									Name: "count",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Create",
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
										Name: "create",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "UserCreate",
										},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "User",
										},
									},
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
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "create",
											},
											Sel: &ast.Ident{
												Name: "Validate",
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "models",
											},
											Sel: &ast.Ident{
												Name: "User",
											},
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "ID",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `""`,
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "FirstName",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `""`,
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "LastName",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `""`,
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Password",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `""`,
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Email",
												},
												Value: &ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "strings",
														},
														Sel: &ast.Ident{
															Name: "ToLower",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "create",
															},
															Sel: &ast.Ident{
																Name: "Email",
															},
														},
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "CreatedAt",
												},
												Value: &ast.CompositeLit{
													Type: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "time",
														},
														Sel: &ast.Ident{
															Name: "Time",
														},
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "UpdatedAt",
												},
												Value: &ast.CompositeLit{
													Type: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "time",
														},
														Sel: &ast.Ident{
															Name: "Time",
														},
													},
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "GroupID",
												},
												Value: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "models",
													},
													Sel: &ast.Ident{
														Name: "GroupIDUser",
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "user",
									},
									Sel: &ast.Ident{
										Name: "SetPassword",
									},
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "create",
										},
										Sel: &ast.Ident{
											Name: "Password",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "u",
												},
												Sel: &ast.Ident{
													Name: "userRepository",
												},
											},
											Sel: &ast.Ident{
												Name: "Create",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ctx",
											},
											&ast.Ident{
												Name: "user",
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Update",
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
										Name: "update",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "UserUpdate",
										},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "models",
										},
										Sel: &ast.Ident{
											Name: "User",
										},
									},
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
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "update",
											},
											Sel: &ast.Ident{
												Name: "Validate",
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "userRepository",
											},
										},
										Sel: &ast.Ident{
											Name: "Get",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ctx",
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "update",
											},
											Sel: &ast.Ident{
												Name: "ID",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "update",
									},
									Sel: &ast.Ident{
										Name: "FirstName",
									},
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "user",
												},
												Sel: &ast.Ident{
													Name: "FirstName",
												},
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "update",
													},
													Sel: &ast.Ident{
														Name: "FirstName",
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "update",
									},
									Sel: &ast.Ident{
										Name: "LastName",
									},
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "user",
												},
												Sel: &ast.Ident{
													Name: "LastName",
												},
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "update",
													},
													Sel: &ast.Ident{
														Name: "LastName",
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "update",
									},
									Sel: &ast.Ident{
										Name: "Password",
									},
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "user",
												},
												Sel: &ast.Ident{
													Name: "SetPassword",
												},
											},
											Args: []ast.Expr{
												&ast.StarExpr{
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "update",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "update",
									},
									Sel: &ast.Ident{
										Name: "Email",
									},
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "user",
												},
												Sel: &ast.Ident{
													Name: "Email",
												},
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "strings",
													},
													Sel: &ast.Ident{
														Name: "ToLower",
													},
												},
												Args: []ast.Expr{
													&ast.StarExpr{
														X: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "update",
															},
															Sel: &ast.Ident{
																Name: "Email",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "u",
												},
												Sel: &ast.Ident{
													Name: "userRepository",
												},
											},
											Sel: &ast.Ident{
												Name: "Update",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ctx",
											},
											&ast.Ident{
												Name: "user",
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "nil",
											},
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "user",
								},
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserUseCase",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Delete",
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
										Name: "id",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "models",
									},
									Sel: &ast.Ident{
										Name: "UUID",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
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
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "u",
												},
												Sel: &ast.Ident{
													Name: "userRepository",
												},
											},
											Sel: &ast.Ident{
												Name: "Delete",
											},
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "ctx",
											},
											&ast.Ident{
												Name: "id",
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X: &ast.Ident{
									Name: "err",
								},
								Op: token.NEQ,
								Y: &ast.Ident{
									Name: "nil",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "err",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "nil",
								},
							},
						},
					},
				},
			},
		},
	}
}
