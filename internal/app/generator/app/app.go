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

	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type App struct {
	app *configs.AppConfig
}

func NewApp(domain *configs.AppConfig) *App {
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
	if a.app.HTTPEnabled {
		decls = append(decls, a.registerHTTP())
	}
	if a.app.GRPCEnabled {
		decls = append(decls, a.registerGRPC())
	}
	if a.app.KafkaEnabled {
		decls = append(decls, a.registerKafka())
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
				Kind:  token.STRING,
				Value: a.app.ProjectConfig.ClockImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: a.app.ProjectConfig.LogImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: a.app.ProjectConfig.UUIDImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: a.app.ProjectConfig.KafkaImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"github.com/jmoiron/sqlx"`,
			},
		},
	}
	for _, entity := range a.app.Entities {
		specs = append(
			specs,
			&ast.ImportSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sUseCases", entity.LowerCamelName())),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s/usecases/%s"`,
						a.app.Module,
						a.app.AppName(),
						entity.DirName(),
					),
				},
			},
			&ast.ImportSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sRepositories", entity.LowerCamelName())),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s/repositories/postgres/%s"`,
						a.app.Module,
						a.app.AppName(),
						entity.DirName(),
					),
				},
			},
			&ast.ImportSpec{
				Name: ast.NewIdent(fmt.Sprintf("%sServices", entity.LowerCamelName())),
				Path: &ast.BasicLit{
					Kind: token.STRING,
					Value: fmt.Sprintf(
						`"%s/internal/app/%s/services/%s"`,
						a.app.Module,
						a.app.AppName(),
						entity.DirName(),
					),
				},
			},
		)
		if a.app.KafkaEnabled {
			specs = append(
				specs,
				&ast.ImportSpec{
					Name: ast.NewIdent(fmt.Sprintf("%sEvents", entity.LowerCamelName())),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(
							`"%s/internal/app/%s/repositories/kafka/%s"`,
							a.app.Module,
							a.app.AppName(),
							entity.DirName(),
						),
					},
				},
				&ast.ImportSpec{
					Name: ast.NewIdent(fmt.Sprintf("%sKafkaHandlers", entity.LowerCamelName())),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(
							`"%s/internal/app/%s/handlers/kafka/%s"`,
							a.app.Module,
							a.app.AppName(),
							entity.DirName(),
						),
					},
				},
			)
		}
		if a.app.HTTPEnabled {
			specs = append(
				specs,
				&ast.ImportSpec{
					Name: ast.NewIdent(fmt.Sprintf("%sHttpHandlers", entity.LowerCamelName())),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(
							`"%s/internal/app/%s/handlers/http/%s"`,
							a.app.Module,
							a.app.AppName(),
							entity.DirName(),
						),
					},
				},
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: a.app.ProjectConfig.HTTPImportPath(),
					},
				},
			)
		}
		if a.app.GRPCEnabled {
			specs = append(
				specs,
				&ast.ImportSpec{
					Name: ast.NewIdent(fmt.Sprintf("%sGrpcHandlers", entity.LowerCamelName())),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(
							`"%s/internal/app/%s/handlers/grpc/%s"`,
							a.app.Module,
							a.app.AppName(),
							entity.DirName(),
						),
					},
				},
				&ast.ImportSpec{
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: a.app.ProjectConfig.GRPCImportPath(),
					},
				},
				&ast.ImportSpec{
					Name: ast.NewIdent(a.app.ProtoPackage),
					Path: &ast.BasicLit{
						Kind: token.STRING,
						Value: fmt.Sprintf(
							`"%s/pkg/%s/v1"`,
							a.app.Module,
							a.app.ProtoPackage,
						),
					},
				},
			)
		}
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
				ast.NewIdent("readDB"),
				ast.NewIdent("writeDB"),
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
			Type: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Logger"),
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
	if a.app.KafkaEnabled {
		args = append(args,
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("kafkaProducer"),
				},
				Type: &ast.StarExpr{X: &ast.SelectorExpr{
					X:   ast.NewIdent("kafka"),
					Sel: ast.NewIdent("Producer"),
				}},
			},
		)
	}
	exprs := []ast.Expr{
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("readDB"),
			Value: ast.NewIdent("readDB"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("writeDB"),
			Value: ast.NewIdent("writeDB"),
		},
		&ast.KeyValueExpr{
			Key:   ast.NewIdent("logger"),
			Value: ast.NewIdent("logger"),
		},
	}
	if a.app.KafkaEnabled {
		exprs = append(exprs,
			&ast.KeyValueExpr{
				Key:   ast.NewIdent("kafkaProducer"),
				Value: ast.NewIdent("kafkaProducer"),
			},
		)
	}
	body := &ast.BlockStmt{
		List: []ast.Stmt{},
	}
	for _, entity := range a.app.Entities {
		exprs = append(exprs,
			&ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetRepositoryPrivateVariableName()),
				Value: ast.NewIdent(entity.GetRepositoryPrivateVariableName()),
			},
			&ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetServicePrivateVariableName()),
				Value: ast.NewIdent(entity.GetServicePrivateVariableName()),
			},
			&ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetUseCasePrivateVariableName()),
				Value: ast.NewIdent(entity.GetUseCasePrivateVariableName()),
			},
		)
		body.List = append(body.List,
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetRepositoryPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sRepositories", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.GetRepositoryConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent("readDB"),
							ast.NewIdent("writeDB"),
							ast.NewIdent("logger"),
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetServicePrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sServices", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.GetServiceConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(entity.GetRepositoryPrivateVariableName()),
							ast.NewIdent("clock"),
							ast.NewIdent("logger"),
							ast.NewIdent("uuidGenerator"),
						},
					},
				},
			})
		if a.app.KafkaEnabled {
			body.List = append(body.List, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetEventProducerPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sEvents", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.EventProducerConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent("kafkaProducer"),
							ast.NewIdent("logger"),
						},
					},
				},
			})
		}
		useCaseArgs := []ast.Expr{
			ast.NewIdent(entity.GetServicePrivateVariableName()),
		}
		if a.app.KafkaEnabled {
			useCaseArgs = append(useCaseArgs, ast.NewIdent(entity.GetEventProducerPrivateVariableName()))
		}
		useCaseArgs = append(useCaseArgs, ast.NewIdent("logger"))
		body.List = append(body.List, &ast.AssignStmt{
			Lhs: []ast.Expr{
				ast.NewIdent(entity.GetUseCasePrivateVariableName()),
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sUseCases", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.GetUseCaseConstructorName()),
					},
					Args: useCaseArgs,
				},
			},
		})

		if a.app.HTTPEnabled {
			body.List = append(body.List, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetHTTPHandlerPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sHttpHandlers", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.GetHTTPHandlerConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(entity.GetUseCasePrivateVariableName()),
							ast.NewIdent("logger"),
						},
					},
				},
			})
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetHTTPHandlerPrivateVariableName()),
				Value: ast.NewIdent(entity.GetHTTPHandlerPrivateVariableName()),
			})
		}
		if a.app.KafkaEnabled {
			body.List = append(body.List, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetKafkaHandlerPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sKafkaHandlers", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.KafkaHandlerConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(entity.GetUseCasePrivateVariableName()),
							ast.NewIdent("logger"),
						},
					},
				},
			})
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetEventProducerPrivateVariableName()),
				Value: ast.NewIdent(entity.GetEventProducerPrivateVariableName()),
			}, &ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetKafkaHandlerPrivateVariableName()),
				Value: ast.NewIdent(entity.GetKafkaHandlerPrivateVariableName()),
			})
		}
		if a.app.GRPCEnabled {
			body.List = append(body.List, &ast.AssignStmt{
				Lhs: []ast.Expr{
					ast.NewIdent(entity.GetGRPCHandlerPrivateVariableName()),
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent(fmt.Sprintf("%sGrpcHandlers", entity.LowerCamelName())),
							Sel: ast.NewIdent(entity.GetGRPCHandlerConstructorName()),
						},
						Args: []ast.Expr{
							ast.NewIdent(entity.GetUseCasePrivateVariableName()),
							ast.NewIdent("logger"),
						},
					},
				},
			})
			exprs = append(exprs, &ast.KeyValueExpr{
				Key:   ast.NewIdent(entity.GetGRPCHandlerPrivateVariableName()),
				Value: ast.NewIdent(entity.GetGRPCHandlerPrivateVariableName()),
			})
		}

	}
	body.List = append(body.List, &ast.ReturnStmt{
		Results: []ast.Expr{
			&ast.UnaryExpr{
				Op: token.AND,
				X: &ast.CompositeLit{
					Type: ast.NewIdent("App"),
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
						ast.NewIdent("readDB"),
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
						ast.NewIdent("writeDB"),
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
					Type: &ast.SelectorExpr{
						X:   ast.NewIdent("log"),
						Sel: ast.NewIdent("Logger"),
					},
				},
			},
		},
	}
	if a.app.KafkaEnabled {
		structType.Fields.List = append(structType.Fields.List,
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent("kafkaProducer"),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("kafka"),
						Sel: ast.NewIdent("Producer"),
					},
				},
			},
		)
	}
	for _, entity := range a.app.Entities {
		structType.Fields.List = append(structType.Fields.List, &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent(entity.GetRepositoryPrivateVariableName()),
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent(fmt.Sprintf("%sRepositories", entity.LowerCamelName())),
					Sel: ast.NewIdent(entity.GetRepositoryTypeName()),
				},
			},
		},
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetServicePrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sServices", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.GetServiceTypeName()),
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetUseCasePrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sUseCases", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.GetUseCaseTypeName()),
					},
				},
			})
		if a.app.HTTPEnabled {
			structType.Fields.List = append(structType.Fields.List, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetHTTPHandlerPrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sHttpHandlers", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.GetHTTPHandlerTypeName()),
					},
				},
			})
		}
		if a.app.KafkaEnabled {
			structType.Fields.List = append(structType.Fields.List, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetEventProducerPrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sEvents", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.EventProducerTypeName()),
					},
				},
			}, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetKafkaHandlerPrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sKafkaHandlers", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.KafkaHandlerTypeName()),
					},
				},
			})
		}
		if a.app.GRPCEnabled {
			structType.Fields.List = append(structType.Fields.List, &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(entity.GetGRPCHandlerPrivateVariableName()),
				},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent(fmt.Sprintf("%sGrpcHandlers", entity.LowerCamelName())),
						Sel: ast.NewIdent(entity.GetGRPCHandlerTypeName()),
					},
				},
			})
		}
	}
	return &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent("App"),
				Type: structType,
			},
		},
	}
}

func (a App) registerGRPC() *ast.FuncDecl {
	stmts := make([]ast.Stmt, 0, len(a.app.Entities)+1)
	for _, entity := range a.app.Entities {
		stmts = append(stmts, &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent("grpcServer"),
					Sel: ast.NewIdent("AddHandler"),
				},
				Args: []ast.Expr{
					&ast.UnaryExpr{
						Op: token.AND,
						X: &ast.SelectorExpr{
							X:   ast.NewIdent(entity.ProtoPackage),
							Sel: ast.NewIdent(entity.GetGRPCServiceDescriptionName()),
						},
					},
					&ast.SelectorExpr{
						X:   ast.NewIdent("a"),
						Sel: ast.NewIdent(entity.GetGRPCHandlerPrivateVariableName()),
					},
				},
			},
		})
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})
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
			List: stmts,
		},
	}
}

func (a App) registerHTTP() *ast.FuncDecl {
	stmts := make([]ast.Stmt, 0, len(a.app.Entities)+1)
	for _, entity := range a.app.Entities {
		stmts = append(stmts,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("httpServer"),
						Sel: ast.NewIdent("Mount"),
					},
					Args: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"/api/v1/%s/%s"`, a.app.AppName(), entity.GetHTTPPath()),
						},
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: ast.NewIdent("a"),
									Sel: ast.NewIdent(
										entity.GetHTTPHandlerPrivateVariableName(),
									),
								},
								Sel: ast.NewIdent("ChiRouter"),
							},
						},
					},
				},
			},
		)
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})
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
			List: stmts,
		},
	}
}

