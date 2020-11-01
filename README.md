# dependency-track-client

[![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)](https://github.com/nscuro/dependency-track-client/actions?query=workflow%3A%22Continuous+Integration%22)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/nscuro/dependency-track-client)](https://pkg.go.dev/github.com/nscuro/dependency-track-client)

Unofficial Go client library and CLI for [Dependency-Track](https://github.com/DependencyTrack/dependency-track)

## Installation

For library usage:

```
GO111MODULE=on go get -v github.com/nscuro/dependency-track-client
```

With `dtrack` command:

```
GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...
```

## Compatibility

* Go >= 1.15
* Dependency-Track >= 4.0.0

## API Coverage

The library primarily covers those parts of the Dependency-Track API that are needed for the CLI application.
If you'd like to use this library, and your desired functionality is not yet available, please consider creating a PR.

## Usage

```
Usage:
  dtrack [command]

Available Commands:
  audit       Audit for vulnerabilities
  bom         Export and Upload BOMs
  help        Help about any command
  report      Generate reports
  version     Display version information

Flags:
  -k, --apikey string            Dependency-Track API Key
  -h, --help                     help for dtrack
      --project string           Project UUID
      --project-name string      Project Name
      --project-version string   Project Version
  -u, --url string               Dependency-Track URL
```

Dependency-Track's URL and the API key can be provided via environment variables as well:

```
$ export DTRACK_URL=https://dependencytrack.example.com
$ export DTRACK_APIKEY=0sl67mjen99zxb2y
```

### Examples

#### Audit

```
$ ./dtrack audit \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom ./bom.xml --autocreate \
    --gate ./examples/qualitygate.yaml
```

#### BOM

##### Export

```
$ ./dtrack bom export \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    -o bom.xml
```

##### Status

```
$ ./dtrack bom status \
    --token e043867f-b055-465f-814b-38f3330c2ec2
```

##### Upload

```
$ ./dtrack bom upload \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom bom.xml --autocreate
```

#### Report

```
$ ./dtrack report \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --template ./examples/project-report.gohtml \
    --output report.html
```
