package domainservices

import (
	_ "embed"
	"encoding/json"
	"sort"
	"strings"

	"github.com/crclz/BigDiskFlame/domainmodels"
	"github.com/crclz/BigDiskFlame/template"
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

	if unit.Size != 0 {
		return
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

	// check

	if depth <= 0 {
		return flame
	}

	// merge small files
	var newChildren []*domainmodels.FileUnit
	var smallFileUnit = &domainmodels.FileUnit{Name: "-", Size: 0}

	for _, u := range unit.Children {
		if u.Size >= minSize {
			newChildren = append(newChildren, u)
		} else {
			smallFileUnit.Size += u.Size
		}
	}

	if smallFileUnit.Size > 0 {
		newChildren = append(newChildren, smallFileUnit)
	}

	// order by size desc
	if len(newChildren) > 0 {
		sort.Slice(newChildren, func(i, j int) bool {
			return newChildren[i].Size > newChildren[j].Size
		})
	}

	for _, u := range newChildren {
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

	var html = strings.ReplaceAll(template.GetHtmlTemplate(), "{ flameDataPlaceHolder }", string(flameDataBytes))

	return html
}
