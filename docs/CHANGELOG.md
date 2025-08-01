<a name="unreleased"></a>
## [Unreleased]


<a name="v0.6.1"></a>
## [v0.6.1] - 2025-07-08
### Feat
- use real uuid
- add validation to order by
- add name to files


<a name="v0.6.0"></a>
## [v0.6.0] - 2025-03-22
### Fear
- update config for readme
- use ast constructor

### Feat
- add http handlers
- add http register to auth app
- fix auth tests
- use auth with values
- fix uuid in proto file
- fix tests
- fix json decode
- remove required validation for bools
- fix tests
- fix tests
- replace pointers with value in grpc handlers
- replace pointers with value in services
- replace pointers with value in repositories
- fix update
- fill many dto to entity
- fill filter dto to entity
- fill update dto to entity
- fill create dto to entity
- fill dto constructors
- add api prefix
- add bearer and task file
- add swagger annotation to server
- add swagger annotation


<a name="v0.5.0"></a>
## [v0.5.0] - 2024-12-09
### Chore
- clean

### Feat
- rename models with entities
- use auth in app
- add http cmd
- add http container
- add http server
- use config in app build
- mount handler to router
- add http handler generator
- fix search
- use private interfaces
- fix tests
- local mock
- add assert to tests
- cleanup layout names
- remove layout
- remove methods
- rename interceptor with usecases
- rename usecases with services
- use private interfaces
- add register func to app
- remove auth from interceptor
- add auth middleware
- add middlewares to grpc
- invoke auth app
- remove permissions from auth
- invoke auth module first
- remove logger
- start app
- update configs
- add app generator
- cleanup
- add logger to repository
- add logger, clock and uuid interfaces to usecase
- add logger to handlers
- rename handlers
- **errors:** update is function
- **interceptors:** add logger
- **interceptors:** clean up interfaces

### Fix
- update imports
- auth flag


<a name="v0.4.8"></a>
## [v0.4.8] - 2024-06-18
### Feat
- add uuid generator
- fix tests
- fix imports
- update errors


<a name="vv0.5.0"></a>
## [vv0.5.0] - 2024-03-08
### Feat
- rename package
- move interfaces to internal pkg
- move interfaces to internal pkg
- move domains to app dir
- move errors to internal pkg
- add password and email generation
- add new arch
- something old
- **configs:** generate interceptors from mod
- **configs:** generate repository from mod
- **configs:** generate grpc handler from mod
- **configs:** add ids to the filter model
- **configs:** generate use case methods by mod
- **configs:** generate use case methods by mod
- **configs:** generate use case from mod
- **configs:** generate interfaces from config
- **fakes:** add separate generator
- **kafka:** add event repository
- **models:** use separate generators
- **perms:** add perms if not exists

### Fix
- fix search
- fix migrations
- update postgresql test
- renaming
- auth
- update auth
- generate user from domain
- generate user from domain
- update rest path
- generate uptrace if enabled

### Refactor
- move templates to generators
- move templates to generators
- sync di from pkg
- update config
- update package name
- update package name
- renaming
- renaming
- rename package
- remove method type
- **generators:** separate mods with layout
- **mods:** separate mods with configs

### Style
- remove extra call


<a name="v0.4.7"></a>
## [v0.4.7] - 2023-04-06
### Docs
- rebuild example

### Feat
- **models:** generate auth and user from ast
- **repositories:** generate interfaces from ast
- **rest:** add auth and user tests
- **rest:** add crud tests


<a name="v0.4.6"></a>
## [v0.4.6] - 2023-04-03
### Feat
- **implementations:** generate auth and user from ast
- **rest:** handle error
- **rest:** update docs
- **rest:** register handlers

### Fix
- imports
- **generators:** fix module name

### Refactor
- separate implementation packages


<a name="v0.4.5"></a>
## [v0.4.5] - 2023-03-22
### Feat
- **auth:** update handler params
- **fx:** sync with file
- **uptrace:** add to fx

### Fix
- **auth:** update auth usage
- **filter:** add search if it enabled


<a name="v0.4.4"></a>
## [v0.4.4] - 2023-03-04
### Build
- update dependencies

### Feat
- clean up di
- **container:** run servers as gorutine
- **domain:** generate files from ast
- **generators:** generate grpc interface
- **grpc:** generate server from ast
- **grpc:** add otel middleware
- **interceptor:** cleanup templates
- **interceptor:** generate implementation from AST
- **interfaces:** generate from AST
- **models:** generate filter from ast
- **models:** generate from AST
- **repositories:** generate implementation of create from AST
- **repositories:** generate implementation from AST
- **usecases:** generate implementation from AST

### Fix
- enable post init
- **container:** generate fx containers from ast
- **docker:** set go version
- **repositories:** enable search only if needed

### Refactor
- cleanup
- cleanup
- move to implementation
- add base crud generator
- move generators to separate packages
- separate model with generators
- **generators:** use project config in errors
- **models:** separate model, update and filter
- **usecases:** remove params

### Style
- remove comments


