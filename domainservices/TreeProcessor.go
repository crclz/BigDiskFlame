package domainservices

import (
	"BigDisk/domainmodels"
	"BigDisk/template"
	_ "embed"
	"encoding/json"
	"strings"
)

type TreeProcessor struct {
}

func NewTreeProcessor() *TreeProcessor {
	return &TreeProcessor{}
}

var singletonTreeProcessor *TreeProcessor = initSingletonTreeProcessor()

func GetSingletonTreeProcessor() *TreeProcessor {
	return singletonTreeProcessor
}

func initSingletonTreeProcessor() *TreeProcessor {
	return NewTreeProcessor()
}

func (p *TreeProcessor) Process(unit *domainmodels.FileUnit) {
	if unit == nil {
		panic("unit is nil")
	}

	for _, u := range unit.Children {
		p.Process(u)

		unit.Size += u.Size
	}
}

func (p *TreeProcessor) ToFlameNode(unit *domainmodels.FileUnit, minSize int64, depth int) *domainmodels.FlameNode {
	var flame = &domainmodels.FlameNode{
		Name:  unit.Name,
		Value: float64(unit.Size) / (1024 * 1024),
	}

	if unit.Size < minSize {
		return flame
	}

	if depth <= 0 {
		return flame
	}

	for _, u := range unit.Children {
		var node = p.ToFlameNode(u, minSize, depth-1)
		flame.Children = append(flame.Children, node)
	}

	return flame
}

func (p *TreeProcessor) GenerateReportHtml(node *domainmodels.FlameNode) string {
	if node == nil {
		panic("node is nil")
	}
	flameDataBytes, err := json.Marshal(node)

	if err != nil {
		panic(err)
	}

	var html = strings.ReplaceAll(template.GetHtmlTemplate(), "{flameDataPlaceHolder}", string(flameDataBytes))

	return html
}
