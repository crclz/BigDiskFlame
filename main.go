package main

import (
	"BigDisk/domainmodels"
	"BigDisk/domainservices"
	"BigDisk/infra"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main_old() {
	var config = &domainmodels.RunConfig{}
	configData, err := ioutil.ReadFile("config.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, config)

	if err != nil {
		panic(err)
	}

	var treeService = infra.GetSingletonTreeService()
	var treeProcessor = domainservices.GetSingletonTreeProcessor()

	unit, err := treeService.GetUnit(config.Path)

	if err != nil {
		panic(err)
	}

	treeProcessor.Process(unit)

	var flame = treeProcessor.ToFlameNode(unit, config.ReportMinSize, config.MaxDepth)

	var html = treeProcessor.GenerateReportHtml(flame)

	// output
	var htmlFilename = fmt.Sprintf("%v.result.html", time.Now().UnixMilli())

	err = ioutil.WriteFile(htmlFilename, []byte(html), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	var treeService = infra.GetSingletonTreeService()

	file, err := os.Open("du-result.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	treeService.GetUnitFromDuResult(file)
}
