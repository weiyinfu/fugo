package util

import (
	"fmt"
	"github.com/samuel/go-thrift/parser"
)

type SS map[string]string

func GetAnnotation(annotations []*parser.Annotation, name string) *parser.Annotation {
	for _, ano := range annotations {
		if ano.Name == name {
			return ano
		}
	}
	return nil
}
func BuildAnnotation(annotations SS) []*parser.Annotation {
	var a []*parser.Annotation
	for k, v := range annotations {
		a = append(a, &parser.Annotation{Name: k, Value: v})
	}
	return a
}
func GetAnnotationValue(annotations []*parser.Annotation, name string) string {
	ano := GetAnnotation(annotations, name)
	if ano == nil {
		return ""
	}
	return ano.Value
}
func HasAnnotation(annotations []*parser.Annotation, name string) bool {
	return GetAnnotation(annotations, name) != nil
}
func GetUniqueTypes(typeList []*parser.Type) []*parser.Type {
	//对typeList进行去重
	had := map[string]bool{}
	var typeSet []*parser.Type
	for _, typeInfo := range typeList {
		if typeInfo == nil {
			continue
		}
		if !had[typeInfo.Name] {
			had[typeInfo.Name] = true
			typeSet = append(typeSet, typeInfo)
		}
	}
	return typeSet
}

func ParseThrift(filepath string) map[string]*parser.Thrift {
	p := parser.Parser{}
	file2content, thriftPath, err := p.ParseFile(filepath)
	if err != nil {
		Logger.Error("parse file error: file=%s,err=%s", thriftPath, err)
		panic(fmt.Sprintf("parse thrift file error=%v,filepath=%v", err, filepath))
	}
	return file2content
}
