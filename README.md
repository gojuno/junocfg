# junocfg [![GoDoc](https://godoc.org/github.com/gojuno/junocfg?status.svg)](http://godoc.org/github.com/gojuno/junocfg) [![Build Status](https://travis-ci.org/gojuno/junocfg.svg?branch=master)](https://travis-ci.org/gojuno/junocfg)

Template based config generator / settings files merge tool

## Installation

```
go get github.com/gojuno/junocfg/...
```

## Usage

Modes:
- merge - merge config data from multiply sources to one dataset
- check-tmpl

### merge multiply config files to one

```
$ junocfg --merge -i public.yaml,secure.yaml -o settings.yaml

$ junocfg --merge -i public.yaml,secure.yaml > settings.yaml
```

### generate template from settings file

```
$ junocfg -t config.yaml.tmpl -i settings.dev.yaml -o config.yaml

$ cat settings.dev.yaml | junocfg -t config.yaml.tmpl > config.yaml
```

### check tmpl

```
$ junocfg --check-tmpl -t config.yaml.tmpl
$ cat config.yaml.tmpl | junocfg --check-tmpl
```

### full pipiline / generate template from multiply config files

merge config data from multiply sources to one dataset and apply it to template

```
$ junocfg --merge -i public.yaml,secure.yaml | junocfg -t config.yaml.tmpl > config.yaml
```

### convert between json & yaml

```
echo '{"a":"a","b":["b1","b2","b3","b4"],"c":{"ca":"ca","cb":["cb1","cb2",{"cb3a":"cb3aa"},"cb4",{"cb5a":"cb5aa"}]},"d":"a","e":{"ea":"ea","eb":"eb","ec":"ec"}}' | junocfg --yaml2json | jq '.'

echo '{"a":"a","b":["b1","b2","b3","b4"],"c":{"ca":"ca","cb":["cb1","cb2",{"cb3a":"cb3aa"},"cb4",{"cb5a":"cb5aa"}]},"d":"a","e":{"ea":"ea","eb":"eb","ec":"ec"}}' | junocfg --json2yaml
```

## Test

```
junocfg --merge -i examples/a.yaml,examples/b.yaml

junocfg --check-tmpl -t examples/c.tmpl
cat examples/c.tmpl | junocfg --check-tmpl

echo '{"a":"a","b":["b1","b2","b3","b4"],"c":{"ca":"ca","cb":["cb1","cb2",{"cb3a":"cb3aa"},"cb4",{"cb5a":"cb5aa"}]},"d":"a","e":{"ea":"ea","eb":"eb","ec":"ec"}}' | junocfg --yaml2json | jq '.' | junocfg --json2yaml
```