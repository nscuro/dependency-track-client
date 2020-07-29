# dependency-track-client

![Build Status](https://github.com/nscuro/dependency-track-client/workflows/Continuous%20Integration/badge.svg)

## Installation

`GO111MODULE=on go get -v github.com/nscuro/dependency-track-client/...`

## Usage

```
$ ./dtrack -h
Usage:
  dtrack [command]

Available Commands:
  audit       audit for vulnerabilities
  bom         retrieve or upload boms
  help        Help about any command
  report      generate reports

Flags:
  -k, --api-key string           dependency-track api key
  -h, --help                     help for dtrack
      --project-name string      project name
      --project-uuid string      project uuid
      --project-version string   project version
  -u, --url string               dependency-track base url

Use "dtrack [command] --help" for more information about a command.
```

Dependency-Track's URL and the API key can be provided via environment variables as well:

```
$ export DTRACK_BASE_URL=https://dependencytrack.evilcorp.com
$ export DTRACK_API_KEY=0sl67mjen99zxb2y
```

### Audit

```
$ ./dtrack audit \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom ./bom.xml --autocreate
```

### BOM

```
$ ./dtrack bom get \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    -o bom.xml
```

```
$ ./dtrack bom upload \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --bom bom.xml --autocreate
```

### Report

```
$ ./dtrack report \
    --project-name Dependency-Track \
    --project-version 3.8.0 \
    --template ./examples/project-report.html \
    --output report.html
```
