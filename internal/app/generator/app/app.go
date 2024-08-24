package app

import (
	"bytes"
	"fmt"
	"github.com/mikalai-mitsin/creathor/internal/pkg/domain"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
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
	return &ast.File{
		Name: &ast.Ident{
			Name: a.domain.DirName(),
		},
		Decls: []ast.Decl{
			a.imports(),
			a.structure(),
			a.constructor(),
			a.start(),
			a.stop(),
		},
	}
}

func (a App) imports() *ast.GenDecl {
	imports := &ast.GenDecl{
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
					Value: fmt.Sprintf(`"%s/internal/app/%s/handlers/grpc"`, a.domain.Module, a.domain.DirName()),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/app/%s/interceptors"`, a.domain.Module, a.domain.DirName()),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/app/%s/repositories/postgres"`, a.domain.Module, a.domain.DirName()),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf(`"%s/internal/app/%s/usecases"`, a.domain.Module, a.domain.DirName()),
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
			&ast.ImportSpec{
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
		},
	}
	return imports
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
				Name: "grpcServer",
			},
			Value: &ast.Ident{
				Name: "grpcServer",
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
			Key:   ast.NewIdent(fmt.Sprintf("%sRepository", a.domain.LowerCamelName())),
			Value: ast.NewIdent(a.domain.Repository.Variable),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(fmt.Sprintf("%sUseCase", a.domain.LowerCamelName())),
			Value: ast.NewIdent(a.domain.UseCase.Variable),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(fmt.Sprintf("%sInterceptor", a.domain.LowerCamelName())),
			Value: ast.NewIdent(a.domain.Interceptor.Variable),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent(fmt.Sprintf("%sHandler", a.domain.LowerCamelName())),
			Value: ast.NewIdent(a.domain.GRPCHandler.Variable),
		},
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(a.domain.Repository.Variable),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "postgres",
							},
							Sel: ast.NewIdent(fmt.Sprintf("New%s", a.domain.Repository.Name)),
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
					ast.NewIdent(a.domain.UseCase.Variable),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "usecases",
							},
							Sel: ast.NewIdent(fmt.Sprintf("New%s", a.domain.UseCase.Name)),
						},
						Args: []ast.Expr{
							ast.NewIdent(a.domain.Repository.Variable),
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
			ast.NewIdent(a.domain.Interceptor.Variable),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "interceptors",
					},
					Sel: ast.NewIdent(fmt.Sprintf("New%s", a.domain.Interceptor.Name)),
				},
				Args: []ast.Expr{
					ast.NewIdent(a.domain.UseCase.Variable),
					&ast.Ident{
						Name: "logger",
					},
				},
			},
		},
	})
	body.List = append(body.List, &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent(a.domain.GRPCHandler.Variable),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("handlers"),
					Sel: ast.NewIdent(fmt.Sprintf("New%s", a.domain.GRPCHandler.Name)),
				},
				Args: []ast.Expr{
					ast.NewIdent(a.domain.Interceptor.Variable),
					&ast.Ident{
						Name: "logger",
					},
				},
			},
		},
	})
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
						ast.NewIdent("grpcServer"),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("grpc"),
							Sel: ast.NewIdent("Server"),
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
						ast.NewIdent(fmt.Sprintf("%sRepository", a.domain.LowerCamelName())),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "postgres",
							},
							Sel: ast.NewIdent(a.domain.Repository.Name),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(fmt.Sprintf("%sUseCase", a.domain.LowerCamelName())),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "usecases",
							},
							Sel: ast.NewIdent(a.domain.UseCase.Name),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(fmt.Sprintf("%sInterceptor", a.domain.LowerCamelName())),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "interceptors",
							},
							Sel: ast.NewIdent(a.domain.Interceptor.Name),
						},
					},
				},
				{
					Names: []*ast.Ident{
						ast.NewIdent(fmt.Sprintf("%sHandler", a.domain.LowerCamelName())),
					},
					Type: &ast.StarExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("handlers"),
							Sel: ast.NewIdent(a.domain.GRPCHandler.Name),
						},
					},
				},
			},
		},
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

func (a App) start() *ast.FuncDecl {
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
		Name: ast.NewIdent("Start"),
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
							X:   ast.NewIdent(a.domain.ProtoModule),
							Sel: ast.NewIdent(fmt.Sprintf("Register%sServiceServer", a.domain.CamelName())),
						},
						Args: []ast.Expr{
							&ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent("grpcServer"),
							},
							&ast.SelectorExpr{
								X:   ast.NewIdent("a"),
								Sel: ast.NewIdent(fmt.Sprintf("%sHandler", a.domain.LowerCamelName())),
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

func (a App) stop() *ast.FuncDecl {
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
		Name: ast.NewIdent("Stop"),
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
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}
