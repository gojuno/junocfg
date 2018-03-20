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

type itemArray []item

func (a itemArray) Len() int { return len(a) }

func (a itemArray) Less(i, j int) bool {
	return strings.Compare(a[i].pathString(), a[j].pathString()) < 1
}

func (a itemArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
