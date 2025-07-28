package app

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/app"
)

type App struct {
	app *app.BaseEntity
}

func NewApp(domain *app.BaseEntity) *App {
	return &App{app: domain}
}

func (a App) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", a.app.AppName(), "app.go")
	err := os.MkdirAll(path.Dir(filename), 0777)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = a.file()
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

func (a App) file() *ast.File {
	decls := []ast.Decl{
		a.imports(),
		a.structure(),
		a.constructor(),
	}
	if a.app.Config.HTTPEnabled {
		decls = append(decls, a.registerHTTP())
	}
	if a.app.Config.GRPCEnabled {
		decls = append(decls, a.registerGRPC())
	}
	return &ast.File{
		Name:  ast.NewIdent(a.app.AppName()),
		Decls: decls,
	}
}

func (a App) imports() *ast.GenDecl {
	specs := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/usecases"`,
					a.app.Module,
					a.app.AppName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/repositories/postgres"`,
					a.app.Module,
					a.app.AppName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/services"`,
					a.app.Module,
					a.app.AppName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, a.app.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/grpc"`, a.app.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/http"`, a.app.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/log"`, a.app.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, a.app.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/jmoiron/sqlx"`,
			},
		},
	}
	if a.app.Config.HTTPEnabled {
		specs = append(specs, &ast.ImportSpec{
			Name: ast.NewIdent("httpHandlers"),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/handlers/http"`,
					a.app.Module,
					a.app.AppName(),
				),
			},
		})
	}
	if a.app.Config.GRPCEnabled {
		specs = append(
			specs,
			&ast.ImportSpec{
				Name: ast.NewIdent("grpcHandlers"),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s/handlers/grpc"`,
						a.app.Module,
						a.app.AppName(),
					),
				},
			}, &ast.ImportSpec{
				Name: ast.NewIdent(a.app.ProtoModule),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/pkg/%s/v1"`,
						a.app.Module,
						a.app.ProtoModule,
					),
				},
			},
		)

	}
	return &ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: specs,
	}
}

func (a App) constructor() *ast.FuncDecl {
	args := []*ast.Field{
		{
			Names: []*ast.Ident{
				ast.NewIdent("db"),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("sqlx"),
					Sel: ast.NewIdent("DB"),
				},
			},
		},
		{
			Names: []*ast.Ident{
				ast.NewIdent("logger"),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("log"),
					Sel: ast.NewIdent("Log"),
				},
			},
		},
		{
			Names: []*ast.Ident{
				ast.NewIdent("clock"),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("clock"),
					Sel: ast.NewIdent("Clock"),
				},
			},
		},
		{
			Names: []*ast.Ident{
				ast.NewIdent("uuidGenerator"),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("uuid"),
					Sel: ast.NewIdent("UUIDv7Generator"),
				},
			},
		},
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("db"),
			Value: ast.NewIdent("db"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("logger"),
			Value: ast.NewIdent("logger"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.app.GetRepositoryPrivateVariableName()),
			Value: ast.NewIdent(a.app.GetRepositoryPrivateVariableName()),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.app.GetServicePrivateVariableName()),
			Value: ast.NewIdent(a.app.GetServicePrivateVariableName()),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
			Value: ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
		},
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(a.app.GetRepositoryPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("postgres"),
							Sel: ast.NewIdent(a.app.GetRepositoryConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent("db"),
							ast.NewIdent("logger"),
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(a.app.GetServicePrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("services"),
							Sel: ast.NewIdent(a.app.GetServiceConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(a.app.GetRepositoryPrivateVariableName()),
							ast.NewIdent("clock"),
							ast.NewIdent("logger"),
							ast.NewIdent("uuidGenerator"),
						},
					},
				},
			},
		},
	}
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("usecases"),
					Sel: ast.NewIdent(a.app.GetUseCaseConstructorName()),
				},
				Args: []ast.Expr{
					ast.NewIdent(a.app.GetServicePrivateVariableName()),
					ast.NewIdent("logger"),
				},
			},
		},
	})

	if a.app.Config.HTTPEnabled {
		body.List = append(body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(a.app.GetHTTPHandlerPrivateVariableName()),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("httpHandlers"),
						Sel: ast.NewIdent(a.app.GetHTTPHandlerConstructorName()),
					},
					Args: []ast.Expr{
						ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
						ast.NewIdent("logger"),
					},
				},
			},
		})
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(a.app.GetHTTPHandlerPrivateVariableName()),
			Value: ast.NewIdent(a.app.GetHTTPHandlerPrivateVariableName()),
		})
	}
	if a.app.Config.GRPCEnabled {
		body.List = append(body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(a.app.GetGRPCHandlerPrivateVariableName()),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("grpcHandlers"),
						Sel: ast.NewIdent(a.app.GetGRPCHandlerConstructorName()),
					},
					Args: []ast.Expr{
						ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
						ast.NewIdent("logger"),
					},
				},
			},
		})
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(a.app.GetGRPCHandlerPrivateVariableName()),
			Value: ast.NewIdent(a.app.GetGRPCHandlerPrivateVariableName()),
		})
	}

	body.List = append(body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: ast.NewIdent("BaseEntity"),
					Elts: exprs,
				},
			},
		},
	})
	return &ast.FuncDecl{
		Name: ast.NewIdent("NewApp"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: args,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: ast.NewIdent("BaseEntity"),
						},
					},
				},
			},
		},
		Body: body,
	}
}

