package grpc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"

	"github.com/mikalai-mitsin/creathor/internal/pkg/astfile"
	"github.com/mikalai-mitsin/creathor/internal/pkg/configs"
)

type HandlerGenerator struct {
	domain *configs.EntityConfig
}

func NewHandlerGenerator(domain *configs.EntityConfig) *HandlerGenerator {
	return &HandlerGenerator{
		domain: domain,
	}
}

func (h HandlerGenerator) file() *ast.File {
	importSpec := []ast.Spec{
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"context"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.EntitiesImportPath(),
			},
		},
		&ast.ImportSpec{
			Name: ast.NewIdent(h.domain.ProtoPackage),
			Path: &ast.BasicLit{
				Kind: token.STRING,
				Value: fmt.Sprintf(
					`"%s/pkg/%s/v1"`,
					h.domain.Module,
					h.domain.ProtoPackage,
				),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.AppConfig.ProjectConfig.PointerImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: h.domain.AppConfig.ProjectConfig.UUIDImportPath(),
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/emptypb"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/timestamppb"`,
			},
		},
		&ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: `"google.golang.org/protobuf/types/known/wrapperspb"`,
			},
		},
	}
	for _, param := range h.domain.GetUpdateModel().Params {
		if param.IsSlice() {
			importSpec = append(importSpec, &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"google.golang.org/protobuf/types/known/structpb"`,
				},
			})
			break
		}
	}
	return &ast.File{
		Name: ast.NewIdent("handlers"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok:   token.IMPORT,
				Specs: importSpec,
			},
		},
	}
}

func (h HandlerGenerator) filename() string {
	return path.Join(
		"internal",
		"app",
		h.domain.AppConfig.AppName(),
		"handlers",
		"grpc",
		h.domain.DirName(),
		h.domain.FileName(),
	)
}

func (h HandlerGenerator) structure() *ast.TypeSpec {
	return &ast.TypeSpec{
		Name: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
		Type: &ast.StructType{
			Fields: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.SelectorExpr{
							X: ast.NewIdent(h.domain.ProtoPackage),
							Sel: &ast.Ident{
								Name: fmt.Sprintf(
									"Unimplemented%sServiceServer",
									h.domain.GetMainModel().Name,
								),
							},
						},
					},
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
	}
}

func (h HandlerGenerator) syncStruct() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		file = h.file()
	}
	structure, structureExists := astfile.FindType(file, h.domain.GetGRPCHandlerTypeName())
	if structure == nil {
		structure = h.structure()
	}
	if !structureExists {
		gd := &ast.GenDecl{
			Tok:   token.TYPE,
			Specs: []ast.Spec{structure},
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

func (h HandlerGenerator) constructor() *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent(fmt.Sprintf("New%s", h.domain.GetGRPCHandlerTypeName())),
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
							X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
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
									Name: fmt.Sprintf(
										"%sServiceServer",
										h.domain.GetMainModel().Name,
									),
								},
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
	}
}

