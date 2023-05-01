package configs

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"go/ast"
	"golang.org/x/exp/slices"
)

type ModelType uint8

const (
	ModelTypeMain = iota
	ModelTypeCreate
	ModelTypeUpdate
	ModelTypeFilter
)

type MethodType uint8

const (
	MethodTypeGet = iota
	MethodTypeList
	MethodTypeCount
	MethodTypeCreate
	MethodTypeUpdate
	MethodTypeDelete
)

type Model struct {
	Type       ModelType
	Name       string
	Variable   string
	Params     []*Param
	Validation bool
	Mock       bool
}

func NewCreateModel(modelConfig *ModelConfig) *Model {
	return &Model{
		Type:       ModelTypeCreate,
		Name:       modelConfig.CreateTypeName(),
		Variable:   "create",
		Params:     modelConfig.Params,
		Validation: true,
		Mock:       true,
	}
}

func NewUpdateModel(modelConfig *ModelConfig) *Model {
	model := &Model{
		Type:     ModelTypeUpdate,
		Name:     modelConfig.UpdateTypeName(),
		Variable: "update",
		Params: []*Param{
			{
				Name: "ID",
				Type: "UUID",
			},
		},
		Validation: true,
		Mock:       true,
	}
	for _, param := range modelConfig.Params {
		model.Params = append(model.Params, &Param{
			Name: param.GetName(),
			Type: fmt.Sprintf("*%s", param.Type),
		})
	}
	return model
}

