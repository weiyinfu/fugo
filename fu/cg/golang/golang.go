package golang

import (
	"fmt"
	"github.com/samuel/go-thrift/parser"
	"github.com/weiyinfu/fugo/fu"
	"github.com/weiyinfu/fugo/fu/cg/common"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

func Handle(api *common.Api, conf *common.Conf) {
	modelString := &strings.Builder{}
	for _, m := range api.ModelList {
		modelString.WriteString(getModel(m))
	}
	serviceString := &strings.Builder{}
	for _, s := range api.ServiceList {
		serviceString.WriteString(getService(s))
	}
	fileContent := fu.Format(`package {{PackageName}}


import (
    "github.com/gin-gonic/gin"
    "github.com/weiyinfu/fugo/fu"
    "net/http"
)

{{{ModelString}}}

{{{ServiceString}}}
`, bson.M{
		"PackageName":   conf.GoPackage,
		"ModelString":   modelString.String(),
		"ServiceString": serviceString.String(),
	})
	fu.WriteFile(conf.GoFile, fileContent)
}

func getService(s *parser.Service) string {
	//为service生成interface
	var methodList []string
	var ginUrlList []string
	for _, m := range s.Methods {
		responseType := getGoType(m.ReturnType, false)
		if len(m.Arguments) != 1 {
			panic(fmt.Sprintf("%v.%v参数个数应该为1", s.Name, m.Name))
		}
		requestType := getGoType(m.Arguments[0].Type, false)
		goMethodName := common.UpperFirst(m.Name)
		methodList = append(methodList, fu.Format(`    {{GoMethodName}}(req * {{RequestType}})*{{ResponseType}}`, bson.M{
			"GoMethodName": goMethodName,
			"RequestType":  requestType,
			"ResponseType": responseType,
		}))
		url := common.GetAnnotation(m.Annotations, "url")
		if url == nil {
			panic(fmt.Sprintf("unknown url to method %v", m.Name))
		}

		methodModel := bson.M{
			"Url":          url.Value,
			"GoMethodName": goMethodName,
			"RequestType":  requestType,
			"ResponseType": responseType,
			"ServiceName":  s.Name,
		}
		ginUrlList = append(ginUrlList, fu.Format(`    x.POST("{{Url}}", func(c *gin.Context) {
        req := &{{RequestType}}{}
        err := fu.Bind(c, req)
        if err != nil {
            panic(err)
        }
        resp := s.{{GoMethodName}}(req)
        c.JSON(http.StatusOK, resp)
    })`, methodModel))
	}
	template := `
type {{Name}}Service interface {
{{MethodList}}
}

func Register{{Name}}(x *gin.Engine, s {{Name}}Service) {
{{{UrlList}}}
}
`
	goName := common.UpperFirst(s.Name)
	return fu.Format(template, bson.M{
		"Name":        goName,
		"MethodList":  strings.Join(methodList, "\n"),
		"UrlList":     strings.Join(ginUrlList, "\n"),
	})
}
func getGoType(t *parser.Type, usePointer bool) string {
	if t.Name == "list" {
		return "[]" + getGoType(t.ValueType, usePointer)
	}
	if t.Name == "map" {
		return fu.Format(`map[{{KeyType}}]{{ValueType}}`, bson.M{
			"KeyType":   getGoType(t.KeyType, usePointer),
			"ValueType": getGoType(t.ValueType, usePointer),
		})
	}
	ma := map[string]string{
		"i32":    "int",
		"string": "string",
		"bool":   "bool",
		"i64":    "int64",
		"double": "float64",
	}
	if ans, ok := ma[t.Name]; ok {
		return ans
	}
	if usePointer {
		return "*" + t.Name
	}
	return t.Name
}
func getModel(m *parser.Struct) string {
	//获取一个model的结构体定义
	template := `
type {{Name}} struct {
{{{Lines}}}
}
`
	var lines []string
	for _, f := range m.Fields {
		goFieldName := strings.ToUpper(f.Name[:1]) + f.Name[1:]
		goType := getGoType(f.Type, true)
		line := fu.Format("    {{GoFieldName}} {{FieldType}} `json:\"{{FieldName}}\"`", bson.M{
			"GoFieldName": goFieldName,
			"FieldType":   goType,
			"FieldName":   f.Name,
		})
		lines = append(lines, line)
	}
	return fu.Format(template, bson.M{
		"Name":  m.Name,
		"Lines": strings.Join(lines, "\n"),
	})
}
