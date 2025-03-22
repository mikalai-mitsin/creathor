package services

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type ServiceAuth struct {
	project *configs.Project
}

func NewServiceAuth(project *configs.Project) *ServiceAuth {
	return &ServiceAuth{project: project}
}

func (u ServiceAuth) Sync() error {
	fileset := token.NewFileSet()
	filename := filepath.Join("internal", "app", "auth", "services", "auth.go")
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
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

func (u ServiceAuth) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("services"),
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("userEntities"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/app/user/entities"`, u.project.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, u.project.Module),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("AuthService"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("authRepository"),
										},
										Type: ast.NewIdent("authRepository"),
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("userRepository"),
										},
										Type: ast.NewIdent("userRepository"),
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("logger"),
										},
										Type: ast.NewIdent("logger"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewAuthService"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("authRepository"),
								},
								Type: ast.NewIdent("authRepository"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("userRepository"),
								},
								Type: ast.NewIdent("userRepository"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("logger"),
								},
								Type: ast.NewIdent("logger"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("*AuthService"),
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
										Type: ast.NewIdent("AuthService"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("authRepository"),
												Value: ast.NewIdent("authRepository"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("userRepository"),
												Value: ast.NewIdent("userRepository"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("logger"),
												Value: ast.NewIdent("logger"),
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
								ast.NewIdent("u"),
							},
							Type: ast.NewIdent("AuthService"),
						},
					},
				},
				Name: ast.NewIdent("CreateToken"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("login"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Login"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("TokenPair"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("user"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("userRepository"),
										},
										Sel: ast.NewIdent("GetByEmail"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										&ast.SelectorExpr{
											X:   ast.NewIdent("login"),
											Sel: ast.NewIdent("Email"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("entities"),
													Sel: ast.NewIdent("TokenPair"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("user"),
											Sel: ast.NewIdent("CheckPassword"),
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("login"),
												Sel: ast.NewIdent("Password"),
											},
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("entities"),
													Sel: ast.NewIdent("TokenPair"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("authRepository"),
										},
										Sel: ast.NewIdent("Create"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("user"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("entities"),
													Sel: ast.NewIdent("TokenPair"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("nil"),
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
								ast.NewIdent("u"),
							},
							Type: ast.NewIdent("AuthService"),
						},
					},
				},
				Name: ast.NewIdent("CreateTokenByUser"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("user"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("userEntities"),
									Sel: ast.NewIdent("User"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("TokenPair"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("authRepository"),
										},
										Sel: ast.NewIdent("Create"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("user"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("entities"),
													Sel: ast.NewIdent("TokenPair"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("tokenPair"),
								ast.NewIdent("nil"),
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
								ast.NewIdent("u"),
							},
							Type: ast.NewIdent("AuthService"),
						},
					},
				},
				Name: ast.NewIdent("RefreshToken"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("refresh"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Token"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("TokenPair"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("pair"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("authRepository"),
										},
										Sel: ast.NewIdent("RefreshToken"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("refresh"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("entities"),
													Sel: ast.NewIdent("TokenPair"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("pair"),
								ast.NewIdent("nil"),
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
								ast.NewIdent("u"),
							},
							Type: ast.NewIdent("AuthService"),
						},
					},
				},
				Name: ast.NewIdent("ValidateToken"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("access"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Token"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("u"),
												Sel: ast.NewIdent("authRepository"),
											},
											Sel: ast.NewIdent("Validate"),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("access"),
										},
									},
								},
							},
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("nil"),
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
								ast.NewIdent("u"),
							},
							Type: ast.NewIdent("AuthService"),
						},
					},
				},
				Name: ast.NewIdent("Auth"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
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
									ast.NewIdent("access"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Token"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("userEntities"),
									Sel: ast.NewIdent("User"),
								},
							},
							{
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("userID"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("authRepository"),
										},
										Sel: ast.NewIdent("GetSubject"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										ast.NewIdent("access"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("userEntities"),
													Sel: ast.NewIdent("User"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("user"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("u"),
											Sel: ast.NewIdent("userRepository"),
										},
										Sel: ast.NewIdent("Get"),
									},
									Args: []ast.Expr{
										ast.NewIdent("ctx"),
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("uuid"),
												Sel: ast.NewIdent("UUID"),
											},
											Args: []ast.Expr{
												ast.NewIdent("userID"),
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("err"),
								Op: token.NEQ,
								Y:  ast.NewIdent("nil"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CompositeLit{
												Type: &ast.SelectorExpr{
													X:   ast.NewIdent("userEntities"),
													Sel: ast.NewIdent("User"),
												},
											},
											ast.NewIdent("err"),
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("user"),
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
	}

}
