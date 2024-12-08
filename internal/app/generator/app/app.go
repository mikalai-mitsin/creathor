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

	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
)

type App struct {
	domain *domain.Domain
}

func NewApp(domain *domain.Domain) *App {
	return &App{domain: domain}
}

func (a App) Sync() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "app", a.domain.DirName(), "app.go")
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
	if a.domain.Config.HTTPEnabled {
		decls = append(decls, a.registerHTTP())
	}
	if a.domain.Config.GRPCEnabled {
		decls = append(decls, a.registerGRPC())
	}
	return &ast.File{
		Name: &ast.Ident{
			Name: a.domain.DirName(),
		},
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
					a.domain.Module,
					a.domain.DirName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/repositories/postgres"`,
					a.domain.Module,
					a.domain.DirName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/services"`,
					a.domain.Module,
					a.domain.DirName(),
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/clock"`, a.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/grpc"`, a.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/http"`, a.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/log"`, a.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s/internal/pkg/uuid"`, a.domain.Module),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/jmoiron/sqlx"`,
			},
		},
	}
	if a.domain.Config.HTTPEnabled {
		specs = append(specs, &ast.ImportSpec{
			Name: ast.NewIdent("httpHandlers"),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/internal/app/%s/handlers/http"`,
					a.domain.Module,
					a.domain.DirName(),
				),
			},
		})
	}
	if a.domain.Config.GRPCEnabled {
		specs = append(
			specs,
			&ast.ImportSpec{
				Name: ast.NewIdent("grpcHandlers"),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s/handlers/grpc"`,
						a.domain.Module,
						a.domain.DirName(),
					),
				},
			}, &ast.ImportSpec{
				Name: ast.NewIdent(a.domain.ProtoModule),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/pkg/%s/v1"`,
						a.domain.Module,
						a.domain.ProtoModule,
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
				{
					Name: "db",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "sqlx",
					},
					Sel: &ast.Ident{
						Name: "DB",
					},
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "logger",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "log",
					},
					Sel: &ast.Ident{
						Name: "Log",
					},
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "clock",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "clock",
					},
					Sel: &ast.Ident{
						Name: "Clock",
					},
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "uuidGenerator",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "uuid",
					},
					Sel: &ast.Ident{
						Name: "UUIDv4Generator",
					},
				},
			},
		},
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key: &ast.Ident{
				Name: "db",
			},
			Value: &ast.Ident{
				Name: "db",
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
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.domain.GetRepositoryPrivateVariableName()),
			Value: ast.NewIdent(a.domain.GetRepositoryPrivateVariableName()),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.domain.GetServicePrivateVariableName()),
			Value: ast.NewIdent(a.domain.GetServicePrivateVariableName()),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
			Value: ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
		},
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(a.domain.GetRepositoryPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "postgres",
							},
							Sel: ast.NewIdent(a.domain.GetRepositoryConstructorName()),
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "db",
							},
							&ast.Ident{
								Name: "logger",
							},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(a.domain.GetServicePrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "services",
							},
							Sel: ast.NewIdent(a.domain.GetServiceConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(a.domain.GetRepositoryPrivateVariableName()),
							&ast.Ident{
								Name: "clock",
							},
							&ast.Ident{
								Name: "logger",
							},
							&ast.Ident{
								Name: "uuidGenerator",
							},
						},
					},
				},
			},
		},
	}
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "usecases",
					},
					Sel: ast.NewIdent(a.domain.GetUseCaseConstructorName()),
				},
				Args: []ast.Expr{
					ast.NewIdent(a.domain.GetServicePrivateVariableName()),
					&ast.Ident{
						Name: "logger",
					},
				},
			},
		},
	})

	if a.domain.Config.HTTPEnabled {
		body.List = append(body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(a.domain.GetHTTPHandlerPrivateVariableName()),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("httpHandlers"),
						Sel: ast.NewIdent(a.domain.GetHTTPHandlerConstructorName()),
					},
					Args: []ast.Expr{
						ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
						&ast.Ident{
							Name: "logger",
						},
					},
				},
			},
		})
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(a.domain.GetHTTPHandlerPrivateVariableName()),
			Value: ast.NewIdent(a.domain.GetHTTPHandlerPrivateVariableName()),
		})
	}
	if a.domain.Config.GRPCEnabled {
		body.List = append(body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(a.domain.GetGRPCHandlerPrivateVariableName()),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("grpcHandlers"),
						Sel: ast.NewIdent(a.domain.GetGRPCHandlerConstructorName()),
					},
					Args: []ast.Expr{
						ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
						&ast.Ident{
							Name: "logger",
						},
					},
				},
			},
		})
		exprs = append(exprs, &ast.KeyValueExpr{
			Key:   ast.NewIdent(a.domain.GetGRPCHandlerPrivateVariableName()),
			Value: ast.NewIdent(a.domain.GetGRPCHandlerPrivateVariableName()),
		})
	}

	body.List = append(body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: &ast.Ident{
						Name: "App",
					},
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
							X: ast.NewIdent("App"),
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
						{
							Name: "db",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "sqlx",
							},
							Sel: &ast.Ident{
								Name: "DB",
							},
						},
					},
				},
				{
					Names: []*ast.Ident{
						{
							Name: "logger",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "log",
							},
							Sel: &ast.Ident{
								Name: "Log",
							},
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.domain.GetRepositoryPrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "postgres",
							},
							Sel: ast.NewIdent(a.domain.GetRepositoryTypeName()),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.domain.GetServicePrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "services",
							},
							Sel: ast.NewIdent(a.domain.GetServiceTypeName()),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(a.domain.GetUseCasePrivateVariableName()),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "usecases",
							},
							Sel: ast.NewIdent(a.domain.GetUseCaseTypeName()),
						},
					},
				},
			},
		},
	}
	if a.domain.Config.HTTPEnabled {
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(a.domain.GetHTTPHandlerPrivateVariableName()),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("httpHandlers"),
					Sel: ast.NewIdent(a.domain.GetHTTPHandlerTypeName()),
				},
			},
		})
	}
	if a.domain.Config.GRPCEnabled {
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(a.domain.GetGRPCHandlerPrivateVariableName()),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("grpcHandlers"),
					Sel: ast.NewIdent(a.domain.GetGRPCHandlerTypeName()),
				},
			},
		})
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{
					Name: "App",
				},
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
						X: ast.NewIdent("App"),
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
							{
								Name: "grpcServer",
							},
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
							X: ast.NewIdent("grpcServer"),
							Sel: &ast.Ident{
								Name: "AddHandler",
							},
						},
						Args: []ast.Expr{
							&ast.UnaryExpr{
								Op: token.AND,
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(a.domain.ProtoModule),
									Sel: ast.NewIdent(a.domain.GetGRPCServiceDescriptionName()),
								},
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent(a.domain.GetGRPCHandlerPrivateVariableName()),
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
						{
							Name: "a",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "App",
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
									Name: "http",
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
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"/%s/"`, a.domain.GetManyVariableName()),
							},
							&ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: ast.NewIdent(
											a.domain.GetHTTPHandlerPrivateVariableName(),
										),
									},
									Sel: &ast.Ident{
										Name: "ChiRouter",
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
	}
}
