package domainservices

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

	childrens, err := os.ReadDir(filepath.Join(path, info.Name()))

	if err != nil {
		fmt.Printf("ReadDir error. Path=%v, err=%v", filepath.Join(path, info.Name()), err)
		return nil, err
	}

	for _, x := range childrens {
		var xu, err = p.GetUnit(filepath.Join(path, x.Name()))

		if err == nil {
			fmt.Printf("GetUnit error. Path=%v, err=%v", path, err)
			continue
		}

		if xu == nil {
			panic("xu is null")
		}

		unit.Children = append(unit.Children, xu)
	}

	return unit, nil
}
