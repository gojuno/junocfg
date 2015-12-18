# junocfg
Template based config  generator

## Installation

```
go get github.com/juno-lab/junocfg
```

## Usage

```
$ junocfg -t config.yaml.tmpl -c settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg -t config.yaml.tmpl > config.yaml
```