<a name="v0.4.3"></a>
## [v0.4.3] - 2023-02-08
### Docs
- rebuild example
- regenerate example

### Feat
- **gateway:** add count to metadata
- **gateway:** add gateway generator

### Fix
- **premissions:** add sub package


<a name="v0.4.2"></a>
## [v0.4.2] - 2023-02-01
### Ci
- **release:** set tag


<a name="v0.4.1"></a>
## [v0.4.1] - 2023-02-01
### Feat
- **apis:** possible to disable grpc or rest api
- **build:** add task
- **build:** use task
- **ci:** update validation rules
- **di:** extract DI to container package
- **grpc:** update proto package
- **grpc:** fill update test
- **models:** add creation validation
- **postgres:** add dto
- **repositories:** add args to the test
- **repositories:** separate repositories implementation

### Refactor
- move models to pkg


<a name="v0.4.0"></a>
## [v0.4.0] - 2023-01-27
### Docs
- regenerate example and update readme

### Feat
- generate project by file
- add default permission to migration
- use seq migrations
- **config:** add config to init
- **crud:** add postgres migrations
- **grpc:** update case
- **grpc:** register handlers
- **grpc:** add implementation
- **models:** replace id type
- **rest:** add openapi generation
- **templates:** skip if file exists

### Fix
- **grpc:** fix delete test

### Style
- remove fmt


<a name="v0.3.3"></a>
## [v0.3.3] - 2023-01-18
### Docs
- updated example

### Feat
- add rest
- **permissions:** add model permissions to repository
- **rest:** add user handler

### Fix
- use error interface
- use error interface
- **user:** update tests


<a name="v0.3.2"></a>
## [v0.3.2] - 2023-01-16
### Ci
- add release action


<a name="v0.3.1"></a>
## [v0.3.1] - 2023-01-16
### Ci
- add release action

### Feat
- setup build scripts
- **ci:** add github action
- **ci:** add gitlab

### Fix
- **repositories:** set default page size

### Style
- cleanup templates


<a name="v0.3.0"></a>
## [v0.3.0] - 2023-01-08
### Feat
- add reach model flag


<a name="v0.2.1"></a>
## [v0.2.1] - 2023-01-08
### Build
- add regex to makefile

### Docs
- add auth flag to readme

### Feat
- add grpc

### Refactor
- range files


<a name="v0.2.0"></a>
## [v0.2.0] - 2023-01-07
### Feat
- **auth:** generate auth as part of layout
- **di:** add postgres module to container


<a name="v0.1.7"></a>
## [v0.1.7] - 2023-01-06
### Feat
- build variables
- **auth:** check user model
- **deployments:** generate helm chart
- **di:** add implementation to DI
- **interceptor:** add tests


<a name="v0.1.6"></a>
## [v0.1.6] - 2022-12-31
### Docs
- cleanup readme

### Feat
- add auth flag
- add repository tests
- **errors:** decode domain error from postgres
- **interceptor:** add permission check
- **models:** add mock
- **usecase:** add count and tests


<a name="v0.1.5"></a>
## [v0.1.5] - 2022-12-28

<a name="v0.1.3"></a>
## [v0.1.3] - 2022-12-28

<a name="v0.1.4"></a>
## [v0.1.4] - 2022-12-27
### Docs
- add example

### Feat
- clean and generate mocks

### Fix
- update errors, repository and config


<a name="v0.1.1"></a>
## [v0.1.1] - 2022-10-06
### Docs
- add readme


<a name="v0.1.0"></a>
## v0.1.0 - 2022-10-06
### Feat
- add changelog
- add docs
- set go version
- add repository implementation
- add implementations
- add uber fx
- remove logrus
- add log
- add clock
- use err package name


[Unreleased]: https://github.com/mikalai-mitsin/creathor/compare/v0.6.1...HEAD
[v0.6.1]: https://github.com/mikalai-mitsin/creathor/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.5.0...v0.6.0
[v0.5.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.8...v0.5.0
[v0.4.8]: https://github.com/mikalai-mitsin/creathor/compare/vv0.5.0...v0.4.8
[vv0.5.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.7...vv0.5.0
[v0.4.7]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.6...v0.4.7
[v0.4.6]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.5...v0.4.6
[v0.4.5]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.4...v0.4.5
[v0.4.4]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.3...v0.4.4
[v0.4.3]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.2...v0.4.3
[v0.4.2]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.1...v0.4.2
[v0.4.1]: https://github.com/mikalai-mitsin/creathor/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.3.3...v0.4.0
[v0.3.3]: https://github.com/mikalai-mitsin/creathor/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/mikalai-mitsin/creathor/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/mikalai-mitsin/creathor/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.2.1...v0.3.0
[v0.2.1]: https://github.com/mikalai-mitsin/creathor/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.7...v0.2.0
[v0.1.7]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.6...v0.1.7
[v0.1.6]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.3...v0.1.5
[v0.1.3]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.4...v0.1.3
[v0.1.4]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.1...v0.1.4
[v0.1.1]: https://github.com/mikalai-mitsin/creathor/compare/v0.1.0...v0.1.1
