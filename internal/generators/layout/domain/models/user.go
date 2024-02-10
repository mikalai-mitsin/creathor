package models

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/018bf/creathor/internal/configs"
)

type ModelUser struct {
	project *configs.Project
}

// NewModelUser
// deprecated
func NewModelUser(project *configs.Project) *ModelUser {
	return &ModelUser{project: project}
}

func (m ModelUser) file() *ast.File {
	return &ast.File{
		Name: &ast.Ident{
			Name: "models",
		},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"fmt"`,
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
							Value: fmt.Sprintf(`"%s/internal/errs"`, m.project.Module),
						},
					},
					&ast.ImportSpec{
						Name: &ast.Ident{
							Name: "validation",
						},
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-ozzo/ozzo-validation/v4/is"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"golang.org/x/crypto/bcrypt"`,
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
								Name: "PermissionIDUserList",
							},
						},
						Type: &ast.Ident{
							Name: "PermissionID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"user_list"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "PermissionIDUserDetail",
							},
						},
						Type: &ast.Ident{
							Name: "PermissionID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"user_detail"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "PermissionIDUserCreate",
							},
						},
						Type: &ast.Ident{
							Name: "PermissionID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"user_create"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "PermissionIDUserUpdate",
							},
						},
						Type: &ast.Ident{
							Name: "PermissionID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"user_update"`,
							},
						},
					},
					&ast.ValueSpec{
						Names: []*ast.Ident{
							{
								Name: "PermissionIDUserDelete",
							},
						},
						Type: &ast.Ident{
							Name: "PermissionID",
						},
						Values: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"user_delete"`,
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
							Name: "User",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "ID",
											},
										},
										Type: &ast.Ident{
											Name: "uuid.UUID",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"id,omitempty\"         json:\"id\"         form:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "FirstName",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"first_name\"           json:\"first_name\" form:\"first_name\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "LastName",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"last_name\"            json:\"last_name\"  form:\"last_name\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Password",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"password\"             json:\"-\"          form:\"-\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Email",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"email\"                json:\"email\"      form:\"email\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "GroupID",
											},
										},
										Type: &ast.Ident{
											Name: "GroupID",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"group_id\"             json:\"group_id\"   form:\"group_id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "CreatedAt",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "time",
											},
											Sel: &ast.Ident{
												Name: "Time",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"created_at,omitempty\" json:\"created_at\" form:\"created_at\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "UpdatedAt",
											},
										},
										Type: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "time",
											},
											Sel: &ast.Ident{
												Name: "Time",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`db:\"updated_at\"           json:\"updated_at\" form:\"updated_at\"`",
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
									Name: "User",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "u",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "ID",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "UUID",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "FirstName",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "LastName",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Email",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "EmailFormat",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "CreatedAt",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "UpdatedAt",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
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
						{
							Names: []*ast.Ident{
								{
									Name: "u",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "User",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "SetPassword",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "password",
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
								&ast.Ident{
									Name: "fromPassword",
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
											Name: "bcrypt",
										},
										Sel: &ast.Ident{
											Name: "GenerateFromPassword",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.ArrayType{
												Elt: &ast.Ident{
													Name: "byte",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "password",
												},
											},
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "bcrypt",
											},
											Sel: &ast.Ident{
												Name: "DefaultCost",
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "u",
									},
									Sel: &ast.Ident{
										Name: "Password",
									},
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.Ident{
										Name: "string",
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "fromPassword",
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
									Name: "User",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "CheckPassword",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "password",
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
												Name: "bcrypt",
											},
											Sel: &ast.Ident{
												Name: "CompareHashAndPassword",
											},
										},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.ArrayType{
													Elt: &ast.Ident{
														Name: "byte",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
											},
											&ast.CallExpr{
												Fun: &ast.ArrayType{
													Elt: &ast.Ident{
														Name: "byte",
													},
												},
												Args: []ast.Expr{
													&ast.Ident{
														Name: "password",
													},
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "NewInvalidParameter",
													},
												},
												Args: []ast.Expr{
													&ast.BasicLit{
														Kind:  token.STRING,
														Value: `"email or password"`,
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
									Name: "User",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "FullName",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "fmt",
										},
										Sel: &ast.Ident{
											Name: "Sprintf",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: `"%s %s"`,
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
											},
											Sel: &ast.Ident{
												Name: "FirstName",
											},
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "u",
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
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "UserFilter",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "PageSize",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "uint64",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"page_size\"   form:\"page_size\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "PageNumber",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "uint64",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"page_number\" form:\"page_number\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Search",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"search\"      form:\"search\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "OrderBy",
											},
										},
										Type: &ast.ArrayType{
											Elt: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"order_by\"    form:\"order_by\"`",
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
									Name: "c",
								},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{
									Name: "UserFilter",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "c",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "PageSize",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "PageNumber",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "OrderBy",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "c",
														},
														Sel: &ast.Ident{
															Name: "Search",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
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
							Name: "UserCreate",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "Email",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"email\"    form:\"email\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Password",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"password\" form:\"password\"`",
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
									Name: "UserCreate",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "u",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Email",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "Email",
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "validation",
													},
													Sel: &ast.Ident{
														Name: "Required",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "validation",
													},
													Sel: &ast.Ident{
														Name: "Required",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
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
							Name: "UserUpdate",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "ID",
											},
										},
										Type: &ast.Ident{
											Name: "uuid.UUID",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"id\"         form:\"id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "FirstName",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"first_name\" form:\"first_name\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "LastName",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"last_name\"  form:\"last_name\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Password",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"password\"   form:\"password\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Email",
											},
										},
										Type: &ast.StarExpr{
											X: &ast.Ident{
												Name: "string",
											},
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"email\"      form:\"email\"`",
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
									Name: "UserUpdate",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "u",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "ID",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "validation",
													},
													Sel: &ast.Ident{
														Name: "Required",
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "UUID",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "FirstName",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "LastName",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "validation",
														},
														Sel: &ast.Ident{
															Name: "Length",
														},
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.INT,
															Value: "6",
														},
														&ast.BasicLit{
															Kind:  token.INT,
															Value: "100",
														},
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Email",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "EmailFormat",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
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
							Name: "SetPassword",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											{
												Name: "UserID",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"user_id\"  form:\"user_id\"`",
										},
									},
									{
										Names: []*ast.Ident{
											{
												Name: "Password",
											},
										},
										Type: &ast.Ident{
											Name: "string",
										},
										Tag: &ast.BasicLit{
											Kind:  token.STRING,
											Value: "`json:\"password\" form:\"password\"`",
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
									Name: "SetPassword",
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Validate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
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
						&ast.AssignStmt{
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
											Name: "validation",
										},
										Sel: &ast.Ident{
											Name: "ValidateStruct",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "u",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "UserID",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "validation",
													},
													Sel: &ast.Ident{
														Name: "Required",
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "is",
													},
													Sel: &ast.Ident{
														Name: "UUID",
													},
												},
											},
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "validation",
												},
												Sel: &ast.Ident{
													Name: "Field",
												},
											},
											Args: []ast.Expr{
												&ast.UnaryExpr{
													Op: token.AND,
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "u",
														},
														Sel: &ast.Ident{
															Name: "Password",
														},
													},
												},
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "validation",
													},
													Sel: &ast.Ident{
														Name: "Required",
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "validation",
														},
														Sel: &ast.Ident{
															Name: "Length",
														},
													},
													Args: []ast.Expr{
														&ast.BasicLit{
															Kind:  token.INT,
															Value: "6",
														},
														&ast.BasicLit{
															Kind:  token.INT,
															Value: "100",
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
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "FromValidationError",
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

func (m ModelUser) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "domain", "models", "user.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = m.file()
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
