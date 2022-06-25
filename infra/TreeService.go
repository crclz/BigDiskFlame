package infra

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/crclz/BigDiskFlame/domainmodels"
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

func (p *TreeService) SelectTrimmedWhereNotWhitespace(x []string) []string {
	var result []string

	for _, s := range x {
		s = strings.Trim(s, " \r\t")

		if len(s) != 0 {
			result = append(result, s)
		}
	}

	return result
}

func (p *TreeService) GetUnitFromDuResult(reader io.Reader) (*domainmodels.FileUnit, error) {
	var scanner = bufio.NewScanner(reader)

	var rootNode = &domainmodels.FileUnit{Name: ""}

	var nodeMap = map[string]*domainmodels.FileUnit{
		"": rootNode,
	} // key="/data00/abc" value=node

	var GetNode func(path string) *domainmodels.FileUnit

	GetNode = func(path string) *domainmodels.FileUnit {
		if path == "" || path == "." {
			return nodeMap[""]
		}

		path = strings.TrimRight(path, "/")

		// 责任链
		var node = nodeMap[path]

		if node != nil {
			return node
		}

		var parentPath = path[:strings.LastIndex(path, "/")]
		var name = path[strings.LastIndex(path, "/")+1:]

		var parentNode = GetNode(parentPath)

		var currentNode = &domainmodels.FileUnit{
			Name:   name,
			IsFile: false,
		}

		parentNode.Children = append(parentNode.Children, currentNode)

		nodeMap[path] = currentNode

		return currentNode
	}

	for scanner.Scan() {
		var text = scanner.Text()
		var parts = regexp.MustCompile("[\t]").Split(text, -1)
		parts = p.SelectTrimmedWhereNotWhitespace(parts)

		if len(parts) != 2 {
			fmt.Printf("Skip line: %v\n", text)
		}

		size, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			panic(err)
		}

		var node = GetNode(parts[1])
		node.Size = size * 1024 // KB -> byte
	}

	return nodeMap[""], nil
}
