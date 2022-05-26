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

	var flame = treeProcessor.ToFlameNode(unit)

	d, err := json.Marshal(flame)
	if err != nil {
		panic(err)
	}

	fmt.Printf("flame: %s\n", d)
}
