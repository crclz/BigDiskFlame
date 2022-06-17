package domainservices

import (
	"BigDisk/domainmodels"
	"BigDisk/template"
	_ "embed"
	"encoding/json"
	"sort"
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
		Value: float64(unit.Size) / (1024 * 1024), // MB display
	}

	if depth <= 0 {
		return flame
	}

	if len(unit.Children) > 0 {
		sort.Slice(unit.Children, func(i, j int) bool {
			return unit.Children[i].Size > unit.Children[j].Size
		})
	}

	for _, u := range unit.Children {
		if u.Size < minSize {
			continue
		}
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
