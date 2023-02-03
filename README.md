Creathor is a CLI tool for generating layout and base CRUD operations on your project.

# Overview

Creathor provides:

* [Standart project layout](https://github.com/golang-standards/project-layout)
* Clean architecture with interfaces, interceptors, usecases, repositories and domain models
* DI with [FX](https://github.com/uber-go/fx)
* Interface of [Logger](https://github.com/uber-go/zap) and clock
* gRPC and RESTful APIs
* PostgreSQL repositories and migrations
* CI/CD configurations for Github and GitLab
* [Changelog](https://keepachangelog.com/en/1.0.0/)
* Dockerfile and helm chart

# Example

[Example](/example) of using `Creathor` with rich models and authentication by this command

```shell
creathor -d example -c creathor.yaml
```

# Installing

```
go install github.com/018bf/creathor@latest
```

# Usage

Config example

```yaml
name: "example"                    // Service name
module: "github.com/018bf/example" // Module name
goVersion: "1.19"                  // go version
auth: true                         // Generate auth, permissions and users api
ci: "github"                       // Add `gitlab` or `github` CI configs
gRPC: true                         // Generate gRPC API
REST: false                        // RESTful API
models:
  - model: "session" // Model name
    params:
      - name: "title"  // Parameter name
        type: "string" // Parameter type one of "int", "int64", "int32", "int16", "int8", "[]int", "[]int64", "[]int32", "[]int16", "[]int8", "uint", "uint64", "uint32", "uint16", "uint8", "[]uint", "[]uint64", "[]uint32", "[]uint16", "[]uint8", "string", "[]string", "time.Time", "[]time.Time",
        search: false  // Add field to the search
      - name: "description"
        type: "string"
        search: true
  - model: "equipment"
    params:
      - name: "name"
        type: "string"
        search: true
      - name: "repeat"
        type: "int"
        search: false
      - name: "weight"
        type: "int"
        search: false
```

To generate code in the current directory and with default config name, use the command `creathor`