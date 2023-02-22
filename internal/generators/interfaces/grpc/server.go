package grpc

import (
	"bytes"
	"fmt"
	"github.com/018bf/creathor/internal/configs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
)

type Server struct {
	Config *configs.Project
}

func (s Server) Sync() error {
	if err := s.SyncServerStruct(); err != nil {
		return err
	}
	if err := s.SyncServerConstructor(); err != nil {
		return err
	}
	if err := s.SyncServerStart(); err != nil {
		return err
	}
	if err := s.SyncServerStop(); err != nil {
		return err
	}
	if err := s.SyncMessageProducer(); err != nil {
		return err
	}
	if err := s.SyncDecodeError(); err != nil {
		return err
	}
	return nil
}

func (s Server) AstServerStruct() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: &ast.Ident{
			Name: "Server",
		},
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "server",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "grpc",
								},
								Sel: &ast.Ident{
									Name: "Server",
								},
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "config",
							},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "configs",
								},
								Sel: &ast.Ident{
									Name: "Config",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s Server) SyncServerStruct() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var structureExists bool
	var structure *ast.TypeSpec
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.TypeSpec); ok && t.Name.String() == "Server" {
			structure = t
			structureExists = true
			return false
		}
		return true
	})
	if structure == nil {
		structure = s.AstServerStruct()
	}
	//for _, param := range m.Params {
	//	ast.Inspect(structure, func(node ast.Node) bool {
	//		if st, ok := node.(*ast.StructType); ok && st.Fields != nil {
	//			for _, field := range st.Fields.List {
	//				for _, fieldName := range field.Names {
	//					if fieldName.Name == param.GetPublicName() {
	//						return false
	//					}
	//				}
	//			}
	//			st.Fields.List = append(st.Fields.List, &ast.Field{
	//				Doc:   nil,
	//				Names: []*ast.Ident{ast.NewIdent(param.GetPublicName())},
	//				Type:  ast.NewIdent(param.Type),
	//				Tag: &ast.BasicLit{
	//					Kind:  token.STRING,
	//					Value: fmt.Sprintf("`json:\"%s\"`", param.GetTag()),
	//				},
	//				Comment: nil,
	//			})
	//			return true
	//		}
	//		return true
	//	})
	//}
	if !structureExists {
		gd := &ast.GenDecl{
			Doc:    nil,
			TokPos: 0,
			Tok:    token.TYPE,
			Lparen: 0,
			Specs:  []ast.Spec{structure},
			Rparen: 0,
		}
		file.Decls = append(file.Decls, gd)
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

func (s Server) AstServerConstructor() *ast.FuncDecl {
	fields := []*ast.Field{
		{
			Names: []*ast.Ident{
				{
					Name: "logger",
				},
			},
			Type: &ast.SelectorExpr{
				X: &ast.Ident{
					Name: "log",
				},
				Sel: &ast.Ident{
					Name: "Logger",
				},
			},
		},
		{
			Names: []*ast.Ident{
				{
					Name: "config",
				},
			},
			Type: &ast.StarExpr{
				X: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "configs",
					},
					Sel: &ast.Ident{
						Name: "Config",
					},
				},
			},
		},
	}
	registerStmts := []ast.Stmt{}
	if s.Config.Auth {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{
					{
						Name: "authMiddleware",
					},
				},
				Type: &ast.StarExpr{
					X: &ast.Ident{
						Name: "AuthMiddleware",
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					{
						Name: "authHandler",
					},
				},
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: s.Config.ProtoPackage(),
					},
					Sel: &ast.Ident{
						Name: "AuthServiceServer",
					},
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					{
						Name: "userHandler",
					},
				},
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: s.Config.ProtoPackage(),
					},
					Sel: &ast.Ident{
						Name: "UserServiceServer",
					},
				},
			},
		)
		registerStmts = append(
			registerStmts,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "examplepb",
						},
						Sel: &ast.Ident{
							Name: "RegisterAuthServiceServer",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "server",
						},
						&ast.Ident{
							Name: "authHandler",
						},
					},
				},
			},
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "examplepb",
						},
						Sel: &ast.Ident{
							Name: "RegisterUserServiceServer",
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "server",
						},
						&ast.Ident{
							Name: "userHandler",
						},
					},
				},
			},
		)
	}
	for _, modelConfig := range s.Config.Models {
		fields = append(
			fields,
			&ast.Field{
				Names: []*ast.Ident{
					{
						Name: modelConfig.GRPCHandlerVariableName(),
					},
				},
				Type: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: s.Config.ProtoPackage(),
					},
					Sel: &ast.Ident{
						Name: fmt.Sprintf("%sServiceServer", modelConfig.ModelName()),
					},
				},
			},
		)
		registerStmts = append(
			registerStmts,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: s.Config.ProtoPackage(),
						},
						Sel: &ast.Ident{
							Name: fmt.Sprintf("Register%sServiceServer", modelConfig.ModelName()),
						},
					},
					Args: []ast.Expr{
						&ast.Ident{
							Name: "server",
						},
						&ast.Ident{
							Name: modelConfig.GRPCHandlerVariableName(),
						},
					},
				},
			},
		)
	}
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "NewServer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: fields,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{
								Name: "Server",
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
							Name: "server",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "grpc",
								},
								Sel: &ast.Ident{
									Name: "NewServer",
								},
							},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "ChainStreamInterceptor",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "otelgrpc",
												},
												Sel: &ast.Ident{
													Name: "StreamServerInterceptor",
												},
											},
										},
									},
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "grpc",
										},
										Sel: &ast.Ident{
											Name: "ChainUnaryInterceptor",
										},
									},
									Args: []ast.Expr{
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "otelgrpc",
												},
												Sel: &ast.Ident{
													Name: "UnaryServerInterceptor",
												},
											},
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "authMiddleware",
											},
											Sel: &ast.Ident{
												Name: "UnaryServerInterceptorAuth",
											},
										},
										&ast.Ident{
											Name: "RequestIDUnaryServerInterceptor",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "grpcZap",
												},
												Sel: &ast.Ident{
													Name: "UnaryServerInterceptor",
												},
											},
											Args: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "logger",
														},
														Sel: &ast.Ident{
															Name: "Logger",
														},
													},
												},
												&ast.CallExpr{
													Fun: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "grpcZap",
														},
														Sel: &ast.Ident{
															Name: "WithMessageProducer",
														},
													},
													Args: []ast.Expr{
														&ast.Ident{
															Name: "DefaultMessageProducer",
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
								Name: "reflection",
							},
							Sel: &ast.Ident{
								Name: "Register",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "server",
							},
						},
					},
				},
				&ast.BlockStmt{
					List: registerStmts,
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "healthServer",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "health",
								},
								Sel: &ast.Ident{
									Name: "NewServer",
								},
							},
						},
					},
				},
				&ast.RangeStmt{
					Key: &ast.Ident{
						Name: "service",
					},
					Tok: token.DEFINE,
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "server",
							},
							Sel: &ast.Ident{
								Name: "GetServiceInfo",
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "healthServer",
										},
										Sel: &ast.Ident{
											Name: "SetServingStatus",
										},
									},
									Args: []ast.Expr{
										&ast.Ident{
											Name: "service",
										},
										&ast.SelectorExpr{
											X: &ast.Ident{
												Name: "grpc_health_v1",
											},
											Sel: &ast.Ident{
												Name: "HealthCheckResponse_SERVING",
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
								Name: "grpc_health_v1",
							},
							Sel: &ast.Ident{
								Name: "RegisterHealthServer",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "server",
							},
							&ast.Ident{
								Name: "healthServer",
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.Ident{
									Name: "Server",
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: "server",
										},
										Value: &ast.Ident{
											Name: "server",
										},
									},
									&ast.KeyValueExpr{
										Key: &ast.Ident{
											Name: "config",
										},
										Value: &ast.Ident{
											Name: "config",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s Server) SyncServerConstructor() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "NewServer" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = s.AstServerConstructor()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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

func (s Server) AstServerStart() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "Server",
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Start",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "_",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
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
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "listener",
						},
						&ast.Ident{
							Name: "err",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "net",
								},
								Sel: &ast.Ident{
									Name: "Listen",
								},
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"tcp\"",
								},
								&ast.SelectorExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "s",
										},
										Sel: &ast.Ident{
											Name: "config",
										},
									},
									Sel: &ast.Ident{
										Name: "BindAddr",
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
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "s",
									},
									Sel: &ast.Ident{
										Name: "server",
									},
								},
								Sel: &ast.Ident{
									Name: "Serve",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "listener",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (s Server) SyncServerStart() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Start" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = s.AstServerStart()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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

func (s Server) AstServerStop() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{
							Name: "s",
						},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{
							Name: "Server",
						},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Stop",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "_",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
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
							X: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "s",
								},
								Sel: &ast.Ident{
									Name: "server",
								},
							},
							Sel: &ast.Ident{
								Name: "GracefulStop",
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

func (s Server) SyncServerStop() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "Stop" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = s.AstServerStop()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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

func (s Server) AstMessageProducer() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "DefaultMessageProducer",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "ctx",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "context",
							},
							Sel: &ast.Ident{
								Name: "Context",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "msg",
							},
						},
						Type: &ast.Ident{
							Name: "string",
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "level",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "zapcore",
							},
							Sel: &ast.Ident{
								Name: "Level",
							},
						},
					},
					{
						Names: []*ast.Ident{
							{
								Name: "code",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "codes",
							},
							Sel: &ast.Ident{
								Name: "Code",
							},
						},
					},
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
					{
						Names: []*ast.Ident{
							{
								Name: "duration",
							},
						},
						Type: &ast.SelectorExpr{
							X: &ast.Ident{
								Name: "zapcore",
							},
							Sel: &ast.Ident{
								Name: "Field",
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
							Name: "logger",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "ctxzap",
								},
								Sel: &ast.Ident{
									Name: "Extract",
								},
							},
							Args: []ast.Expr{
								&ast.Ident{
									Name: "ctx",
								},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "params",
						},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CompositeLit{
							Type: &ast.ArrayType{
								Elt: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "zap",
									},
									Sel: &ast.Ident{
										Name: "Field",
									},
								},
							},
							Elts: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "zap",
										},
										Sel: &ast.Ident{
											Name: "String",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"grpc.code\"",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "code",
												},
												Sel: &ast.Ident{
													Name: "String",
												},
											},
										},
									},
								},
								&ast.Ident{
									Name: "duration",
								},
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "zap",
										},
										Sel: &ast.Ident{
											Name: "Any",
										},
									},
									Args: []ast.Expr{
										&ast.BasicLit{
											Kind:  token.STRING,
											Value: "\"request_id\"",
										},
										&ast.CallExpr{
											Fun: &ast.SelectorExpr{
												X: &ast.Ident{
													Name: "ctx",
												},
												Sel: &ast.Ident{
													Name: "Value",
												},
											},
											Args: []ast.Expr{
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "log",
													},
													Sel: &ast.Ident{
														Name: "RequestIDKey",
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "sts",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "status",
											},
											Sel: &ast.Ident{
												Name: "Convert",
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
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "msg",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "sts",
											},
											Sel: &ast.Ident{
												Name: "Message",
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
									Name: "v",
								},
								Tok: token.DEFINE,
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "sts",
										},
										Sel: &ast.Ident{
											Name: "Details",
										},
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												&ast.Ident{
													Name: "errParams",
												},
											},
											Tok: token.DEFINE,
											Rhs: []ast.Expr{
												&ast.CompositeLit{
													Type: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "errs",
														},
														Sel: &ast.Ident{
															Name: "Params",
														},
													},
												},
											},
										},
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												&ast.Ident{
													Name: "badRequest",
												},
												&ast.Ident{
													Name: "ok",
												},
											},
											Tok: token.DEFINE,
											Rhs: []ast.Expr{
												&ast.TypeAssertExpr{
													X: &ast.Ident{
														Name: "v",
													},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "errdetails",
															},
															Sel: &ast.Ident{
																Name: "BadRequest",
															},
														},
													},
												},
											},
										},
										&ast.IfStmt{
											Cond: &ast.Ident{
												Name: "ok",
											},
											Body: &ast.BlockStmt{
												List: []ast.Stmt{
													&ast.RangeStmt{
														Key: &ast.Ident{
															Name: "_",
														},
														Value: &ast.Ident{
															Name: "violation",
														},
														Tok: token.DEFINE,
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "badRequest",
																},
																Sel: &ast.Ident{
																	Name: "GetFieldViolations",
																},
															},
														},
														Body: &ast.BlockStmt{
															List: []ast.Stmt{
																&ast.AssignStmt{
																	Lhs: []ast.Expr{
																		&ast.IndexExpr{
																			X: &ast.Ident{
																				Name: "errParams",
																			},
																			Index: &ast.CallExpr{
																				Fun: &ast.SelectorExpr{
																					X: &ast.Ident{
																						Name: "violation",
																					},
																					Sel: &ast.Ident{
																						Name: "GetField",
																					},
																				},
																			},
																		},
																	},
																	Tok: token.ASSIGN,
																	Rhs: []ast.Expr{
																		&ast.CallExpr{
																			Fun: &ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "violation",
																				},
																				Sel: &ast.Ident{
																					Name: "GetDescription",
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
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												&ast.Ident{
													Name: "errorInfo",
												},
												&ast.Ident{
													Name: "ok",
												},
											},
											Tok: token.DEFINE,
											Rhs: []ast.Expr{
												&ast.TypeAssertExpr{
													X: &ast.Ident{
														Name: "v",
													},
													Type: &ast.StarExpr{
														X: &ast.SelectorExpr{
															X: &ast.Ident{
																Name: "errdetails",
															},
															Sel: &ast.Ident{
																Name: "ErrorInfo",
															},
														},
													},
												},
											},
										},
										&ast.IfStmt{
											Cond: &ast.Ident{
												Name: "ok",
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
														X: &ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "errorInfo",
																},
																Sel: &ast.Ident{
																	Name: "GetMetadata",
																},
															},
														},
														Body: &ast.BlockStmt{
															List: []ast.Stmt{
																&ast.AssignStmt{
																	Lhs: []ast.Expr{
																		&ast.IndexExpr{
																			X: &ast.Ident{
																				Name: "errParams",
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
												},
											},
										},
										&ast.AssignStmt{
											Lhs: []ast.Expr{
												&ast.Ident{
													Name: "params",
												},
											},
											Tok: token.ASSIGN,
											Rhs: []ast.Expr{
												&ast.CallExpr{
													Fun: &ast.Ident{
														Name: "append",
													},
													Args: []ast.Expr{
														&ast.Ident{
															Name: "params",
														},
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "zap",
																},
																Sel: &ast.Ident{
																	Name: "Object",
																},
															},
															Args: []ast.Expr{
																&ast.BasicLit{
																	Kind:  token.STRING,
																	Value: "\"params\"",
																},
																&ast.Ident{
																	Name: "errParams",
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
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "logger",
									},
									Sel: &ast.Ident{
										Name: "Check",
									},
								},
								Args: []ast.Expr{
									&ast.Ident{
										Name: "level",
									},
									&ast.Ident{
										Name: "msg",
									},
								},
							},
							Sel: &ast.Ident{
								Name: "Write",
							},
						},
						Args: []ast.Expr{
							&ast.Ident{
								Name: "params",
							},
						},
						Ellipsis: 2996,
					},
				},
			},
		},
	}
}

