package main

import "BigDisk/infra"

func main() {
	var treeService = infra.GetSingletonTreeService()

	treeService.GetUnit(".")
}
