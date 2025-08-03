Creathor is a CLI tool for generating layout and base CRUD operations on your project.

# Overview

Creathor provides:

* [Standart project layout](https://github.com/golang-standards/project-layout)
* Clean architecture with interfaces, usecases, services, repositories and domain entities
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
ci: "github"
gRPC: true
http: true
gateway: false
uptrace: true
apps:
  - name: posts
    entities:
      - name: post
        params:
          - name: "Body"
            type: "string"
      - name: tag
        params:
          - name: "post_id"
            type: "uuid.UUID"
          - name: "Value"
            type: "string"
      - name: like
        params:
          - name: "PostID"
            type: "uuid.UUID"
          - name: "Value"
            type: "string"
          - name: "user_id"
            type: "uuid.UUID"
  - name: articles
    entities:
      - name: article
        params:
          - name: "Title"
            type: "string"
          - name: "Subtitle"
            type: "string"
          - name: "Body"
            type: "string"
          - name: "is_published"
            type: "string"
```

To generate code in the current directory and with default config name, use the command `creathor`