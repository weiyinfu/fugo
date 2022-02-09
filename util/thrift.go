package util

import (
	"fmt"

	"code.byted.org/gopkg/logs"
	"code.byted.org/gopkg/thriftparser"
)

type SS map[string]string

func GetAnnotation(annotations []*thriftparser.Annotation, name string) *thriftparser.Annotation {
	for _, ano := range annotations {
		if ano.Name == name {
			return ano
		}
	}
	return nil
}
func BuildAnnotation(annotations SS) []*thriftparser.Annotation {
	var a []*thriftparser.Annotation
	for k, v := range annotations {
		a = append(a, &thriftparser.Annotation{Name: k, Value: v})
	}
	return a
}
func GetAnnotationValue(annotations []*thriftparser.Annotation, name string) string {
	ano := GetAnnotation(annotations, name)
	if ano == nil {
		return ""
	}
	return ano.Value
}
func HasAnnotation(annotations []*thriftparser.Annotation, name string) bool {
	return GetAnnotation(annotations, name) != nil
}
func GetUniqueTypes(typeList []*thriftparser.Type) []*thriftparser.Type {
	//对typeList进行去重
	had := map[string]bool{}
	var typeSet []*thriftparser.Type
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

func ParseThrift(filepath string) map[string]*thriftparser.Thrift {
	p := thriftparser.Parser{}
	file2content, thriftPath, err := p.ParseFile(filepath)
	if err != nil {
		logs.Error("parse file error: file=%s,err=%s", thriftPath, err)
		panic(fmt.Sprintf("parse thrift file error=%v,filepath=%v", err, filepath))
	}
	return file2content
}
