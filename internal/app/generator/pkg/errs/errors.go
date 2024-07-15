package errs

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/tmpl"

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
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
							Value: "\"reflect\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"go.uber.org/zap/zapcore\"",
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
					&ast.ValueSpec{
						Names: []*ast.Ident{
							&ast.Ident{
								Name: "ErrorCodeClosedRequest",
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
							Name: "Param",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Key",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"key\"`",
										},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Value",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"value\"`",
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
									Name: "p",
								},
							},
							Type: &ast.Ident{
								Name: "Param",
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
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "p",
										},
										Sel: &ast.Ident{
											Name: "Key",
										},
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "p",
										},
										Sel: &ast.Ident{
											Name: "Value",
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
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Params",
						},
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								Name: "Param",
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
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "param",
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
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "param",
													},
													Sel: &ast.Ident{
														Name: "Key",
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "param",
													},
													Sel: &ast.Ident{
														Name: "Value",
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
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "Err",
											},
										},
										Type: &ast.Ident{
											Name: "error",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"-\"`",
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
												Value: &ast.Ident{
													Name: "nil",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Err",
												},
												Value: &ast.Ident{
													Name: "nil",
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
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "err",
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
											Name: "ErrorCodeInternal",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Unexpected behavior.\"",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "err",
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
									&ast.Ident{
										Name: "details",
									},
								},
							},
						},
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
										&ast.Ident{
											Name: "message",
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
					Name: "NewEntityNotFoundError",
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
											Name: "ErrorCodeNotFound",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Entity not found.\"",
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
					Name: "NewBadTokenError",
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
											Name: "ErrorCodePermissionDenied",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Bad token.\"",
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
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "NewError",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "ErrorCodePermissionDenied",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Permission denied.\"",
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
					Name: "NewSubscriptionAlreadyCancelledError",
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
											Name: "ErrorCodeFailedPrecondition",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"Subscription is already cancelled.\"",
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
					Name: "NewInactivePlanError",
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
											Name: "ErrorCodeFailedPrecondition",
										},
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"This plan is inactive.\"",
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
					Name: "Cause",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Err",
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
							Type: &ast.Ident{
								Name: "Error",
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
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"message\"",
									},
									&ast.SelectorExpr{
										X: &ast.Ident{
											Name: "e",
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
										Name: "encoder",
									},
									Sel: &ast.Ident{
										Name: "AddUint",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"code\"",
									},
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "uint",
										},
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
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
												Name: "encoder",
											},
											Sel: &ast.Ident{
												Name: "AddObject",
											},
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: "\"params\"",
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "e",
												},
												Sel: &ast.Ident{
													Name: "Params",
												},
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
					Name: "WithCause",
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
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Err",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "err",
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
								Type: &ast.Ellipsis{
									Ellipsis: 2682,
									Elt: &ast.Ident{
										Name: "Param",
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
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.CallExpr{
									Fun: &ast.Ident{
										Name: "len",
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "e",
											},
											Sel: &ast.Ident{
												Name: "Params",
											},
										},
									},
								},
								Op: token.EQL,
								Y: &ast.BasicLit{
									Kind:  token.INT,
									Value: "0",
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
							Else: &ast.BlockStmt{
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
											&ast.CallExpr{
												Fun: &ast.Ident{
													Name: "append",
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "e",
														},
														Sel: &ast.Ident{
															Name: "Params",
														},
													},
													&ast.Ident{
														Name: "params",
													},
												},
												Ellipsis: 2792,
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
						&ast.DeclStmt{
							Decl: &ast.GenDecl{
								Tok: token.VAR,
								Specs: []ast.Spec{
									&ast.ValueSpec{
										Names: []*ast.Ident{
											&ast.Ident{
												Name: "target",
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
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "ok",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
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
												Name: "tgt",
											},
											&ast.UnaryExpr{
												Op: token.AND,
												X: &ast.Ident{
													Name: "target",
												},
											},
										},
									},
								},
							},
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
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "target",
									},
									Sel: &ast.Ident{
										Name: "Code",
									},
								},
								Op: token.NEQ,
								Y: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Code",
									},
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
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "target",
									},
									Sel: &ast.Ident{
										Name: "Message",
									},
								},
								Op: token.NEQ,
								Y: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "e",
									},
									Sel: &ast.Ident{
										Name: "Message",
									},
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
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "param",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "target",
								},
								Sel: &ast.Ident{
									Name: "Params",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.UnaryExpr{
											Op: token.NOT,
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "slices",
													},
													Sel: &ast.Ident{
														Name: "Contains",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "e",
														},
														Sel: &ast.Ident{
															Name: "Params",
														},
													},
													&ast.Ident{
														Name: "param",
													},
												},
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
								},
							},
						},
						&ast.RangeStmt{
							Key: &ast.Ident{
								Name: "_",
							},
							Value: &ast.Ident{
								Name: "param",
							},
							Tok: token.DEFINE,
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "e",
								},
								Sel: &ast.Ident{
									Name: "Params",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.UnaryExpr{
											Op: token.NOT,
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "slices",
													},
													Sel: &ast.Ident{
														Name: "Contains",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "target",
														},
														Sel: &ast.Ident{
															Name: "Params",
														},
													},
													&ast.Ident{
														Name: "param",
													},
												},
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
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.Ident{
									Name: "true",
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
					Name: "SetCause",
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
										Name: "Err",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.Ident{
									Name: "err",
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
								Type: &ast.Ident{
									Name: "Params",
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
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "append",
									},
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "e",
											},
											Sel: &ast.Ident{
												Name: "Params",
											},
										},
										&ast.CompositeLit{
											Type: &ast.Ident{
												Name: "Param",
											},
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key: &ast.Ident{
														Name: "Key",
													},
													Value: &ast.Ident{
														Name: "key",
													},
												},
												&ast.KeyValueExpr{
													Key: &ast.Ident{
														Name: "Value",
													},
													Value: &ast.Ident{
														Name: "value",
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
					Value: "\"reflect\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"go.uber.org/zap/zapcore\"",
				},
			},
		},
	}
}

