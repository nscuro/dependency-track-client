# dependency-track-client

![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)

## Installation

`GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...`

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
    --bom ./bom.xml --autocreate
```

#### BOM

##### Export

```
$ ./dtrack bom export \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    -o bom.xml
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
