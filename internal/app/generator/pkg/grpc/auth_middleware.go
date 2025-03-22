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

type AuthMiddleware struct {
	project *configs.Project
}

func NewAuthMiddleware(project *configs.Project) *AuthMiddleware {
	return &AuthMiddleware{project: project}
}
func (m AuthMiddleware) filename() string {
	return path.Join("internal", "pkg", "grpc", "auth_middleware.go")
}

func (m AuthMiddleware) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("grpc"),
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
							Value: fmt.Sprintf(`"%s/internal/app/auth/entities"`, m.project.Module),
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
							Value: `"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/codes"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"google.golang.org/grpc/status"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("ctxKey"),
						Type: ast.NewIdent("int"),
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("UserKey"),
						},
						Type: ast.NewIdent("ctxKey"),
						Values: []ast.Expr{
							&ast.BinaryExpr{
								X:  ast.NewIdent("iota"),
								Op: token.ADD,
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: "1",
								},
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("headerAuthorize"),
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"authorization"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("expectedScheme"),
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"bearer"`,
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("AuthFromMD"),
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
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("string"),
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
								ast.NewIdent("val"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("metautils"),
												Sel: ast.NewIdent("ExtractIncoming"),
											},
											Args: []ast.Expr{
												ast.NewIdent("ctx"),
											},
										},
										Sel: ast.NewIdent("Get"),
									},
									Args: []ast.Expr{
										ast.NewIdent("headerAuthorize"),
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X:  ast.NewIdent("val"),
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.STRING,
									Value: `""`,
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("status"),
													Sel: ast.NewIdent("Errorf"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("codes"),
														Sel: ast.NewIdent("Unauthenticated"),
													},
													&ast.BinaryExpr{
														X: &ast.BasicLit{
															Kind:  token.STRING,
															Value: `"Request unauthenticated with "`,
														},
														Op: token.ADD,
														Y:  ast.NewIdent("expectedScheme"),
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("splits"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("strings"),
										Sel: ast.NewIdent("SplitN"),
									},
									Args: []ast.Expr{
										ast.NewIdent("val"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `" "`,
										},
										&ast.BasicLit{
											Kind:  token.INT,
											Value: "2",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.CallExpr{
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										ast.NewIdent("splits"),
									},
								},
								Op: token.LSS,
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: "2",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("status"),
													Sel: ast.NewIdent("Errorf"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("codes"),
														Sel: ast.NewIdent("Unauthenticated"),
													},
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"Bad authorization string"`,
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("strings"),
										Sel: ast.NewIdent("EqualFold"),
									},
									Args: []ast.Expr{
										&ast.IndexExpr{
											X: ast.NewIdent("splits"),
											Index: &ast.BasicLit{
												Kind:  token.INT,
												Value: "0",
											},
										},
										ast.NewIdent("expectedScheme"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
											},
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("status"),
													Sel: ast.NewIdent("Errorf"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("codes"),
														Sel: ast.NewIdent("Unauthenticated"),
													},
													&ast.BinaryExpr{
														X: &ast.BasicLit{
															Kind:  token.STRING,
															Value: `"Request unauthenticated with "`,
														},
														Op: token.ADD,
														Y:  ast.NewIdent("expectedScheme"),
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.IndexExpr{
									X: ast.NewIdent("splits"),
									Index: &ast.BasicLit{
										Kind:  token.INT,
										Value: "1",
									},
								},
								ast.NewIdent("nil"),
							},
						},
					},
				},
			},
		},
	}
}

func (m AuthMiddleware) astStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent("AuthMiddleware"),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("logger"),
						},
						Type: ast.NewIdent("Logger"),
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("authUseCase"),
						},
						Type: ast.NewIdent("AuthUseCase"),
					},
				},
			},
		},
	}
}

func (m AuthMiddleware) syncStruct() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == "AuthMiddleware" {
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

func (m AuthMiddleware) astConstructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("NewAuthMiddleware"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("authUseCase"),
						},
						Type: ast.NewIdent("AuthUseCase"),
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("logger"),
						},
						Type: ast.NewIdent("Logger"),
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("AuthMiddleware"),
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
								Type: ast.NewIdent("AuthMiddleware"),
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   ast.NewIdent("authUseCase"),
										Value: ast.NewIdent("authUseCase"),
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
	}
}

func (m AuthMiddleware) syncConstructor() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewAuthMiddleware" {
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

func (m AuthMiddleware) astAuthMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("AuthMiddleware"),
					},
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("context"),
							Sel: ast.NewIdent("Context"),
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
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{
									ast.NewIdent("token"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("entities"),
									Sel: ast.NewIdent("Token"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("stringToken"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("AuthFromMD"),
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
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
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("context"),
											Sel: ast.NewIdent("WithValue"),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("UserKey"),
											&ast.SelectorExpr{
												X:   ast.NewIdent("entities"),
												Sel: ast.NewIdent("Guest"),
											},
										},
									},
									ast.NewIdent("nil"),
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  ast.NewIdent("stringToken"),
						Op: token.EQL,
						Y: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `""`,
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("context"),
											Sel: ast.NewIdent("WithValue"),
										},
										Args: []ast.Expr{
											ast.NewIdent("ctx"),
											ast.NewIdent("UserKey"),
											&ast.SelectorExpr{
												X:   ast.NewIdent("entities"),
												Sel: ast.NewIdent("Guest"),
											},
										},
									},
									ast.NewIdent("nil"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("token"),
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("entities"),
								Sel: ast.NewIdent("Token"),
							},
							Args: []ast.Expr{
								ast.NewIdent("stringToken"),
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
									X:   ast.NewIdent("m"),
									Sel: ast.NewIdent("authUseCase"),
								},
								Sel: ast.NewIdent("Auth"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("token"),
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("newCtx"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("context"),
								Sel: ast.NewIdent("WithValue"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								ast.NewIdent("UserKey"),
								ast.NewIdent("user"),
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("newCtx"),
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (m AuthMiddleware) syncAuthMethod() error {
	fileset := token.NewFileSet()
	filename := m.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Auth" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = m.astAuthMethod()
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

func (m AuthMiddleware) astUnaryServerInterceptorMethod() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("m"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("AuthMiddleware"),
					},
				},
			},
		},
		Name: ast.NewIdent("UnaryServerInterceptor"),
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
							ast.NewIdent("req"),
						},
						Type: ast.NewIdent("any"),
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("_"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("grpc"),
								Sel: ast.NewIdent("UnaryServerInfo"),
							},
						},
					},
					{
						Names: []*ast.Ident{
							ast.NewIdent("handler"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("grpc"),
							Sel: ast.NewIdent("UnaryHandler"),
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
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("newCtx"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("m"),
								Sel: ast.NewIdent("Auth"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent("handler"),
							Args: []ast.Expr{
								ast.NewIdent("newCtx"),
								ast.NewIdent("req"),
							},
						},
					},
				},
			},
		},
	}
}

func (m AuthMiddleware) syncUnaryServerInterceptorMethod() error {
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

func (m AuthMiddleware) Sync() error {
	if err := m.syncStruct(); err != nil {
		return err
	}
	if err := m.syncConstructor(); err != nil {
		return err
	}
	if err := m.syncAuthMethod(); err != nil {
		return err
	}
	if err := m.syncUnaryServerInterceptorMethod(); err != nil {
		return err
	}
	return nil
}
