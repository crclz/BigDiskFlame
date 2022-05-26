package infra

import (
	"BigDisk/domainmodels"
	"fmt"
	"os"
	"path/filepath"
)

type TreeService struct {
}

func NewTreeService() *TreeService {
	return &TreeService{}
}

var singletonTreeService *TreeService = initSingletonTreeService()

func GetSingletonTreeService() *TreeService {
	return singletonTreeService
}

func initSingletonTreeService() *TreeService {
	return NewTreeService()
}

func (p *TreeService) GetUnit(path string) (*domainmodels.FileUnit, error) {
	var info, err = os.Lstat(path)

	if err != nil {
		fmt.Printf("Lstat error. Path=%v, err=%v\n", path, err)
		return nil, err
	}

	if !info.IsDir() {
		return &domainmodels.FileUnit{
			Name:   info.Name(),
			IsFile: true,
			Size:   info.Size(),
		}, nil
	}

	var unit = &domainmodels.FileUnit{
		Name: filepath.Base(path),
		Size: 0,
	}

	childrens, err := os.ReadDir(path)

	if err != nil {
		fmt.Printf("ReadDir error. Path=%v, err=%v\n", path, err)
		return nil, err
	}

	for _, x := range childrens {
		var xu, err = p.GetUnit(filepath.Join(path, x.Name()))

		if err != nil {
			continue
		}

		if xu == nil {
			panic("xu is null")
		}

		unit.Children = append(unit.Children, xu)
	}

	return unit, nil
}
