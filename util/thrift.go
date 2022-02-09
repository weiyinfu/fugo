package util

import (
	"github.com/cloudwego/thriftgo/parser"
)

type SS map[string][]string

func GetAnnotation(annotations []*parser.Annotation, name string) *parser.Annotation {
	for _, ano := range annotations {
		if ano.Key == name {
			return ano
		}
	}
	return nil
}
func BuildAnnotation(annotations SS) []*parser.Annotation {
	var a []*parser.Annotation
	for k, v := range annotations {
		a = append(a, &parser.Annotation{Key: k, Values: v})
	}
	return a
}
func GetAnnotationValue(annotations []*parser.Annotation, name string) string {
	ano := GetAnnotation(annotations, name)
	if ano == nil || len(ano.Values) == 0 {
		return ""
	}
	return ano.Values[0]
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