func (s Server) SyncMessageProducer() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "DefaultMessageProducer" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = s.AstMessageProducer()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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

func (s Server) AstDecodeError() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: &ast.Ident{
			Name: "decodeError",
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
						Type: &ast.Ident{
							Name: "error",
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
										Name: "domainError",
									},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X: &ast.Ident{
											Name: "errs",
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
									Name: "domainError",
								},
							},
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "stat",
									},
								},
								Tok: token.DEFINE,
								Rhs: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{
											X: &ast.Ident{
												Name: "status",
											},
											Sel: &ast.Ident{
												Name: "New",
											},
										},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "codes",
													},
													Sel: &ast.Ident{
														Name: "Code",
													},
												},
												Args: []ast.Expr{
													&ast.SelectorExpr{
														X: &ast.Ident{
															Name: "domainError",
														},
														Sel: &ast.Ident{
															Name: "Code",
														},
													},
												},
											},
											&ast.SelectorExpr{
												X: &ast.Ident{
													Name: "domainError",
												},
												Sel: &ast.Ident{
													Name: "Message",
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
													Name: "withDetails",
												},
											},
											Type: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X: &ast.Ident{
														Name: "status",
													},
													Sel: &ast.Ident{
														Name: "Status",
													},
												},
											},
										},
									},
								},
							},
							&ast.SwitchStmt{
								Tag: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "domainError",
									},
									Sel: &ast.Ident{
										Name: "Code",
									},
								},
								Body: &ast.BlockStmt{
									List: []ast.Stmt{
										&ast.CaseClause{
											List: []ast.Expr{
												&ast.SelectorExpr{
													X: &ast.Ident{
														Name: "errs",
													},
													Sel: &ast.Ident{
														Name: "ErrorCodeInvalidArgument",
													},
												},
											},
											Body: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "d",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.UnaryExpr{
															Op: token.AND,
															X: &ast.CompositeLit{
																Type: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "errdetails",
																	},
																	Sel: &ast.Ident{
																		Name: "BadRequest",
																	},
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
													X: &ast.SelectorExpr{
														X: &ast.Ident{
															Name: "domainError",
														},
														Sel: &ast.Ident{
															Name: "Params",
														},
													},
													Body: &ast.BlockStmt{
														List: []ast.Stmt{
															&ast.AssignStmt{
																Lhs: []ast.Expr{
																	&ast.SelectorExpr{
																		X: &ast.Ident{
																			Name: "d",
																		},
																		Sel: &ast.Ident{
																			Name: "FieldViolations",
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
																					Name: "d",
																				},
																				Sel: &ast.Ident{
																					Name: "FieldViolations",
																				},
																			},
																			&ast.UnaryExpr{
																				Op: token.AND,
																				X: &ast.CompositeLit{
																					Type: &ast.SelectorExpr{
																						X: &ast.Ident{
																							Name: "errdetails",
																						},
																						Sel: &ast.Ident{
																							Name: "BadRequest_FieldViolation",
																						},
																					},
																					Elts: []ast.Expr{
																						&ast.KeyValueExpr{
																							Key: &ast.Ident{
																								Name: "Field",
																							},
																							Value: &ast.Ident{
																								Name: "key",
																							},
																						},
																						&ast.KeyValueExpr{
																							Key: &ast.Ident{
																								Name: "Description",
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
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "withDetails",
														},
														&ast.Ident{
															Name: "err",
														},
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "stat",
																},
																Sel: &ast.Ident{
																	Name: "WithDetails",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "d",
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
																				Name: "status",
																			},
																			Sel: &ast.Ident{
																				Name: "Error",
																			},
																		},
																		Args: []ast.Expr{
																			&ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "codes",
																				},
																				Sel: &ast.Ident{
																					Name: "Internal",
																				},
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
															},
														},
													},
												},
											},
										},
										&ast.CaseClause{
											Body: []ast.Stmt{
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "d",
														},
													},
													Tok: token.DEFINE,
													Rhs: []ast.Expr{
														&ast.UnaryExpr{
															Op: token.AND,
															X: &ast.CompositeLit{
																Type: &ast.SelectorExpr{
																	X: &ast.Ident{
																		Name: "errdetails",
																	},
																	Sel: &ast.Ident{
																		Name: "ErrorInfo",
																	},
																},
																Elts: []ast.Expr{
																	&ast.KeyValueExpr{
																		Key: &ast.Ident{
																			Name: "Reason",
																		},
																		Value: &ast.SelectorExpr{
																			X: &ast.Ident{
																				Name: "domainError",
																			},
																			Sel: &ast.Ident{
																				Name: "Message",
																			},
																		},
																	},
																	&ast.KeyValueExpr{
																		Key: &ast.Ident{
																			Name: "Domain",
																		},
																		Value: &ast.BasicLit{
																			Kind:  token.STRING,
																			Value: "\"\"",
																		},
																	},
																	&ast.KeyValueExpr{
																		Key: &ast.Ident{
																			Name: "Metadata",
																		},
																		Value: &ast.SelectorExpr{
																			X: &ast.Ident{
																				Name: "domainError",
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
												},
												&ast.AssignStmt{
													Lhs: []ast.Expr{
														&ast.Ident{
															Name: "withDetails",
														},
														&ast.Ident{
															Name: "err",
														},
													},
													Tok: token.ASSIGN,
													Rhs: []ast.Expr{
														&ast.CallExpr{
															Fun: &ast.SelectorExpr{
																X: &ast.Ident{
																	Name: "stat",
																},
																Sel: &ast.Ident{
																	Name: "WithDetails",
																},
															},
															Args: []ast.Expr{
																&ast.Ident{
																	Name: "d",
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
																				Name: "status",
																			},
																			Sel: &ast.Ident{
																				Name: "Error",
																			},
																		},
																		Args: []ast.Expr{
																			&ast.SelectorExpr{
																				X: &ast.Ident{
																					Name: "codes",
																				},
																				Sel: &ast.Ident{
																					Name: "Internal",
																				},
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
												Name: "withDetails",
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
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "status",
								},
								Sel: &ast.Ident{
									Name: "Error",
								},
							},
							Args: []ast.Expr{
								&ast.SelectorExpr{
									X: &ast.Ident{
										Name: "codes",
									},
									Sel: &ast.Ident{
										Name: "Internal",
									},
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
				},
			},
		},
	}
}

func (s Server) SyncDecodeError() error {
	fileset := token.NewFileSet()
	filename := path.Join("internal", "interfaces", "grpc", "server.go")
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	var methodExist bool
	var method *ast.FuncDecl
	ast.Inspect(file, func(node ast.Node) bool {
		if t, ok := node.(*ast.FuncDecl); ok && t.Name.String() == "decodeError" {
			methodExist = true
			method = t
			return false
		}
		return true
	})
	if method == nil {
		method = s.AstDecodeError()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
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
