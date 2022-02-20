package common

import (
	"github.com/samuel/go-thrift/parser"
	"github.com/weiyinfu/fugo/fu"
	"strings"
)

type Conf struct {
	SrcFile   string
	GoFile    string
	GoPackage string
	TsFile    string
}
type Api struct {
	ModelList   []*parser.Struct
	ServiceList []*parser.Service
}
func GetAnnotation(a []*parser.Annotation, name string) *parser.Annotation {
	for _, i := range a {
		if i.Name == name {
			return i
		}
	}
	return nil
}
func Parse(srcFile string) *Api {
	p := parser.Parser{}
	ma, _, err := p.ParseFile(srcFile)
	if err != nil {
		fu.Logger.Errorf("parse file error %v", err)
		panic(err)
		return nil
	}
	var modelList []*parser.Struct
	var serviceList []*parser.Service
	for filename, th := range ma {
		fu.Logger.Infof("parsing %s", filename)
		for _, structure := range th.Structs {
			modelList = append(modelList, structure)
		}
		for _, service := range th.Services {
			serviceList = append(serviceList, service)
		}
	}
	return &Api{
		ModelList: modelList, ServiceList: serviceList,
	}
}

func UpperFirst(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}
