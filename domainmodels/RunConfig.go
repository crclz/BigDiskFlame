package domainmodels

type RunConfig struct {
	Path          string `json:"path"`
	ReportMinSize int64  `json:"reportMinSize"`
	MaxDepth      int    `json:"maxDepth"`
}
