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

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

const destinationPath = "."

type HandlerGenerator struct {
	entityConfig configs.EntityConfig
}

func NewHandlerGenerator(entityConfig configs.EntityConfig) *HandlerGenerator {
	return &HandlerGenerator{
		entityConfig: entityConfig,
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
	return path.Join(
		"internal",
		"app",
		h.entityConfig.AppConfig.AppName(),
		"handlers",
		"http",
		h.entityConfig.DirName(),
		h.entityConfig.FileName(),
	)
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
							Value: h.entityConfig.AppConfig.ProjectConfig.ErrsImportPath(),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: h.entityConfig.AppConfig.ProjectConfig.UUIDImportPath(),
						},
					},
					&ast.ImportSpec{
						Name: ast.NewIdent("httpServer"),
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: h.entityConfig.AppConfig.ProjectConfig.HTTPImportPath(),
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
						Name: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{
											ast.NewIdent(h.entityConfig.GetUseCasePrivateVariableName()),
										},
										Type: ast.NewIdent(h.entityConfig.GetUseCaseInterfaceName()),
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
				Name: ast.NewIdent(h.entityConfig.GetHTTPHandlerConstructorName()),
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									ast.NewIdent(h.entityConfig.GetUseCasePrivateVariableName()),
								},
								Type: ast.NewIdent(h.entityConfig.GetUseCaseInterfaceName()),
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
									X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
										Type: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: ast.NewIdent(
													h.entityConfig.GetUseCasePrivateVariableName(),
												),
												Value: ast.NewIdent(
													h.entityConfig.GetUseCasePrivateVariableName(),
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
							Text: fmt.Sprintf(
								"// @Summary Create %s",
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.entityConfig.GetOneVariableName()),
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
							Text: fmt.Sprintf(
								"// @Param form body %s true \"Create %s request\"",
								h.entityConfig.GetHTTPCreateDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf(
								"// @Success 201 {object} %s \"Created %s\"",
								h.entityConfig.GetHTTPItemDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
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
							Text: fmt.Sprintf(
								"// @Router /api/v1/%s/%s/ [POST]",
								h.entityConfig.AppConfig.AppName(),
								h.entityConfig.GetHTTPPath(),
							),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPCreateDTOConstructorName()),
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
								ast.NewIdent(h.entityConfig.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.entityConfig.GetUseCasePrivateVariableName(),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
							Text: fmt.Sprintf(
								"// @Summary Get %s by id",
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.entityConfig.GetOneVariableName()),
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
							Text: fmt.Sprintf(
								"// @Success 200 {object} %s \"Requested %s\"",
								h.entityConfig.GetHTTPItemDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
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
							Text: fmt.Sprintf(
								"// @Router /api/v1/%s/%s/{id} [GET]",
								h.entityConfig.AppConfig.AppName(),
								h.entityConfig.GetHTTPPath(),
							),
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
										Sel: ast.NewIdent("MustParse"),
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
								ast.NewIdent(h.entityConfig.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.entityConfig.GetUseCasePrivateVariableName(),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
							Text: fmt.Sprintf(
								"// @Summary List of %s",
								h.entityConfig.GetManyVariableName(),
							),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.entityConfig.GetOneVariableName()),
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
							Text: fmt.Sprintf(
								"// @Param filter query %s true \"Filter of %s\"",
								h.entityConfig.GetHTTPFilterDTOName(),
								h.entityConfig.GetManyVariableName(),
							),
						},
						{
							Text: fmt.Sprintf(
								"// @Success 200 {object} %s \"Filtered list of %s\"",
								h.entityConfig.GetHTTPListDTOName(),
								h.entityConfig.GetManyVariableName(),
							),
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
							Text: fmt.Sprintf(
								"// @Router /api/v1/%s/%s/ [GET]",
								h.entityConfig.AppConfig.AppName(),
								h.entityConfig.GetHTTPPath(),
							),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPFilterDTOConstructorName()),
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
								ast.NewIdent(h.entityConfig.GetManyVariableName()),
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
												h.entityConfig.GetUseCasePrivateVariableName(),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPListDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.entityConfig.GetManyVariableName()),
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
							Text: fmt.Sprintf(
								"// @Summary Update %s",
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.entityConfig.GetOneVariableName()),
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
							Text: fmt.Sprintf(
								"// @Param form body %s true \"Update %s request\"",
								h.entityConfig.GetHTTPUpdateDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf(
								"// @Success 200 {object} %s \"Updated %s\"",
								h.entityConfig.GetHTTPItemDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
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
							Text: fmt.Sprintf(
								"// @Router /api/v1/%s/%s/{id} [PATCH]",
								h.entityConfig.AppConfig.AppName(),
								h.entityConfig.GetHTTPPath(),
							),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPUpdateDTOConstructorName()),
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
								ast.NewIdent(h.entityConfig.GetOneVariableName()),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X: ast.NewIdent("h"),
											Sel: ast.NewIdent(
												h.entityConfig.GetUseCasePrivateVariableName(),
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
									Fun: ast.NewIdent(h.entityConfig.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
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
							Text: fmt.Sprintf(
								"// @Summary Delete %s by id",
								h.entityConfig.GetOneVariableName(),
							),
						},
						{
							Text: fmt.Sprintf("// @Tags %s", h.entityConfig.GetOneVariableName()),
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
							Text: fmt.Sprintf(
								"// @Success 200 {object} %s \"Updated %s\"",
								h.entityConfig.GetHTTPItemDTOName(),
								h.entityConfig.GetOneVariableName(),
							),
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
							Text: fmt.Sprintf(
								"// @Router /api/v1/%s/%s/{id} [DELETE]",
								h.entityConfig.AppConfig.AppName(),
								h.entityConfig.GetHTTPPath(),
							),
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
								ast.NewIdent("delDTO"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.entityConfig.GetHTTPDeleteDTOConstructorName()),
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
								ast.NewIdent("del"),
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("delDTO"),
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
								&ast.Ident{
									Name: h.entityConfig.GetOneVariableName(),
								},
								ast.NewIdent("err"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent("h"),
											Sel: ast.NewIdent(h.entityConfig.GetUseCasePrivateVariableName()),
										},
										Sel: ast.NewIdent("Delete"),
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "r",
												},
												Sel: &ast.Ident{
													Name: "Context",
												},
											},
										},
										&ast.Ident{
											Name: "del",
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
												},
											},
										},
									},
									&ast.ReturnStmt{},
								},
							},
						},
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "response",
								},
								&ast.Ident{
									Name: "err",
								},
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: ast.NewIdent(h.entityConfig.GetHTTPItemDTOConstructorName()),
									Args: []ast.Expr{
										ast.NewIdent(h.entityConfig.GetOneVariableName()),
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
									&ast.ExprStmt{
										X: &ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "errs",
												},
												Sel: &ast.Ident{
													Name: "RenderToHTTPResponse",
												},
											},
											Args: []ast.Expr{
												&ast.Ident{
													Name: "err",
												},
												&ast.Ident{
													Name: "w",
												},
												&ast.Ident{
													Name: "r",
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
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "Status",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "r",
									},
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
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "render",
									},
									Sel: &ast.Ident{
										Name: "JSON",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "w",
									},
									&ast.Ident{
										Name: "r",
									},
									&ast.Ident{
										Name: "response",
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
								X: ast.NewIdent(h.entityConfig.GetHTTPHandlerTypeName()),
							},
						},
					},
				},
				Name: ast.NewIdent("router"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("chi"),
									Sel: ast.NewIdent("Router"),
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								ast.NewIdent("router"),
							},
							Tok: token.DEFINE,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("chi"),
										Sel: ast.NewIdent("NewRouter"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   ast.NewIdent("router"),
									Sel: ast.NewIdent("Route"),
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
															ast.NewIdent("g"),
														},
														Type: &ast.SelectorExpr{
															X:   ast.NewIdent("chi"),
															Sel: ast.NewIdent("Router"),
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
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Post"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("Create"),
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Get"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("List"),
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Get"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("Get"),
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Patch"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("Update"),
															},
														},
													},
												},
												&ast.ExprStmt{
													X: &ast.CallExpr{
														Fun: &ast.SelectorExpr{
															X:   ast.NewIdent("g"),
															Sel: ast.NewIdent("Delete"),
														},
														Args: []ast.Expr{
															&ast.BasicLit{
																Kind:  token.STRING,
																Value: "\"/{id}\"",
															},
															&ast.SelectorExpr{
																X:   ast.NewIdent("h"),
																Sel: ast.NewIdent("Delete"),
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
								ast.NewIdent("router"),
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
								X: &ast.Ident{
									Name: h.entityConfig.GetHTTPHandlerTypeName(),
								},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "RegisterHTTP",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{
										Name: "httpServer",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "httpServer",
										},
										Sel: &ast.Ident{
											Name: "Server",
										},
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
										Name: "httpServer",
									},
									Sel: &ast.Ident{
										Name: "Mount",
									},
								},
								Args: []ast.Expr{
									&ast.BasicLit{
										Kind: token.STRING,
										Value: fmt.Sprintf(
											`"/api/v1/%s/%s"`,
											h.entityConfig.AppConfig.AppName(),
											h.entityConfig.GetHTTPPath(),
										),
									},
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "h",
											},
											Sel: &ast.Ident{
												Name: "router",
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
