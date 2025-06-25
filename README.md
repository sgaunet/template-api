# template-api

[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/template-api)](https://goreportcard.com/report/github.com/sgaunet/template-api)
![GitHub Downloads](https://img.shields.io/github/downloads/sgaunet/template-api/total)
![GitHub Release](https://img.shields.io/github/v/release/sgaunet/template-api)
![Test Coverage](https://raw.githubusercontent.com/wiki/sgaunet/template-api/coverage-badge.svg)
[![License](https://img.shields.io/github/license/sgaunet/template-api.svg)](LICENSE)
[![linter](https://github.com/sgaunet/template-api/actions/workflows/linter.yml/badge.svg)](https://github.com/sgaunet/template-api/actions/workflows/linter.yml)
[![coverage](https://github.com/sgaunet/template-api/actions/workflows/coverage.yml/badge.svg)](https://github.com/sgaunet/template-api/actions/workflows/coverage.yml)
[![snapshot](https://github.com/sgaunet/template-api/actions/workflows/snapshot.yml/badge.svg)](https://github.com/sgaunet/template-api/actions/workflows/snapshot.yml)
[![release](https://github.com/sgaunet/template-api/actions/workflows/release.yml/badge.svg)](https://github.com/sgaunet/template-api/actions/workflows/release.yml)

template-api is an API REST template project.

```bash
# install gonew
go install golang.org/x/tools/cmd/gonew@latest
# use gonew to create your project based on this template
gonew github.com/sgaunet/template-api gitplatform.com/username/awesome_new_project
cd awesome_new_project
git init
git add .
git remote add origin git@gitplatform.com:username/awesome_new_project
git push -u origin master
```

## Run

```
$ cat cfg.yaml
dbdsn: postgres://user:password@host:port/dbname?sslmode=disable
$ template-api -cfg cfg.yaml
...
```

## Install

* Download the binary in the release section
* Or use the docker image 


# Development

This project is using :

* golang
* [task for development](https://taskfile.dev/#/)
* docker
* [docker buildx](https://github.com/docker/buildx)
* docker manifest
* [goreleaser](https://goreleaser.com/)
* [pre-commit](https://pre-commit.com/)

There are hooks executed in the precommit stage. Once the project cloned on your disk, please install pre-commit:

```
brew install pre-commit
```

Install tools:

```
task dev:prereq
```

And install the hooks:

```
task dev:install
```

To launch manually the pre-commmit hook:

```
task dev:pre-commit
```