func NewMainModel(modelConfig *ModelConfig) *Model {
	model := &Model{
		Type:     ModelTypeMain,
		Name:     modelConfig.ModelName(),
		Variable: modelConfig.Variable(),
		Params: []*Param{
			{
				Name:   "ID",
				Type:   "UUID",
				Search: false,
			},
			{
				Name:   "CreatedAt",
				Type:   "time.Time",
				Search: false,
			},
			{
				Name:   "UpdatedAt",
				Type:   "time.Time",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	model.Params = append(model.Params, modelConfig.Params...)
	return model
}

func NewFilterModel(modelConfig *ModelConfig) *Model {
	model := &Model{
		Type:     ModelTypeFilter,
		Name:     modelConfig.FilterTypeName(),
		Variable: "filter",
		Params: []*Param{
			{
				Name:   "PageSize",
				Type:   "*uint64",
				Search: false,
			},
			{
				Name:   "PageNumber",
				Type:   "*uint64",
				Search: false,
			},
			{
				Name:   "Search",
				Type:   "*string",
				Search: false,
			},
			{
				Name:   "OrderBy",
				Type:   "[]string",
				Search: false,
			},
			{
				Name:   "IDs",
				Type:   "[]UUID",
				Search: false,
			},
		},
		Validation: true,
		Mock:       true,
	}
	return model
}

type Method struct {
	Name   string
	Type   MethodType
	Args   []*ast.Field
	Return []*ast.Field
}

type UseCase struct {
	Name     string
	Variable string
	Methods  []*Method
}

func NewUseCase(m *ModelConfig) *UseCase {
	return &UseCase{
		Name:     m.UseCaseTypeName(),
		Variable: m.UseCaseVariableName(),
		Methods: []*Method{
			{
				Name: "Create",
				Type: MethodTypeCreate,
				Args: []*ast.Field{
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
							ast.NewIdent("create"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.CreateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "List",
				Type: MethodTypeList,
				Args: []*ast.Field{
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("[]*models.%s", m.ModelName())),
					},
					{
						Type: ast.NewIdent("uint64"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Get",
				Type: MethodTypeGet,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Update",
				Type: MethodTypeUpdate,
				Args: []*ast.Field{
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
							ast.NewIdent("update"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.UpdateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Delete",
				Type: MethodTypeDelete,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
}

type Repository struct {
	Name     string
	Variable string
	Methods  []*Method
}

func NewRepository(m *ModelConfig) *Repository {
	return &Repository{
		Name:     m.RepositoryTypeName(),
		Variable: m.RepositoryVariableName(),
		Methods: []*Method{
			{
				Name: "Create",
				Type: MethodTypeCreate,
				Args: []*ast.Field{
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
							ast.NewIdent(m.Variable()),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "List",
				Type: MethodTypeList,
				Args: []*ast.Field{
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("[]*models.%s", m.ModelName())),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Count",
				Type: MethodTypeCount,
				Args: []*ast.Field{
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("uint64"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Get",
				Type: MethodTypeGet,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Update",
				Type: MethodTypeUpdate,
				Args: []*ast.Field{
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
							ast.NewIdent(m.Variable()),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Delete",
				Type: MethodTypeDelete,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
}

type Interceptor struct {
	Auth     bool
	Events   bool
	Name     string
	Variable string
	Methods  []*Method
}

func NewInterceptor(m *ModelConfig) *Interceptor {
	interceptor := &Interceptor{
		Auth:   m.Auth,
		Events: m.KafkaEnabled,
		Name:   m.InterceptorTypeName(),
		Methods: []*Method{
			{
				Name: "Create",
				Type: MethodTypeCreate,
				Args: []*ast.Field{
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
							ast.NewIdent("create"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.CreateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "List",
				Type: MethodTypeList,
				Args: []*ast.Field{
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
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.FilterTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent(fmt.Sprintf("[]*models.%s", m.ModelName())),
					},
					{
						Type: ast.NewIdent("uint64"),
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Get",
				Type: MethodTypeGet,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Update",
				Type: MethodTypeUpdate,
				Args: []*ast.Field{
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
							ast.NewIdent("update"),
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.UpdateTypeName()),
							},
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("models"),
								Sel: ast.NewIdent(m.ModelName()),
							},
						},
					},
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
			{
				Name: "Delete",
				Type: MethodTypeDelete,
				Args: []*ast.Field{
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
							ast.NewIdent("id"),
						},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("models"),
							Sel: ast.NewIdent("UUID"),
						},
					},
				},
				Return: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
	}
	if interceptor.Auth {
		for _, method := range interceptor.Methods {
			method.Args = append(method.Args, &ast.Field{
				Names: []*ast.Ident{ast.NewIdent("requestUser")},
				Type: &ast.StarExpr{
					X: &ast.SelectorExpr{
						X:   ast.NewIdent("models"),
						Sel: ast.NewIdent("User"),
					},
				},
			})
		}
	}
	return interceptor
}

func (i *Interceptor) GetMethod(t MethodType) *Method {
	index := slices.IndexFunc(i.Methods, func(method *Method) bool { return method.Type == t })
	if index >= 0 {
		return i.Methods[index]
	}
	return nil
}

func (m *Interceptor) GetCreateMethod() *Method {
	return m.GetMethod(MethodTypeCreate)
}
func (m *Interceptor) GetUpdateMethod() *Method {
	return m.GetMethod(MethodTypeUpdate)
}
func (m *Interceptor) GetDeleteMethod() *Method {
	return m.GetMethod(MethodTypeDelete)
}
func (m *Interceptor) GetListMethod() *Method {
	return m.GetMethod(MethodTypeList)
}

func (m *Interceptor) GetGetMethod() *Method {
	return m.GetMethod(MethodTypeGet)
}

type Mod struct {
	Name        string
	Module      string
	Filename    string
	Models      []*Model
	UseCase     *UseCase
	Repository  *Repository
	Interceptor *Interceptor
	Auth        bool
}

func (m *Mod) GetMainModel() *Model {
	index := slices.IndexFunc(m.Models, func(model *Model) bool { return model.Type == ModelTypeMain })
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) GetCrateModel() *Model {
	index := slices.IndexFunc(m.Models, func(model *Model) bool { return model.Type == ModelTypeCreate })
	if index >= 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) GetUpdateModel() *Model {
	index := slices.IndexFunc(m.Models, func(model *Model) bool { return model.Type == ModelTypeUpdate })
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) GetFilterModel() *Model {
	index := slices.IndexFunc(m.Models, func(model *Model) bool { return model.Type == ModelTypeFilter })
	if index > 0 {
		return m.Models[index]
	}
	return nil
}

func (m *Mod) PermissionIDCreate() string {
	return fmt.Sprintf("PermissionID%sCreate", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDUpdate() string {
	return fmt.Sprintf("PermissionID%sUpdate", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDDelete() string {
	return fmt.Sprintf("PermissionID%sDelete", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDDetail() string {
	return fmt.Sprintf("PermissionID%sDetail", strcase.ToCamel(m.Name))
}

func (m *Mod) PermissionIDList() string {
	return fmt.Sprintf("PermissionID%sList", strcase.ToCamel(m.Name))
}
