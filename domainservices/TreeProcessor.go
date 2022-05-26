package domainservices

import (
	"BigDisk/domainmodels"
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

func (p *TreeProcessor) ToFlameNode(unit *domainmodels.FileUnit) *domainmodels.FlameNode {
	var flame = &domainmodels.FlameNode{
		Name:  unit.Name,
		Value: unit.Size,
	}

	for _, u := range unit.Children {
		var node = p.ToFlameNode(u)
		flame.Children = append(flame.Children, node)
	}

	return flame
}
