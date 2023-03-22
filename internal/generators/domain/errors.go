package domain

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

type Errors struct {
	project *configs.Project
}

func NewErrors(project *configs.Project) *Errors {
	return &Errors{
		project: project,
	}
}

func (i Errors) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name: &ast.Ident{
			Name: "errs",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"bytes\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"database/sql\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"encoding/json\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"errors\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"fmt\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"reflect\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"text/template\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/lib/pq\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"go.uber.org/zap/zapcore\"",
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "validation",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-ozzo/ozzo-validation/v4\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "ErrorCode",
						},
						Type: &ast.Ident{
							Name: "uint",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Params",
						},
						Type: &ast.MapType{
							Key: &ast.Ident{
								Name: "string",
							},
							Value: &ast.Ident{
								Name: "string",
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "p",
								},
							},
							Type: &ast.Ident{
								Name: "Params",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "MarshalLogObject",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "encoder",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "zapcore",
									},
									Sel: &ast.Ident{
										Name: "ObjectEncoder",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "key",
							},
							Value: &ast.Ident{
								Name: "value",
							},
							Tok: token.DEFINE,
							X: &ast.Ident{
								Name: "p",
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "encoder",
												},
												Sel: &ast.Ident{
													Name: "AddString",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "key",
												},
												&ast.Ident{
													Name: "value",
												},
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
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeOK",
							},
						},
						Type: &ast.Ident{
							Name: "ErrorCode",
						},
						Values: []ast.Expr{
							&ast.Ident{
								Name: "iota",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeCanceled",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeUnknown",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeInvalidArgument",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeDeadlineExceeded",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeNotFound",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeAlreadyExists",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodePermissionDenied",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeResourceExhausted",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeFailedPrecondition",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeAborted",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeOutOfRange",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeUnimplemented",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeInternal",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeUnavailable",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeDataLoss",
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeUnauthenticated",
							},
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Error",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Code",
											},
										},
										Type: &ast.Ident{
											Name: "ErrorCode",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"code\"`",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Message",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"message\"`",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Params",
											},
										},
										Type: &ast.Ident{
											Name: "Params",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"params\"`",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "WithParam",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "key",
									},
									&ast.Ident{
										Name: "value",
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
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "AddParam",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "key",
									},
									&ast.Ident{
										Name: "value",
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "e",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "WithParams",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "params",
									},
								},
								Type: &ast.MapType{
									Key: &ast.Ident{
										Name: "string",
									},
									Value: &ast.Ident{
										Name: "string",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "key",
							},
							Value: &ast.Ident{
								Name: "value",
							},
							Tok: token.DEFINE,
							X: &ast.Ident{
								Name: "params",
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
												},
												Sel: &ast.Ident{
													Name: "AddParam",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "key",
												},
												&ast.Ident{
													Name: "value",
												},
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "e",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.Ident{
								Name: "Error",
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Error",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "string",
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
									Name: "data",
								},
								&ast.Ident{
									Name: "_",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "json",
										},
										Sel: &ast.Ident{
											Name: "Marshal",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "e",
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "data",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Is",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "tgt",
									},
								},
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "bool",
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
									Name: "target",
								},
								&ast.Ident{
									Name: "ok",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.TypeAssertExpr{
									X: &ast.Ident{
										Name: "tgt",
									},
									Type: &ast.StarExpr{
										X: &ast.Ident{
											Name: "Error",
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X: &ast.Ident{
									Name: "ok",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "false",
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "reflect",
										},
										Sel: &ast.Ident{
											Name: "DeepEqual",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "e",
										},
										&ast.Ident{
											Name: "target",
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
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "SetCode",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "code",
									},
								},
								Type: &ast.Ident{
									Name: "ErrorCode",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Code",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "code",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "SetMessage",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "message",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Message",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "message",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "SetParams",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "params",
									},
								},
								Type: &ast.MapType{
									Key: &ast.Ident{
										Name: "string",
									},
									Value: &ast.Ident{
										Name: "string",
									},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Params",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "params",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						&ast.Field{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "Error",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "AddParam",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "key",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "value",
									},
								},
								Type: &ast.Ident{
									Name: "string",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.IndexExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "e",
										},
										Sel: &ast.Ident{
											Name: "Params",
										},
									},
									Index: &ast.Ident{
										Name: "key",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "value",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "code",
									},
								},
								Type: &ast.Ident{
									Name: "ErrorCode",
								},
							},
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "message",
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
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "code",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.Ident{
													Name: "message",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
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
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewUnexpectedBehaviorError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "details",
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
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodeInternal",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Unexpected behavior.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
														},
													},
													Elts: []ast.Expr{
														&ast.KeyValueExpr{
															Key: &ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"details\"",
															},
															Value: &ast.Ident{
																Name: "details",
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
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewInvalidFormError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewError",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ErrorCodeInvalidArgument",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"The form sent is not valid, please correct the errors below.\"",
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
					Name: "NewInvalidParameter",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "message",
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
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
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
									Name: "e",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewError",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ErrorCodeInvalidArgument",
										},
										&ast.Ident{
											Name: "message",
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "e",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "FromValidationError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "err",
									},
								},
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
								},
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
											&ast.Ident{
												Name: "validationErrors",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "validation",
											},
											Sel: &ast.Ident{
												Name: "Errors",
											},
										},
									},
								},
							},
						},
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "validationErrorObject",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "validation",
											},
											Sel: &ast.Ident{
												Name: "ErrorObject",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "errors",
									},
									Sel: &ast.Ident{
										Name: "As",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "validationErrors",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "e",
											},
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "NewError",
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "ErrorCodeInvalidArgument",
													},
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: "\"The form sent is not valid, please correct the errors below.\"",
													},
												},
											},
										},
									},
									&ast.RangeStmt{
										Key: &ast.Ident{
											Name: "key",
										},
										Value: &ast.Ident{
											Name: "value",
										},
										Tok: token.DEFINE,
										X: &ast.Ident{
											Name: "validationErrors",
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.TypeSwitchStmt{
													Assign: &ast.AssignStmt{
														Lhs: []ast.Expr{
															&ast.Ident{
																Name: "t",
															},
														},
														Tok: token.DEFINE,
														Rhs: []ast.Expr{
															&ast.TypeAssertExpr{
																X: &ast.Ident{
																	Name: "value",
																},
															},
														},
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.CaseClause{
																List: []ast.Expr{
																	&ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "validation",
																		},
																		Sel: &ast.Ident{
																			Name: "ErrorObject",
																		},
																	},
																},
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "e",
																				},
																				Sel: &ast.Ident{
																					Name: "AddParam",
																				},
																			},
																			Args: []ast.Expr{
																				&ast.Ident{
																					Name: "key",
																				},
																				&ast.CallExpr{
																					Fun: &ast.Ident{
																						Name: "renderErrorMessage",
																					},
																					Args: []ast.Expr{
																						&ast.Ident{
																							Name: "t",
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
															&ast.CaseClause{
																List: []ast.Expr{
																	&ast.StarExpr{
																		X: &ast.Ident{
																			Name: "Error",
																		},
																	},
																},
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "e",
																				},
																				Sel: &ast.Ident{
																					Name: "AddParam",
																				},
																			},
																			Args: []ast.Expr{
																				&ast.Ident{
																					Name: "key",
																				},
																				&ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "t",
																					},
																					Sel: &ast.Ident{
																						Name: "Message",
																					},
																				},
																			},
																		},
																	},
																},
															},
															&ast.CaseClause{
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "e",
																				},
																				Sel: &ast.Ident{
																					Name: "AddParam",
																				},
																			},
																			Args: []ast.Expr{
																				&ast.Ident{
																					Name: "key",
																				},
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "value",
																						},
																						Sel: &ast.Ident{
																							Name: "Error",
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
												},
											},
										},
									},
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.Ident{
												Name: "e",
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "errors",
									},
									Sel: &ast.Ident{
										Name: "As",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "validationErrorObject",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "NewInvalidParameter",
												},
												Args: []ast.Expr{
													&ast.CallExpr{
														Fun: &ast.Ident{
															Name: "renderErrorMessage",
														},
														Args: []ast.Expr{
															&ast.Ident{
																Name: "validationErrorObject",
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
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "renderErrorMessage",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "object",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "validation",
									},
									Sel: &ast.Ident{
										Name: "ErrorObject",
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.Ident{
									Name: "string",
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
									Name: "parse",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "template",
												},
												Sel: &ast.Ident{
													Name: "New",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"message\"",
												},
											},
										},
										Sel: &ast.Ident{
											Name: "Parse",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "object",
												},
												Sel: &ast.Ident{
													Name: "Message",
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
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"\"",
											},
										},
									},
								},
							},
						},
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "tpl",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "bytes",
											},
											Sel: &ast.Ident{
												Name: "Buffer",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "_",
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "parse",
										},
										Sel: &ast.Ident{
											Name: "Execute",
										},
									},
									Args: []ast.Expr{
										&ast.UnaryExpr{
											Op: token.AND,
											X: &ast.Ident{
												Name: "tpl",
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "object",
												},
												Sel: &ast.Ident{
													Name: "Params",
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
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "tpl",
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
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "FromPostgresError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Names: []*ast.Ident{
									&ast.Ident{
										Name: "err",
									},
								},
								Type: &ast.Ident{
									Name: "error",
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
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
									Name: "e",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: &ast.Ident{
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodeInternal",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Unexpected behavior.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
														},
													},
													Elts: []ast.Expr{
														&ast.KeyValueExpr{
															Key: &ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"error\"",
															},
															Value: &ast.CallExpr{
																Fun: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "err",
																	},
																	Sel: &ast.Ident{
																		Name: "Error",
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
							},
						},
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "pqErr",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "pq",
												},
												Sel: &ast.Ident{
													Name: "Error",
												},
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "errors",
									},
									Sel: &ast.Ident{
										Name: "As",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
									&ast.UnaryExpr{
										Op: token.AND,
										X: &ast.Ident{
											Name: "pqErr",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
												},
												Sel: &ast.Ident{
													Name: "AddParam",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"details\"",
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "pqErr",
													},
													Sel: &ast.Ident{
														Name: "Detail",
													},
												},
											},
										},
									},
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
												},
												Sel: &ast.Ident{
													Name: "AddParam",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"message\"",
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "pqErr",
													},
													Sel: &ast.Ident{
														Name: "Message",
													},
												},
											},
										},
									},
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
												},
												Sel: &ast.Ident{
													Name: "AddParam",
												},
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"postgres_code\"",
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "fmt",
														},
														Sel: &ast.Ident{
															Name: "Sprint",
														},
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X: &ast.Ident{
																Name: "pqErr",
															},
															Sel: &ast.Ident{
																Name: "Code",
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
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "errors",
									},
									Sel: &ast.Ident{
										Name: "Is",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "err",
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sql",
										},
										Sel: &ast.Ident{
											Name: "ErrNoRows",
										},
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.Ident{
												Name: "e",
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "NewEntityNotFound",
												},
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "e",
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewEntityNotFound",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodeNotFound",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Entity not found.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
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
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewPermissionDeniedError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodePermissionDenied",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Permission denied.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
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
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewBadToken",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodeUnauthenticated",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Bad token.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
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
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewPermissionDenied",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
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
											Name: "Error",
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Code",
												},
												Value: &ast.Ident{
													Name: "ErrorCodePermissionDenied",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Message",
												},
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: "\"Permission denied.\"",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Params",
												},
												Value: &ast.CompositeLit{
													Type: &ast.MapType{
														Key: &ast.Ident{
															Name: "string",
														},
														Value: &ast.Ident{
															Name: "string",
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
				},
			},
		},
		Imports: []*ast.ImportSpec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"bytes\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"database/sql\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"encoding/json\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"errors\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"fmt\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"reflect\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"text/template\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/lib/pq\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"go.uber.org/zap/zapcore\"",
				},
			},
			&ast.ImportSpec{
				Name: &ast.Ident{
					Name: "validation",
				},
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-ozzo/ozzo-validation/v4\"",
				},
			},
		},
	}
}

func (i Errors) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "errs", "errors.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.file()
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
