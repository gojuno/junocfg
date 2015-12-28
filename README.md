# junocfg
Template based config  generator

## Installation

```
go get github.com/juno-lab/junocfg
```

## Usage

```
$ junocfg -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg -t config.yaml.tmpl > config.yaml
```

## Check

```
$ junocfg --check -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg --check -t config.yaml.tmpl > config.yaml
```