func (a App) registerKafka() *ast.FuncDecl {
	stmts := make([]ast.Stmt, 0, len(a.app.Entities)+1)
	for _, entity := range a.app.Entities {
		stmts = append(stmts,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("consumer"),
						Sel: ast.NewIdent("AddHandler"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "kafka",
								},
								Sel: &ast.Ident{
									Name: "NewHandler",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.CreatedTopicName()),
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.KafkaCreatedConsumerGroup()),
								},
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: &ast.Ident{
											Name: entity.GetKafkaHandlerPrivateVariableName(),
										},
									},
									Sel: &ast.Ident{
										Name: "Created",
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
						X:   ast.NewIdent("consumer"),
						Sel: ast.NewIdent("AddHandler"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "kafka",
								},
								Sel: &ast.Ident{
									Name: "NewHandler",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.UpdatedTopicName()),
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.KafkaUpdatedConsumerGroup()),
								},
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: &ast.Ident{
											Name: entity.GetKafkaHandlerPrivateVariableName(),
										},
									},
									Sel: &ast.Ident{
										Name: "Updated",
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
						X:   ast.NewIdent("consumer"),
						Sel: ast.NewIdent("AddHandler"),
					},
					Args: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "kafka",
								},
								Sel: &ast.Ident{
									Name: "NewHandler",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.DeletedTopicName()),
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf(`"%s"`, entity.KafkaDeletedConsumerGroup()),
								},
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "a",
										},
										Sel: &ast.Ident{
											Name: entity.GetKafkaHandlerPrivateVariableName(),
										},
									},
									Sel: &ast.Ident{
										Name: "Deleted",
									},
								},
							},
						},
					},
				},
			},
		)
	}
	stmts = append(stmts, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})
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
		Name: ast.NewIdent("RegisterKafka"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							ast.NewIdent("consumer"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("kafka"),
								Sel: ast.NewIdent("Consumer"),
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
			List: stmts,
		},
	}
}
