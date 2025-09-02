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

type Generator struct {
	project *configs.Project
}

func NewGenerator(project *configs.Project) *Generator {
	return &Generator{
		project: project,
	}
}

func (i Generator) file() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("errs"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"encoding/json"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"errors"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"slices"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"go.uber.org/zap/zapcore"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("ErrorCode"),
						Type: ast.NewIdent("uint"),
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeOK"),
						},
						Type: ast.NewIdent("ErrorCode"),
						Values: []ast.Expr{
							ast.NewIdent("iota"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeCanceled"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeUnknown"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeInvalidArgument"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeDeadlineExceeded"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeNotFound"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeAlreadyExists"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodePermissionDenied"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeResourceExhausted"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeFailedPrecondition"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeAborted"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeOutOfRange"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeUnimplemented"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeInternal"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeUnavailable"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeDataLoss"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeUnauthenticated"),
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("ErrorCodeClosedRequest"),
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Param"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Key"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"key\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Value"),
										},
										Type: ast.NewIdent("string"),
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
						{
							Names: []*ast.Ident{
								ast.NewIdent("p"),
							},
							Type: ast.NewIdent("Param"),
						},
					},
				},
				Name: ast.NewIdent("MarshalLogObject"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("encoder"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("zapcore"),
									Sel: ast.NewIdent("ObjectEncoder"),
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("encoder"),
									Sel: ast.NewIdent("AddString"),
								},
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("p"),
										Sel: ast.NewIdent("Key"),
									},
									&ast.SelectorExpr{
										X:   ast.NewIdent("p"),
										Sel: ast.NewIdent("Value"),
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
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Params"),
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("Param"),
						},
					},
				},
			},
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("p"),
							},
							Type: ast.NewIdent("Params"),
						},
					},
				},
				Name: ast.NewIdent("MarshalLogObject"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("encoder"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("zapcore"),
									Sel: ast.NewIdent("ObjectEncoder"),
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
						&ast.RangeStmt{
							Key:   ast.NewIdent("_"),
							Value: ast.NewIdent("param"),
							Tok:   token.DEFINE,
							X:     ast.NewIdent("p"),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("encoder"),
												Sel: ast.NewIdent("AddString"),
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X:   ast.NewIdent("param"),
													Sel: ast.NewIdent("Key"),
												},
												&ast.SelectorExpr{
													X:   ast.NewIdent("param"),
													Sel: ast.NewIdent("Value"),
												},
											},
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
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent("Error"),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent("Code"),
										},
										Type: ast.NewIdent("ErrorCode"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"code\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Message"),
										},
										Type: ast.NewIdent("string"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"message\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Params"),
										},
										Type: ast.NewIdent("Params"),
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"params\"`",
										},
									},
									{
										Names: []*ast.Ident{
											ast.NewIdent("Err"),
										},
										Type: ast.NewIdent("error"),
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
				Name: ast.NewIdent("NewError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("code"),
								},
								Type: ast.NewIdent("ErrorCode"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("message"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
										Type: ast.NewIdent("Error"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Code"),
												Value: ast.NewIdent("code"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Message"),
												Value: ast.NewIdent("message"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Params"),
												Value: ast.NewIdent("nil"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Err"),
												Value: ast.NewIdent("nil"),
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
				Name: ast.NewIdent("NewUnexpectedBehaviorError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("details"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodeInternal"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"Unexpected behavior."`,
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("err"),
									Sel: ast.NewIdent("AddParam"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"details"`,
									},
									ast.NewIdent("details"),
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("err"),
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewInvalidFormError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodeInvalidArgument"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"The form sent is not valid, please correct the errors below."`,
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewInvalidParameter"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("message"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodeInvalidArgument"),
										ast.NewIdent("message"),
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewEntityNotFoundError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodeNotFound"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"Name not found."`,
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewBadTokenError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodePermissionDenied"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"Bad token."`,
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewPermissionDeniedError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodePermissionDenied"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"Permission denied."`,
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewUnauthenticatedError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("NewError"),
									Args: []ast.Expr{
										ast.NewIdent("ErrorCodeFailedPrecondition"),
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"Unauthenticated error."`,
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("Cause"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Err"),
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
								ast.NewIdent("e"),
							},
							Type: ast.NewIdent("Error"),
						},
					},
				},
				Name: ast.NewIdent("MarshalLogObject"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("encoder"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("zapcore"),
									Sel: ast.NewIdent("ObjectEncoder"),
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("encoder"),
									Sel: ast.NewIdent("AddString"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"message"`,
									},
									&ast.SelectorExpr{
										X:   ast.NewIdent("e"),
										Sel: ast.NewIdent("Message"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("encoder"),
									Sel: ast.NewIdent("AddUint"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"code"`,
									},
									&ast.CallExpr{
										Fun: ast.NewIdent("uint"),
										Args: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("Code"),
											},
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
											X:   ast.NewIdent("encoder"),
											Sel: ast.NewIdent("AddObject"),
										},
										Args: []ast.Expr{
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `"params"`,
											},
											&ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("Params"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("WithParam"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("key"),
									ast.NewIdent("value"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("AddParam"),
								},
								Args: []ast.Expr{
									ast.NewIdent("key"),
									ast.NewIdent("value"),
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("e"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("WithCause"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Err"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("err"),
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("e"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("WithParams"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("params"),
								},
								Type: &ast.Ellipsis{
									Ellipsis: 2531,
									Elt:      ast.NewIdent("Param"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
									Fun: ast.NewIdent("len"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("e"),
											Sel: ast.NewIdent("Params"),
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
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("Params"),
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											ast.NewIdent("params"),
										},
									},
								},
							},
							Else: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											&ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("Params"),
											},
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("append"),
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("e"),
														Sel: ast.NewIdent("Params"),
													},
													ast.NewIdent("params"),
												},
												Ellipsis: 2641,
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("e"),
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
								ast.NewIdent("e"),
							},
							Type: ast.NewIdent("Error"),
						},
					},
				},
				Name: ast.NewIdent("Error"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("data"),
								ast.NewIdent("_"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("json"),
										Sel: ast.NewIdent("Marshal"),
									},
									Args: []ast.Expr{
										ast.NewIdent("e"),
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("string"),
									Args: []ast.Expr{
										ast.NewIdent("data"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("Is"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("tgt"),
								},
								Type: ast.NewIdent("error"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("bool"),
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
											ast.NewIdent("target"),
										},
										Type: &ast.StarExpr{
											X: ast.NewIdent("Error"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Init: &ast.AssignStmt{
								Lhs: []ast.Expr{
									ast.NewIdent("ok"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("errors"),
											Sel: ast.NewIdent("As"),
										},
										Args: []ast.Expr{
											ast.NewIdent("tgt"),
											&ast.UnaryExpr{
												Op: token.AND,
												X:  ast.NewIdent("target"),
											},
										},
									},
								},
							},
							Cond: &ast.UnaryExpr{
								Op: token.NOT,
								X:  ast.NewIdent("ok"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("false"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("target"),
									Sel: ast.NewIdent("Code"),
								},
								Op: token.NEQ,
								Y: &ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Code"),
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("false"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.BinaryExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("target"),
									Sel: ast.NewIdent("Message"),
								},
								Op: token.NEQ,
								Y: &ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Message"),
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ReturnStmt{
										Results: []ast.Expr{
											ast.NewIdent("false"),
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key:   ast.NewIdent("_"),
							Value: ast.NewIdent("param"),
							Tok:   token.DEFINE,
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("target"),
								Sel: ast.NewIdent("Params"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.UnaryExpr{
											Op: token.NOT,
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("slices"),
													Sel: ast.NewIdent("Contains"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("e"),
														Sel: ast.NewIdent("Params"),
													},
													ast.NewIdent("param"),
												},
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.ReturnStmt{
													Results: []ast.Expr{
														ast.NewIdent("false"),
													},
												},
											},
										},
									},
								},
							},
						},
						&ast.RangeStmt{
							Key:   ast.NewIdent("_"),
							Value: ast.NewIdent("param"),
							Tok:   token.DEFINE,
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("e"),
								Sel: ast.NewIdent("Params"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.IfStmt{
										Cond: &ast.UnaryExpr{
											Op: token.NOT,
											X: &ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("slices"),
													Sel: ast.NewIdent("Contains"),
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("target"),
														Sel: ast.NewIdent("Params"),
													},
													ast.NewIdent("param"),
												},
											},
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.ReturnStmt{
													Results: []ast.Expr{
														ast.NewIdent("false"),
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
								ast.NewIdent("true"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("SetCode"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("code"),
								},
								Type: ast.NewIdent("ErrorCode"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Code"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("code"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("SetCause"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Err"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("err"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("SetMessage"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("message"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Message"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("message"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("SetParams"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("params"),
								},
								Type: ast.NewIdent("Params"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Params"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								ast.NewIdent("params"),
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
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("AddParam"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("key"),
									ast.NewIdent("value"),
								},
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("Params"),
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent("append"),
									Args: []ast.Expr{
										&ast.SelectorExpr{
											X:   ast.NewIdent("e"),
											Sel: ast.NewIdent("Params"),
										},
										&ast.CompositeLit{
											Type: ast.NewIdent("Param"),
											Elts: []ast.Expr{
												&ast.KeyValueExpr{
													Key:   ast.NewIdent("Key"),
													Value: ast.NewIdent("key"),
												},
												&ast.KeyValueExpr{
													Key:   ast.NewIdent("Value"),
													Value: ast.NewIdent("value"),
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
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"encoding/json"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"errors"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"reflect"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"go.uber.org/zap/zapcore"`,
				},
			},
		},
	}
}

func (i Generator) fileHttp() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("errs"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"errors\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/go-chi/render\"",
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
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("GetHTTPStatus"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("int"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.SwitchStmt{
							Tag: &ast.SelectorExpr{
								X:   ast.NewIdent("e"),
								Sel: ast.NewIdent("Code"),
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeOK"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusOK"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeUnknown"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeInvalidArgument"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusBadRequest"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeDeadlineExceeded"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeNotFound"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusNotFound"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeAlreadyExists"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusBadRequest"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodePermissionDenied"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusForbidden"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeResourceExhausted"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeFailedPrecondition"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusBadRequest"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeAborted"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeOutOfRange"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeUnimplemented"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusMethodNotAllowed"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeInternal"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeUnavailable"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusServiceUnavailable"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeDataLoss"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
													},
												},
											},
										},
									},
									&ast.CaseClause{
										List: []ast.Expr{
											ast.NewIdent("ErrorCodeUnauthenticated"),
										},
										Body: []ast.Stmt{
											&ast.ReturnStmt{
												Results: []ast.Expr{
													&ast.SelectorExpr{
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusUnauthorized"),
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
														X:   ast.NewIdent("http"),
														Sel: ast.NewIdent("StatusInternalServerError"),
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
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								ast.NewIdent("e"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent("Error"),
							},
						},
					},
				},
				Name: ast.NewIdent("Render"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("w"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("r"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("Request"),
									},
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("Status"),
								},
								Args: []ast.Expr{
									ast.NewIdent("r"),
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("e"),
											Sel: ast.NewIdent("GetHTTPStatus"),
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("JSON"),
								},
								Args: []ast.Expr{
									ast.NewIdent("w"),
									ast.NewIdent("r"),
									ast.NewIdent("e"),
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
				Name: ast.NewIdent("RenderToHTTPResponse"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("w"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("http"),
									Sel: ast.NewIdent("ResponseWriter"),
								},
							},
							{
								Names: []*ast.Ident{
									ast.NewIdent("r"),
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("Request"),
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
											ast.NewIdent("e"),
										},
										Type: &ast.StarExpr{
											X: ast.NewIdent("Error"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("As"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("e"),
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
														X:   ast.NewIdent("render"),
														Sel: ast.NewIdent("Render"),
													},
													Args: []ast.Expr{
														ast.NewIdent("w"),
														ast.NewIdent("r"),
														ast.NewIdent("e"),
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
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("render"),
															Sel: ast.NewIdent("Status"),
														},
														Args: []ast.Expr{
															ast.NewIdent("r"),
															&ast.SelectorExpr{
																X:   ast.NewIdent("http"),
																Sel: ast.NewIdent("StatusInternalServerError"),
															},
														},
													},
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("Status"),
								},
								Args: []ast.Expr{
									ast.NewIdent("r"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("http"),
										Sel: ast.NewIdent("StatusInternalServerError"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("PlainText"),
								},
								Args: []ast.Expr{
									ast.NewIdent("w"),
									ast.NewIdent("r"),
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("err"),
											Sel: ast.NewIdent("Error"),
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
					Value: "\"errors\"",
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/go-chi/render\"",
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

func (i Generator) fileGrpc() *ast.File {
	return &ast.File{
		Name: ast.NewIdent("errs"),
	}
}

func (i Generator) filePostgres() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("errs"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"database/sql"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"errors"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"fmt"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/lib/pq"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.CONST,
				Specs: []ast.Spec{
					&ast.ValueSpec{
						Names: []*ast.Ident{
							ast.NewIdent("sqlConflictCode"),
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"23505"`,
							},
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("FromPostgresError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("e"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.UnaryExpr{
									Op: token.AND,
									X: &ast.CompositeLit{
										Type: ast.NewIdent("Error"),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Code"),
												Value: ast.NewIdent("ErrorCodeInternal"),
											},
											&ast.KeyValueExpr{
												Key: ast.NewIdent("Message"),
												Value: &ast.BasicLit{
													Kind:  token.STRING,
													Value: `"Unexpected behavior."`,
												},
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Params"),
												Value: ast.NewIdent("nil"),
											},
											&ast.KeyValueExpr{
												Key:   ast.NewIdent("Err"),
												Value: ast.NewIdent("err"),
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
											ast.NewIdent("pqErr"),
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("pq"),
												Sel: ast.NewIdent("Error"),
											},
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("As"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("pqErr"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("AddParam"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"details"`,
												},
												&ast.SelectorExpr{
													X:   ast.NewIdent("pqErr"),
													Sel: ast.NewIdent("Detail"),
												},
											},
										},
									},
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("AddParam"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"message"`,
												},
												&ast.SelectorExpr{
													X:   ast.NewIdent("pqErr"),
													Sel: ast.NewIdent("Message"),
												},
											},
										},
									},
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("e"),
												Sel: ast.NewIdent("AddParam"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"postgres_code"`,
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X:   ast.NewIdent("fmt"),
														Sel: ast.NewIdent("Sprint"),
													},
													Args: []ast.Expr{
														&ast.SelectorExpr{
															X:   ast.NewIdent("pqErr"),
															Sel: ast.NewIdent("Code"),
														},
													},
												},
											},
										},
									},
									&ast.IfStmt{
										Cond: &ast.BinaryExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("pqErr"),
												Sel: ast.NewIdent("Code"),
											},
											Op: token.EQL,
											Y:  ast.NewIdent("sqlConflictCode"),
										},
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														ast.NewIdent("e"),
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.CallExpr{
																	Fun: ast.NewIdent("NewInvalidFormError"),
																},
																Sel: ast.NewIdent("WithCause"),
															},
															Args: []ast.Expr{
																ast.NewIdent("err"),
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
									X:   ast.NewIdent("e"),
									Sel: ast.NewIdent("AddParam"),
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: `"error"`,
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X:   ast.NewIdent("err"),
											Sel: ast.NewIdent("Error"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("Is"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.SelectorExpr{
										X:   ast.NewIdent("sql"),
										Sel: ast.NewIdent("ErrNoRows"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("e"),
										},
										Tok: token.ASSIGN,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("NewEntityNotFoundError"),
											},
										},
									},
								},
							},
						},
						&ast.ReturnStmt{
							Results: []ast.Expr{
								ast.NewIdent("e"),
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
					Value: `"database/sql"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"errors"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"fmt"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/lib/pq"`,
				},
			},
		},
	}
}
func (i Generator) fileKafka() *ast.File {
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
							Value: "\"errors\"",
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: "\"github.com/IBM/sarama\"",
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "FromKafkaError",
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
											&ast.Ident{
												Name: "prErr",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "sarama",
												},
												Sel: &ast.Ident{
													Name: "ProducerError",
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
											Name: "prErr",
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
													Value: "\"error\"",
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "prErr",
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
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"errors\"",
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: "\"github.com/IBM/sarama\"",
				},
			},
		},
	}
}

