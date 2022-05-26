package domainmodels

type FlameNode struct {
	Name     string       `json:"name"`
	Value    int64        `json:"value"`
	Children []*FlameNode `json:"children"`
}
