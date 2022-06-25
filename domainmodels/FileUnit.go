package domainmodels

type FileUnit struct {
	Name     string
	IsFile   bool
	Size     int64
	Children []*FileUnit
}
