package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/crclz/BigDiskFlame/domainmodels"
	"github.com/crclz/BigDiskFlame/domainservices"
	"github.com/crclz/BigDiskFlame/infra"
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
	var duFilename = flag.String("du", "du-result.txt", "To make du result: du ${path} > du-result.txt")
	var minSizeMB = flag.Int64("min-size-mb", 200, "200. Has impact on html performance")
	var maxDepth = flag.Int("max-depth", 8, "8. Has impact on html performance")
	var outputHtml = flag.String("out-html", "", "The output html filename.")

	flag.Parse()

	if *outputHtml == "" {
		*outputHtml = fmt.Sprintf("%v.disk.html", time.Now().Unix())
	}

	var treeService = infra.GetSingletonTreeService()
	var treeProcessor = domainservices.GetSingletonTreeProcessor()

	file, err := os.Open(*duFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	root, err := treeService.GetUnitFromDuResult(file)

	treeProcessor.Process(root)

	if err != nil {
		panic(err)
	}

	var flameRoot = treeProcessor.ToFlameNode(root, *minSizeMB*1024*1024, *maxDepth)

	var html = treeProcessor.GenerateReportHtml(flameRoot)

	ioutil.WriteFile(*outputHtml, []byte(html), 0644)
}