func (i Errors) fileHttp() *ast.File {
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
							Value: "\"encoding/json\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"net/http\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "ClientClosedRequest",
							},
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.INT,
								Value: "499",
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "GetHTTPStatus",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
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
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "int",
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.SwitchStmt{
							Tag: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "e",
								},
								Sel: &ast.Ident{
									Name: "Code",
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeOK",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusOK",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeCanceled",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.Ident{
														Name: "ClientClosedRequest",
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeUnknown",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeInvalidArgument",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusBadRequest",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeDeadlineExceeded",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeNotFound",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusNotFound",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeAlreadyExists",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusBadRequest",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodePermissionDenied",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusForbidden",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeResourceExhausted",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeFailedPrecondition",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusBadRequest",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeAborted",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeOutOfRange",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeUnimplemented",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusMethodNotAllowed",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeInternal",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeUnavailable",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusServiceUnavailable",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeDataLoss",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											&ast.Ident{
												Name: "ErrorCodeUnauthenticated",
											},
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusUnauthorized",
														},
													},
												},
											},
										},
									},
									&ast.CaseClause{
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "http",
														},
														Sel: &ast.Ident{
															Name: "StatusInternalServerError",
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
					Name: "RenderToHTTPResponse",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "e",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{
										Name: "Error",
									},
								},
							},
							{
								Names: []*ast.Ident{
									{
										Name: "writer",
									},
								},
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "http",
									},
									Sel: &ast.Ident{
										Name: "ResponseWriter",
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "writer",
									},
									Sel: &ast.Ident{
										Name: "WriteHeader",
									},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.Ident{
											Name: "GetHTTPStatus",
										},
										Args: []ast.Expr{
											&ast.Ident{
												Name: "e",
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
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "json",
												},
												Sel: &ast.Ident{
													Name: "NewEncoder",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "writer",
												},
											},
										},
										Sel: &ast.Ident{
											Name: "Encode",
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
					},
				},
			},
		},
		Imports: []*ast.ImportSpec{
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"encoding/json\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"net/http\"",
				},
			},
		},
	}

}

func (i Errors) fileGrpc() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "errs",
		},
	}
}

