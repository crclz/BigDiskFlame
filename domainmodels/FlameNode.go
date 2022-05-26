package domainmodels

type FlameNode struct {
	Name     string       `json:"name"`
	Value    float64      `json:"value"`
	Children []*FlameNode `json:"children"`
}
