# junocfg

Template based config  generator

## Installation

```
go get github.com/juno-lab/junocfg
```

## Usage

### work with one settings file

```
$ junocfg -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg -t config.yaml.tmpl > config.yaml
```

### work with multiply config files

`junocfg` merge config data from multiply sources to one dataset and apply it to template

```
$ junocfg -t config.yaml.tmpl -i public.yaml,secure.yaml -o config.yaml

$ junocfg -t config.yaml.tmpl -i public.yaml,secure.yaml > config.yaml
```

## Check

```
$ junocfg --check -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg --check -t config.yaml.tmpl > config.yaml
```
