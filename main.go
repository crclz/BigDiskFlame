package main

import (
	"BigDisk/domainservices"
	"BigDisk/infra"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	var treeService = infra.GetSingletonTreeService()
	var treeProcessor = domainservices.GetSingletonTreeProcessor()

	var unit, err = treeService.GetUnit("C:/Users/chr")

	if err != nil {
		panic(err)
	}

	treeProcessor.Process(unit)

	var flame = treeProcessor.ToFlameNode(unit, 1024*1024*50, 9)

	var html = treeProcessor.GenerateReportHtml(flame)

	// output
	var htmlFilename = fmt.Sprintf("%v.result.html", time.Now().UnixMilli())

	err = ioutil.WriteFile(htmlFilename, []byte(html), 0644)
	if err != nil {
		panic(err)
	}
}
