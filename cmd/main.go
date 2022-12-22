package main

import (
	"homework/new2"
	"homework/old"
)

//https://godoc.org/golang.org/x/perf/cmd/benchstat

func main() {

	// исходная версия
	old.GetTransformDocuments()

	// оптимизированная версия
	new2.GetTransformDocuments()

}