func (a App) structure() *ast.GenDecl {
	structType := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("db"),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("sqlx"),
							Sel: ast.NewIdent("DB"),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent("logger"),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("log"),
							Sel: ast.NewIdent("Log"),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.app.GetRepositoryPrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("postgres"),
							Sel: ast.NewIdent(a.app.GetRepositoryTypeName()),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.app.GetServicePrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("services"),
							Sel: ast.NewIdent(a.app.GetServiceTypeName()),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.app.GetUseCasePrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("usecases"),
							Sel: ast.NewIdent(a.app.GetUseCaseTypeName()),
						},
					},
				},
			},
		},
	}
	if a.app.Config.HTTPEnabled {
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(a.app.GetHTTPHandlerPrivateVariableName()),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("httpHandlers"),
					Sel: ast.NewIdent(a.app.GetHTTPHandlerTypeName()),
				},
			},
		})
	}
	if a.app.Config.GRPCEnabled {
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(a.app.GetGRPCHandlerPrivateVariableName()),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("grpcHandlers"),
					Sel: ast.NewIdent(a.app.GetGRPCHandlerTypeName()),
				},
			},
		})
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("BaseEntity"),
				Type: structType,
			},
		},
	}
}

func (a App) registerGRPC() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("a"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("BaseEntity"),
					},
				},
			},
		},
		Name: ast.NewIdent("RegisterGRPC"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("grpcServer"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("grpc"),
								Sel: ast.NewIdent("Server"),
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
							X:   ast.NewIdent("grpcServer"),
							Sel: ast.NewIdent("AddHandler"),
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(a.app.ProtoModule),
									Sel: ast.NewIdent(a.app.GetGRPCServiceDescriptionName()),
								},
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent(a.app.GetGRPCHandlerPrivateVariableName()),
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
	}
}

func (a App) registerHTTP() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("a"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent("BaseEntity"),
					},
				},
			},
		},
		Name: ast.NewIdent("RegisterHTTP"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("httpServer"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("http"),
								Sel: ast.NewIdent("Server"),
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
							X:   ast.NewIdent("httpServer"),
							Sel: ast.NewIdent("Mount"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"/api/v1/%s/"`, a.app.GetHTTPPath()),
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: ast.NewIdent("a"),
										Sel: ast.NewIdent(
											a.app.GetHTTPHandlerPrivateVariableName(),
										),
									},
									Sel: ast.NewIdent("ChiRouter"),
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
	}
}
