package http

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

const destinationPath = "."

type HandlerGenerator struct {
	domain *domain.Domain
}

func NewHandlerGenerator(domain *domain.Domain) *HandlerGenerator {
	return &HandlerGenerator{
		domain: domain,
	}
}

func (h *HandlerGenerator) Sync() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	if err := os.MkdirAll(path.Dir(filename), 0777); err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
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

func (h *HandlerGenerator) filename() string {
	return path.Join("internal", "app", h.domain.DirName(), "handlers", "http", h.domain.FileName())
}

func (h *HandlerGenerator) file() *ast.File {
	return &ast.File{
		Package: 1,
		Name:    ast.NewIdent("handlers"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-chi/chi/v5"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/go-chi/render"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/errs"`, h.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, h.domain.Module),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"net/http"`,
						},
					},
				},
			},
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
										},
										Type: ast.NewIdent(h.domain.GetUseCaseInterfaceName()),
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
				Name: ast.NewIdent(h.domain.GetHTTPHandlerConstructorName()),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Type: ast.NewIdent(h.domain.GetUseCaseInterfaceName()),
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
								Type: &ast.StarExpr{
									X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
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
										Type: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: ast.NewIdent(
													h.domain.GetUseCasePrivateVariableName(),
												),
												Value: ast.NewIdent(
													h.domain.GetUseCasePrivateVariableName(),
												),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Create",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @Summary Create %s", h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Security BearerAuth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: fmt.Sprintf("// @Param form body %s true \"Create %s request\"", h.domain.GetHTTPCreateDTOName(), h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Success 201 {object} %s \"Created %s\"", h.domain.GetHTTPItemDTOName(), h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: fmt.Sprintf("// @Router /api/v1/%s/ [POST]", h.domain.GetManyVariableName()),
						},
					},
				},
				Name: ast.NewIdent("Create"),
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
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("createDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPCreateDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent("r"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("create"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("createDTO"),
										Sel: ast.NewIdent("toEntity"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(h.domain.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.domain.GetUseCasePrivateVariableName(),
											),
										},
										Sel: ast.NewIdent("Create"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("create"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.domain.GetOneVariableName()),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
										Sel: ast.NewIdent("StatusCreated"),
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
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Get",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @Summary Get %s by id", h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Security BearerAuth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: "// @Param id path string true \"UUID\"",
						},
						{
							Text: fmt.Sprintf("// @Success 200 {object} %s \"Requested %s\"", h.domain.GetHTTPItemDTOName(), h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: fmt.Sprintf("// @Router /api/v1/%s/{id}/ [GET]", h.domain.GetManyVariableName()),
						},
					},
				},
				Name: ast.NewIdent("Get"),
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
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("id"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("uuid"),
										Sel: ast.NewIdent("UUID"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("chi"),
												Sel: ast.NewIdent("URLParam"),
											},
											Args: []ast.Expr{
												ast.NewIdent("r"),
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"id"`,
												},
											},
										},
									},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(h.domain.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.domain.GetUseCasePrivateVariableName(),
											),
										},
										Sel: ast.NewIdent("Get"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("id"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.domain.GetOneVariableName()),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
										Sel: ast.NewIdent("StatusOK"),
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
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// List",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @Summary List of %s", h.domain.GetManyVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Security BearerAuth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: fmt.Sprintf("// @Param filter query %s true \"Filter of %s\"", h.domain.GetHTTPFilterDTOName(), h.domain.GetManyVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Success 200 {array} %s \"Filtered list of %s\"", h.domain.GetHTTPListDTOName(), h.domain.GetManyVariableName()),
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: fmt.Sprintf("// @Router /api/v1/%s/ [GET]", h.domain.GetManyVariableName()),
						},
					},
				},
				Name: ast.NewIdent("List"),
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
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("filterDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPFilterDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent("r"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("filter"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("filterDTO"),
										Sel: ast.NewIdent("toEntity"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(h.domain.GetManyVariableName()),
								ast.NewIdent("count"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.domain.GetUseCasePrivateVariableName(),
											),
										},
										Sel: ast.NewIdent("List"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("filter"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPListDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.domain.GetManyVariableName()),
										ast.NewIdent("count"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
										Sel: ast.NewIdent("StatusOK"),
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
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Update",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @Summary Update %s", h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Security BearerAuth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: "// @Param id path string true \"UUID\"",
						},
						{
							Text: fmt.Sprintf("// @Param form body %s true \"Update %s request\"", h.domain.GetHTTPUpdateDTOName(), h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Success 200 {object} %s \"Updated %s\"", h.domain.GetHTTPItemDTOName(), h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: fmt.Sprintf("// @Router /api/v1/%s/{id}/ [PATCH]", h.domain.GetManyVariableName()),
						},
					},
				},
				Name: ast.NewIdent("Update"),
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
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("updateDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPUpdateDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent("r"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("update"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("updateDTO"),
										Sel: ast.NewIdent("toEntity"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent(h.domain.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.domain.GetUseCasePrivateVariableName(),
											),
										},
										Sel: ast.NewIdent("Update"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("r"),
												Sel: ast.NewIdent("Context"),
											},
										},
										ast.NewIdent("update"),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("response"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.domain.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.domain.GetOneVariableName()),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
										Sel: ast.NewIdent("StatusOK"),
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
									ast.NewIdent("response"),
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
								ast.NewIdent("h"),
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: "// Delete",
						},
						{
							Text: "//",
						},
						{
							Text: fmt.Sprintf("// @Summary Delete %s by id", h.domain.GetOneVariableName()),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.domain.GetOneVariableName()),
						},
						{
							Text: "// @Security BearerAuth",
						},
						{
							Text: "// @Accept json",
						},
						{
							Text: "// @Produce json",
						},
						{
							Text: "// @Param id path string true \"UUID\"",
						},
						{
							Text: "// @Success 204 \"No content\"",
						},
						{
							Text: "// @Failure 400 {object} errs.Error \"Invalid request body or validation error\"",
						},
						{
							Text: "// @Failure 401 {object} errs.Error \"Unauthorized\"",
						},
						{
							Text: "// @Failure 404 {object} errs.Error \"Not found\"",
						},
						{
							Text: "// @Failure 500 {object} errs.Error \"Internal server error\"",
						},
						{
							Text: fmt.Sprintf("// @Router /api/v1/%s/{id}/ [DELETE]", h.domain.GetManyVariableName()),
						},
					},
				},
				Name: ast.NewIdent("Delete"),
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
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("id"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("uuid"),
										Sel: ast.NewIdent("UUID"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X:   ast.NewIdent("chi"),
												Sel: ast.NewIdent("URLParam"),
											},
											Args: []ast.Expr{
												ast.NewIdent("r"),
												&ast.BasicLit{
													Kind:  token.STRING,
													Value: `"id"`,
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
									ast.NewIdent("err"),
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.SelectorExpr{
												X: ast.NewIdent("h"),
												Sel: ast.NewIdent(
													h.domain.GetUseCasePrivateVariableName(),
												),
											},
											Sel: ast.NewIdent("Delete"),
										},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X:   ast.NewIdent("r"),
													Sel: ast.NewIdent("Context"),
												},
											},
											ast.NewIdent("id"),
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
												X:   ast.NewIdent("errs"),
												Sel: ast.NewIdent("RenderToHTTPResponse"),
											},
											Args: []ast.Expr{
												ast.NewIdent("err"),
												ast.NewIdent("w"),
												ast.NewIdent("r"),
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
										Sel: ast.NewIdent("StatusNoContent"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("render"),
									Sel: ast.NewIdent("NoContent"),
								},
								Args: []ast.Expr{
									ast.NewIdent("w"),
									ast.NewIdent("r"),
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
									Name: "h",
								},
							},
							Type: &ast.StarExpr{
								X: ast.NewIdent(h.domain.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "ChiRouter",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "chi",
									},
									Sel: &ast.Ident{
										Name: "Router",
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
									Name: "router",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "chi",
										},
										Sel: &ast.Ident{
											Name: "NewRouter",
										},
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "router",
									},
									Sel: &ast.Ident{
										Name: "Route",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: "\"/\"",
									},
									&ast.FuncLit{
										Type: &ast.FuncType{
											Params: &ast.FieldList{
												List: []*ast.Field{
													{
														Names: []*ast.Ident{
															{
																Name: "g",
															},
														},
														Type: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "chi",
															},
															Sel: &ast.Ident{
																Name: "Router",
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
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Post",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Create",
																},
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Get",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "List",
																},
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Get",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Get",
																},
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Patch",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Update",
																},
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "g",
															},
															Sel: &ast.Ident{
																Name: "Delete",
															},
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "h",
																},
																Sel: &ast.Ident{
																	Name: "Delete",
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
									Name: "router",
								},
							},
						},
					},
				},
			},
		},
	}
}
