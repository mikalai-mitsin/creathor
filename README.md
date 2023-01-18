Creathor is a CLI tool for generating layout and base CRUD operations on your project.
# Overview

Creathor provides:
* [Standart project layout](https://github.com/golang-standards/project-layout)
* Clean architecture with interfaces, interceptors, usecases, repositories and domain models 
* DI with [FX](https://github.com/uber-go/fx)
* Interface of [Logger](https://github.com/uber-go/zap) and clock
* gRPC and RESTful APIs
* PostgreSQL repositories
* CI/CD configurations for Github and GitLab
* [Changelog](https://keepachangelog.com/en/1.0.0/)
* Dockerfile and helm chart

# Example
[Example](/example) of using `Creathor` with rich models and authentication by this command
```shell
creathor -d ./example -ci github -a init --name example --module "github.com/018bf/example" --model '{"model":"session", "params": {"title": "string", "description": "string"}}' --model '{"model":"equipment", "params": {"name": "string", "repeat": "int", "weight": "int"}}'
```

# Installing
```
go install github.com/018bf/creathor@latest
```

# Usage

To initialize a project in the current directory, use the command `init`
```
creathor init --name tracker --module github.com/018bf/tracker --model user --model token --model equipment --model session --model approach --model '{"model":"mark", "params": {"name": "string", "title": "string", "weight": "int"}}'
```

To add auth and permission checks add `-a` or `--auth` flag
```
creathor -a init --name tracker --module github.com/018bf/tracker --model user --model token --model equipment --model session --model approach
```

You can override the project path with the `-d` or `--destination` flag
```
creathor -d /Users/me/Projects/mysimpletracker init --name tracker --module github.com/018bf/tracker --model user --model token --model equipment --model session --model approach
```

To add new models use the `model` command
```
creathor model --model category --model item
```