func (h HandlerGenerator) syncConstructor() error {
	fileset := token.NewFileSet()
	filename := h.filename()
	file, err := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, h.domain.GetGRPCHandlerConstructorName())
	if method == nil {
		method = h.constructor()
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

func (h HandlerGenerator) create() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetCreateModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("input"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Create"),
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
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetCreateModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Create"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("decode%s", h.domain.GetMainModel().Name),
							),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncCreateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Create")
	if method == nil {
		method = h.create()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}
	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) get() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("uuid"),
				Sel: ast.NewIdent("MustParse"),
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetId"),
					},
				},
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Get"),
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
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(
									fmt.Sprintf("%sGet", h.domain.GetMainModel().Name),
								),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Get"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("decode%s", h.domain.GetMainModel().Name),
							),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncGetMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Get")
	if method == nil {
		method = h.get()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) list() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetFilterModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("filter"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("List"),
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
					{
						Names: []*ast.Ident{
							ast.NewIdent("filter"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetFilterModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(
									fmt.Sprintf("List%s", h.domain.GetMainModel().Name),
								),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("items"),
						ast.NewIdent("count"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("List"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(h.domain.GetGRPCMainListDecodeName()),
							Args: []ast.Expr{
								ast.NewIdent("items"),
								ast.NewIdent("count"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncListMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "List")
	if method == nil {
		method = h.list()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) update() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: ast.NewIdent(fmt.Sprintf("encode%s", h.domain.GetUpdateModel().Name)),
			Args: []ast.Expr{
				ast.NewIdent("input"),
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Update"),
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
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetUpdateModel().Name),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(h.domain.GetMainModel().Name),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						ast.NewIdent("item"),
						ast.NewIdent("err"),
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent("s"),
									Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
								},
								Sel: ast.NewIdent("Update"),
							},
							Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: ast.NewIdent(
								fmt.Sprintf("decode%s", h.domain.GetMainModel().Name),
							),
							Args: []ast.Expr{
								ast.NewIdent("item"),
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncUpdateMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Update")
	if method == nil {
		method = h.update()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) delete() *ast.FuncDecl {
	args := []ast.Expr{
		ast.NewIdent("ctx"),
		&ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("uuid"),
				Sel: ast.NewIdent("MustParse"),
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("input"),
						Sel: ast.NewIdent("GetId"),
					},
				},
			},
		},
	}
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
					},
				},
			},
		},
		Name: ast.NewIdent("Delete"),
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
					{
						Names: []*ast.Ident{
							ast.NewIdent("input"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X: ast.NewIdent(h.domain.ProtoPackage),
								Sel: ast.NewIdent(
									fmt.Sprintf("%sDelete", h.domain.GetMainModel().Name),
								),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("emptypb"),
								Sel: ast.NewIdent("Empty"),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
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
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("s"),
										Sel: ast.NewIdent(h.domain.GetUseCasePrivateVariableName()),
									},
									Sel: ast.NewIdent("Delete"),
								},
								Args: args,
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
									ast.NewIdent("nil"),
									ast.NewIdent("err"),
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							Op: token.AND,
							X: &ast.CompositeLit{
								Type: &ast.SelectorExpr{
									X:   ast.NewIdent("emptypb"),
									Sel: ast.NewIdent("Empty"),
								},
							},
						},
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}
}

func (h HandlerGenerator) syncDeleteMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindFunc(file, "Delete")
	if method == nil {
		method = h.delete()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) syncRegisterMethod() error {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, h.filename(), nil, parser.ParseComments)
	if err != nil {
		return err
	}
	method, methodExist := astfile.FindMethod(file, h.domain.GRPCHandlerTypeName(), "RegisterGRPC")
	if method == nil {
		method = h.registerGRPC()
	}
	if !methodExist {
		file.Decls = append(file.Decls, method)
	}

	buff := &bytes.Buffer{}
	if err := printer.Fprint(buff, fileset, file); err != nil {
		return err
	}
	if err := os.WriteFile(h.filename(), buff.Bytes(), 0777); err != nil {
		return err
	}
	return nil
}

func (h HandlerGenerator) registerGRPC() *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						ast.NewIdent("s"),
					},
					Type: &ast.StarExpr{
						X: ast.NewIdent(h.domain.GetGRPCHandlerTypeName()),
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
									X:   ast.NewIdent(h.domain.ProtoPackage),
									Sel: ast.NewIdent(h.domain.GetGRPCServiceDescriptionName()),
								},
							},
							ast.NewIdent("s"),
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

func (h HandlerGenerator) Sync() error {
	err := os.MkdirAll(path.Dir(h.filename()), 0777)
	if err != nil {
		return err
	}
	if err := h.syncStruct(); err != nil {
		return err
	}
	if err := h.syncConstructor(); err != nil {
		return err
	}
	if err := h.syncCreateMethod(); err != nil {
		return err
	}
	if err := h.syncGetMethod(); err != nil {
		return err
	}
	if err := h.syncListMethod(); err != nil {
		return err
	}
	if err := h.syncUpdateMethod(); err != nil {
		return err
	}
	if err := h.syncDeleteMethod(); err != nil {
		return err
	}
	if err := h.syncRegisterMethod(); err != nil {
		return err
	}
	return nil
}
