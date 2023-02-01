<a name="unreleased"></a>
## [Unreleased]


<a name="0.4.2"></a>
## [0.4.2] - 2023-02-01
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
- **grpc:** fill update test
- **grpc:** update proto package
- **models:** add creation validation
- **postgres:** add dto
- **repositories:** separate repositories implementation
- **repositories:** add args to the test

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


[Unreleased]: https://github.com/018bf/creathor/compare/0.4.2...HEAD
[0.4.2]: https://github.com/018bf/creathor/compare/v0.4.1...0.4.2
[v0.4.1]: https://github.com/018bf/creathor/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/018bf/creathor/compare/v0.3.3...v0.4.0
[v0.3.3]: https://github.com/018bf/creathor/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/018bf/creathor/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/018bf/creathor/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/018bf/creathor/compare/v0.2.1...v0.3.0
[v0.2.1]: https://github.com/018bf/creathor/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/018bf/creathor/compare/v0.1.7...v0.2.0
[v0.1.7]: https://github.com/018bf/creathor/compare/v0.1.6...v0.1.7
[v0.1.6]: https://github.com/018bf/creathor/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/018bf/creathor/compare/v0.1.3...v0.1.5
[v0.1.3]: https://github.com/018bf/creathor/compare/v0.1.4...v0.1.3
[v0.1.4]: https://github.com/018bf/creathor/compare/v0.1.1...v0.1.4
[v0.1.1]: https://github.com/018bf/creathor/compare/v0.1.0...v0.1.1
