Creathor is a CLI tool for generating layout and base CRUD operations on your project.
# Overview

Creathor provides:
* [Standart project layout](https://github.com/golang-standards/project-layout)
* Clean architecture with interfaces, interceptors, usecases, repositories and domain models 
* DI with [FX](https://github.com/uber-go/fx)
* Interface of [Logger](https://github.com/uber-go/zap) and clock
* [Changelog](https://keepachangelog.com/en/1.0.0/)
* Dockerfile and helm chart

# Installing
```
go install github.com/018bf/creathor@latest
```

# Usage

To initialize a project in the current directory, use the command `init`
```
creathor init --name tracker --module github.com/018bf/tracker --model user --model token --model equipment --model session --model approach
```

You can override the project path with the `-d` or `--destination` flag
```
creathor -d /Users/me/Projects/mysimpletracker init --name tracker --module github.com/018bf/tracker --model user --model token --model equipment --model session --model approach
```

To add new models use the `model` command
```
creathor model --model category --model item
```
creathor -d backend init --name ogorod --module git.jetbrains.space/elasticsoft/ogorod/Backend.git --model user --model notification --model notificationChannel --model token --model culture --model cultureCategory --model handbook --model handbookCategory --model handbookCard --model handbookLink --model work --model myWork --model workTag --model region --model masterClass --model masterClassCategory --model masterClassStep --model article --model articleCategory --model articleComment --model question --model questionLike --model answer --model answerLike