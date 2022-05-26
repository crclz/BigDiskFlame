package main

import (
	"BigDisk/domainservices"
	"BigDisk/infra"
	"encoding/json"
	"fmt"
)

func main() {
	var treeService = infra.GetSingletonTreeService()
	var treeProcessor = domainservices.GetSingletonTreeProcessor()

	var unit, err = treeService.GetUnit("./domainmodels")

	if err != nil {
		panic(err)
	}

	fmt.Printf("unit addr: %p\n", unit)
	treeProcessor.Process(unit)

	d, err := json.Marshal(unit)
	if err != nil {
		panic(err)
	}

	fmt.Printf("unit: %s\n", d)
}
