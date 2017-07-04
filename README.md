# junocfg [![GoDoc](https://godoc.org/github.com/gojuno/junocfg?status.svg)](http://godoc.org/github.com/gojuno/junocfg) [![Build Status](https://travis-ci.org/gojuno/junocfg.svg?branch=master)](https://travis-ci.org/gojuno/junocfg)

Template based config  generator

## Installation

```
go get github.com/gojuno/junocfg
```

## Usage

### generate template from one settings file

```
$ junocfg -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg -t config.yaml.tmpl > config.yaml
```

### generate template from multiply config files

`junocfg` merge config data from multiply sources to one dataset and apply it to template

```
$ junocfg -t config.yaml.tmpl -i public.yaml,secure.yaml -o config.yaml

$ junocfg -t config.yaml.tmpl -i public.yaml,secure.yaml > config.yaml
```

### merge multiply config files to one

```
$ junocfg --merge -i public.yaml,secure.yaml -o settings.yaml

$ junocfg --merge -i public.yaml,secure.yaml > settings.yaml
```


### check

```
$ junocfg --check -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg --check -t config.yaml.tmpl > config.yaml
```

### full pipiline

```
$ junocfg --merge -i public.yaml,secure.yaml -o settings.yaml

$ junocfg --merge -i public.yaml,secure.yaml | junocfg --check -t config.yaml.tmpl > config.yaml
```