func (i Errors) filePostgres() *ast.File {
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
							Value: "\"database/sql\"",
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
							Value: "\"github.com/lib/pq\"",
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "sqlConflictCode",
							},
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: "\"23505\"",
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
							{
								Names: []*ast.Ident{
									{
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
							{
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
												Value: &ast.Ident{
													Name: "nil",
												},
											},
											&ast.KeyValueExpr{
												Key: &ast.Ident{
													Name: "Err",
												},
												Value: &ast.Ident{
													Name: "err",
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
											{
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
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "pqErr",
												},
												Sel: &ast.Ident{
													Name: "Code",
												},
											},
											Op: token.EQL,
											Y: &ast.Ident{
												Name: "sqlConflictCode",
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
															Fun: &ast.SelectorExpr{
																X: &ast.CallExpr{
																	Fun: &ast.Ident{
																		Name: "NewInvalidFormError",
																	},
																},
																Sel: &ast.Ident{
																	Name: "WithCause",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "err",
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
										Value: "\"error\"",
									},
									&ast.CallExpr{
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
													Name: "NewEntityNotFoundError",
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
		},
		Imports: []*ast.ImportSpec{
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"database/sql\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"errors\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"fmt\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/lib/pq\"",
				},
			},
		},
	}
}

func (i Errors) fileValidation() *ast.File {
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
							Value: "\"errors\"",
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
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"text/template\"",
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "NewFromValidationError",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
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
							{
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
											{
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
											{
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "e",
													},
													Sel: &ast.Ident{
														Name: "WithCause",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "validationErrors",
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
												Fun: &ast.SelectorExpr{
													X: &ast.CallExpr{
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
													Sel: &ast.Ident{
														Name: "WithCause",
													},
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
							{
								Names: []*ast.Ident{
									{
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
							{
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
											{
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
		},
		Imports: []*ast.ImportSpec{
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"bytes\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"errors\"",
				},
			},
			{
				Name: &ast.Ident{
					Name: "validation",
				},
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-ozzo/ozzo-validation/v4\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"text/template\"",
				},
			},
		},
	}
}

var destinationPath = "."

func (i Errors) Sync() error {
	if err := i.syncErrs(); err != nil {
		return err
	}
	if err := i.syncHttp(); err != nil {
		return err
	}
	if err := i.syncGrpc(); err != nil {
		return err
	}
	if err := i.syncPostgres(); err != nil {
		return err
	}
	if err := i.syncValidation(); err != nil {
		return err
	}
	return nil
}

func (i Errors) syncErrs() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "errors.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
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

	test := &tmpl.Template{
		SourcePath: "templates/internal/pkg/errs/errors_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"pkg",
			"errs",
			"errors_test.go",
		),
		Name: "domain errors tests",
	}
	if err := test.RenderToFile(i.project); err != nil {
		return err
	}
	return nil
}

func (i Errors) syncGrpc() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "grpc.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.fileGrpc()
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}

	test := &tmpl.Template{
		SourcePath: "templates/internal/pkg/errs/grpc_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"pkg",
			"errs",
			"grpc_test.go",
		),
		Name: "grpc errors tests",
	}
	if err := test.RenderToFile(i.project); err != nil {
		return err
	}
	return nil
}

func (i Errors) syncHttp() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "http.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.fileHttp()
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}

	test := &tmpl.Template{
		SourcePath: "templates/internal/pkg/errs/http_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"pkg",
			"errs",
			"http_test.go",
		),
		Name: "http errors tests",
	}
	if err := test.RenderToFile(i.project); err != nil {
		return err
	}
	return nil
}

func (i Errors) syncPostgres() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "postgres.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.filePostgres()
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}

	test := &tmpl.Template{
		SourcePath: "templates/internal/pkg/errs/postgres_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"pkg",
			"errs",
			"postgres_test.go",
		),
		Name: "postgres errors tests",
	}
	if err := test.RenderToFile(i.project); err != nil {
		return err
	}
	return nil
}

func (i Errors) syncValidation() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "validation.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.fileValidation()
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(filename, buff.Bytes(), 0777); err != nil {
		return err
	}

	test := &tmpl.Template{
		SourcePath: "templates/internal/pkg/errs/validation_test.go.tmpl",
		DestinationPath: path.Join(
			destinationPath,
			"internal",
			"pkg",
			"errs",
			"validation_test.go",
		),
		Name: "validation errors tests",
	}
	if err := test.RenderToFile(i.project); err != nil {
		return err
	}
	return nil
}
