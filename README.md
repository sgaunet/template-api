[![Go Report Card](https://goreportcard.com/badge/github.com/sgaunet/template-api)](https://goreportcard.com/report/github.com/sgaunet/template-api)

# template-api

template-api is an API REST template project.

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
