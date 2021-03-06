# Clair

[![Build Status](https://api.travis-ci.org/coreos/clair.svg?branch=master "Build Status")](https://travis-ci.org/coreos/clair)
[![Docker Repository on Quay](https://quay.io/repository/coreos/clair/status "Docker Repository on Quay")](https://quay.io/repository/coreos/clair)
[![Go Report Card](https://goreportcard.com/badge/coreos/clair "Go Report Card")](https://goreportcard.com/report/coreos/clair)
[![GoDoc](https://godoc.org/github.com/coreos/clair?status.svg "GoDoc")](https://godoc.org/github.com/coreos/clair)
[![IRC Channel](https://img.shields.io/badge/freenode-%23clair-blue.svg "IRC Channel")](http://webchat.freenode.net/?channels=clair)

Clair is an open source project for the static analysis of vulnerabilities in [appc] and [docker] containers.

Vulnerability data is continuously imported from a known set of sources and correlated with the indexed contents of container images in order to produce lists of vulnerabilities that threaten a container.
When vulnerability data changes upstream, the previous state and new state of the vulnerability along with the images they affect can be sent via webhook to a configured endpoint.
New data sources can be [added programmatically] at compile-time or data can be injected via HTTP API at runtime.

Our goal is to enable a more transparent view of the security of container-based infrastructure.
Thus, the project was named `Clair` after the French term which translates to *clear*, *bright*, *transparent*.

[appc]: https://github.com/appc/spec
[docker]: https://github.com/docker/docker/blob/master/image/spec/v1.md
[added programmatically]: #custom-data-sources

## Common Use Cases

### Manual Auditing

You're building an application and want to depend on a third-party container image that you found by searching the internet.
To make sure that you do not knowingly introduce a new vulnerability into your production service, you decide to scan the container for vulnerabilities.
You `docker pull` the container to your development machine and start an instance of Clair.
Once it finishes updating, you use the [local image analysis tool] to analyze the container.
You realize this container is vulnerable to many critical CVEs, so you decide to use another one.

[local image analysis tool]: https://github.com/coreos/clair/tree/master/contrib/analyze-local-images

### Container Registry Integration

Your company has a continuous-integration pipeline and you want to stop deployments if they introduce a dangerous vulnerability.
A developer merges some code into the master branch of your codebase.
The first step of your continuous-integration pipeline automates the testing and building of your container and pushes a new container to your container registry.
Your container registry notifies Clair which causes the download and indexing of the images for the new container.
Clair detects some vulnerabilities and sends a webhook to your continuous deployment tool to prevent this vulnerable build from seeing the light of day.

## Hello Heartbleed

### Requirements

An instance of [PostgreSQL] 9.4+ is required.
All instructions assume the user has already setup this instance.
During the first run, Clair will bootstrap its database with vulnerability data from its data sources.
This can take several minutes.

[PostgreSQL]: http://postgresql.org

### Docker

The easiest way to get an instance of Clair running is to simply pull down the latest copy from Quay.

```sh
$ mkdir $HOME/clair_config
$ curl -L https://raw.githubusercontent.com/coreos/clair/config.example.yaml -o $HOME/clair_config/config.yaml
$ $EDITOR $HOME/clair_config/config.yaml # Add the URI for your postgres database
$ docker run quay.io/coreos/clair -p 6060-6061:6060-6061 -v $HOME/clair_config:/config -config=config.yaml
```

### Source

To build Clair, you need to latest stable version of [Go] and a working [Go environment].

[Go]: https://github.com/golang/go/releases
[Go environment]: https://golang.org/doc/code.html

```sh
$ go get github.com/coreos/clair
$ go install github.com/coreos/clair/cmd/clair
$ $EDITOR config.yaml # Add the URI for your postgres database
$ ./$GOBIN/clair -config=config.yaml
```

## Architecture

### At a glance

![Simple Clair Diagram](img/simple_diagram.png)

### Documentation

Documentation can be found in a README.md file located in the directory of the component.

- [Notifier](https://github.com/coreos/clair/blob/master/notifier/README.md)
- [v1 API](https://github.com/coreos/clair/blob/master/api/v1/README.md)

### Vulnerability Analysis

There are two major ways to perform analysis of programs: [Static Analysis] and [Dynamic Analysis].
Clair has been designed to perform *static analysis*; containers never need to be executed.
Rather, the filesystem of the container image is inspected and *features* are indexed into a database.
Features are anything that when present could be an indication of a vulnerability (e.g. the presence of a file or an installed software package).
By indexing the features of an image into the database, images only need to be rescanned when new features are added.

[Static Analysis]: https://en.wikipedia.org/wiki/Static_program_analysis
[Dynamic Analysis]: https://en.wikipedia.org/wiki/Dynamic_program_analysis

### Data Sources

| Data Source                   | Versions                                               | Format |
|-------------------------------|--------------------------------------------------------|--------|
| [Debian Security Bug Tracker] | 6, 7, 8, unstable                                      | [dpkg] |
| [Ubuntu CVE Tracker]          | 12.04, 12.10, 13.04, 14.04, 14.10, 15.04, 15.10, 16.04 | [dpkg] |
| [Red Hat Security Data]       | 5, 6, 7                                                | [rpm]  |

[Debian Security Bug Tracker]: https://security-tracker.debian.org/tracker
[Ubuntu CVE Tracker]: https://launchpad.net/ubuntu-cve-tracker
[Red Hat Security Data]: https://www.redhat.com/security/data/metrics
[dpkg]: https://en.wikipedia.org/wiki/dpkg
[rpm]: http://www.rpm.org


### Custom Data Sources

In addition to the default data sources, Clair has been designed in a way that allows extension without forking the project.
*Fetchers*, which are Go packages that implement the fetching of upstream vulnerability data, are registered in [init()] similar to drivers for Go's standard [database/sql] package.
A fetcher can live in its own repository and custom versions of clair can contain a small patch that adds the import statements of the desired fetchers in `main.go`.

[init()]: https://golang.org/doc/effective_go.html#init
[database/sql]: https://godoc.org/database/sql

## Related Links

- [Talk](https://www.youtube.com/watch?v=PA3oBAgjnkU) and [Slides](https://docs.google.com/presentation/d/1toUKgqLyy1b-pZlDgxONLduiLmt2yaLR0GliBB7b3L0/pub?start=false&loop=false&slide=id.p) @ ContainerDays NYC 2015
- [Quay](https://quay.io): the first container registry to integrate with Clair
- [Dockyard](https://github.com/containerops/dockyard): an open source container registry with Clair integration