func (i Generator) fileValidation() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("errs"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"bytes"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"errors"`,
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("validation"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"text/template"`,
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("NewFromValidationError"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("err"),
								},
								Type: ast.NewIdent("error"),
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: ast.NewIdent("Error"),
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
											ast.NewIdent("validationErrors"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("validation"),
											Sel: ast.NewIdent("Errors"),
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
											ast.NewIdent("validationErrorObject"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("validation"),
											Sel: ast.NewIdent("ErrorObject"),
										},
									},
								},
							},
						},
						&ast.IfStmt{
							Cond: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("As"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("validationErrors"),
									},
								},
							},
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.AssignStmt{
										Lhs: []ast.Expr{
											ast.NewIdent("e"),
										},
										Tok: token.DEFINE,
										Rhs: []ast.Expr{
											&ast.CallExpr{
												Fun: ast.NewIdent("NewError"),
												Args: []ast.Expr{
													ast.NewIdent("ErrorCodeInvalidArgument"),
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"The form sent is not valid, please correct the errors below."`,
													},
												},
											},
										},
									},
									&ast.RangeStmt{
										Key:   ast.NewIdent("key"),
										Value: ast.NewIdent("value"),
										Tok:   token.DEFINE,
										X:     ast.NewIdent("validationErrors"),
										Body: &ast.BlockStmt{
											List: []ast.Stmt{
												&ast.TypeSwitchStmt{
													Assign: &ast.AssignStmt{
														Lhs: []ast.Expr{
															ast.NewIdent("t"),
														},
														Tok: token.DEFINE,
														Rhs: []ast.Expr{
															&ast.TypeAssertExpr{
																X: ast.NewIdent("value"),
															},
														},
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.CaseClause{
																List: []ast.Expr{
																	&ast.SelectorExpr{
																		X:   ast.NewIdent("validation"),
																		Sel: ast.NewIdent("ErrorObject"),
																	},
																},
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X:   ast.NewIdent("e"),
																				Sel: ast.NewIdent("AddParam"),
																			},
																			Args: []ast.Expr{
																				ast.NewIdent("key"),
																				&ast.CallExpr{
																					Fun: ast.NewIdent("renderErrorMessage"),
																					Args: []ast.Expr{
																						ast.NewIdent("t"),
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
																		X: ast.NewIdent("Error"),
																	},
																},
																Body: []ast.Stmt{
																	&ast.ExprStmt{
																		X: &ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X:   ast.NewIdent("e"),
																				Sel: ast.NewIdent("AddParam"),
																			},
																			Args: []ast.Expr{
																				ast.NewIdent("key"),
																				&ast.SelectorExpr{
																					X:   ast.NewIdent("t"),
																					Sel: ast.NewIdent("Message"),
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
																				X:   ast.NewIdent("e"),
																				Sel: ast.NewIdent("AddParam"),
																			},
																			Args: []ast.Expr{
																				ast.NewIdent("key"),
																				&ast.CallExpr{
																					Fun: &ast.SelectorExpr{
																						X:   ast.NewIdent("value"),
																						Sel: ast.NewIdent("Error"),
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
													X:   ast.NewIdent("e"),
													Sel: ast.NewIdent("WithCause"),
												},
												Args: []ast.Expr{
													ast.NewIdent("validationErrors"),
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
									X:   ast.NewIdent("errors"),
									Sel: ast.NewIdent("As"),
								},
								Args: []ast.Expr{
									ast.NewIdent("err"),
									&ast.UnaryExpr{
										Op: token.AND,
										X:  ast.NewIdent("validationErrorObject"),
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
														Fun: ast.NewIdent("NewInvalidParameter"),
														Args: []ast.Expr{
															&ast.CallExpr{
																Fun: ast.NewIdent("renderErrorMessage"),
																Args: []ast.Expr{
																	ast.NewIdent("validationErrorObject"),
																},
															},
														},
													},
													Sel: ast.NewIdent("WithCause"),
												},
												Args: []ast.Expr{
													ast.NewIdent("validationErrorObject"),
												},
											},
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
				Name: ast.NewIdent("renderErrorMessage"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent("object"),
								},
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("validation"),
									Sel: ast.NewIdent("ErrorObject"),
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: ast.NewIdent("string"),
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("parse"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("template"),
												Sel: ast.NewIdent("New"),
											},
											Args: []ast.Expr{
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"message"`,
												},
											},
										},
										Sel: ast.NewIdent("Parse"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("object"),
												Sel: ast.NewIdent("Message"),
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
											&ast.BasicLit{
												Kind:  token.STRING,
												Value: `""`,
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
											ast.NewIdent("tpl"),
										},
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("bytes"),
											Sel: ast.NewIdent("Buffer"),
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("_"),
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("parse"),
										Sel: ast.NewIdent("Execute"),
									},
									Args: []ast.Expr{
										&ast.UnaryExpr{
											Op: token.AND,
											X:  ast.NewIdent("tpl"),
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("object"),
												Sel: ast.NewIdent("Params"),
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
										X:   ast.NewIdent("tpl"),
										Sel: ast.NewIdent("String"),
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
					Value: `"bytes"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"errors"`,
				},
			},
			{
				Name: ast.NewIdent("validation"),
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
				},
			},
			{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"text/template"`,
				},
			},
		},
	}
}

var destinationPath = "."

func (i Generator) Sync() error {
	if err := i.syncErrs(); err != nil {
		return err
	}
	if i.project.GRPCEnabled {
		if err := i.syncGrpc(); err != nil {
			return err
		}
	}
	if i.project.HTTPEnabled {
		if err := i.syncHttp(); err != nil {
			return err
		}
	}
	if err := i.syncPostgres(); err != nil {
		return err
	}
	if i.project.KafkaEnabled {
		if err := i.syncKafka(); err != nil {
			return err
		}
	}
	if err := i.syncValidation(); err != nil {
		return err
	}
	return nil
}

func (i Generator) syncErrs() error {
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

func (i Generator) syncGrpc() error {
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

func (i Generator) syncHttp() error {
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

func (i Generator) syncPostgres() error {
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

func (i Generator) syncKafka() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "pkg", "errs", "kafka.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = i.fileKafka()
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

func (i Generator) syncValidation() error {
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
