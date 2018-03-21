package junocfg

import (
	"strings"
)

type item struct {
	path  []string
	value interface{}
}

func (i *item) pathString() string {
	return strings.Join(i.path, " / ")
}

type ItemArray []item

func (a ItemArray) Len() int { return len(a) }

func (a ItemArray) Less(i, j int) bool {
	return strings.Compare(a[i].pathString(), a[j].pathString()) < 1
}

func (a ItemArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
