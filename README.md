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
go install github.com/mikalai-mitsin/creathor@latest
```

# Usage

Config example

```yaml
name: "example"
module: "github.com/mikalai-mitsin/example"
goVersion: "1.22"
auth: true
ci: "github"
gRPC: true
http: true
gateway: false
uptrace: true
domains:
  - model: Post
    params:
      - name: "title"
        type: "string"
      - name: "body"
        type: "string"
      - name: "is_private"
        type: "bool"
      - name: "tags"
        type: "[]string"
      - name: "published_at"
        type: "time.Time"
      - name: "author_id"
        type: "uuid.UUID"
  - model: Comment
    params:
      - name: "text"
        type: "string"
      - name: "author_id"
        type: "uuid.UUID"
      - name: "post_id"
        type: "uuid.UUID"
```

To generate code in the current directory and with default config name, use the command `creathor`