package ts

import (
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
	fileContent := fu.Format(`import {cli} from "@/api"

{{{ModelString}}}

{{{ServiceString}}}
`, bson.M{
		"PackageName":   conf.GoPackage,
		"ModelString":   modelString.String(),
		"ServiceString": serviceString.String(),
	})
	fu.WriteFile(conf.TsFile, fileContent)
}

func getService(s *parser.Service) string {
	//为service生成interface
	var methodList []string
	for _, m := range s.Methods {
		methodList = append(methodList, fu.Format(`
    {{MethodName}}(req:{{ReqName}}):Promise<{{RespName}}> {
        return cli.post("{{Url}}", req).then(resp=>{
            if (resp.status != 200) throw new Error("网络错误")
            return resp.data;
        })
    },
`, bson.M{
			"MethodName": m.Name,
			"ReqName":    getTsType(m.Arguments[0].Type),
			"RespName":   getTsType(m.ReturnType),
			"Url":        common.GetAnnotation(m.Annotations, "url").Value,
		}))
	}
	template := `
export const {{Name}}Service = {
{{{MethodList}}}
}
`
	return fu.Format(template, bson.M{
		"Name":       s.Name,
		"MethodList": strings.Join(methodList, "\n"),
	})
}
func getTsType(t *parser.Type) string {
	if t.Name == "list" {
		return getTsType(t.ValueType) + "[]"
	}
	if t.Name == "map" {
		return fu.Format(`{[index:{{KeyType}}]:{{ValueType}}}`, bson.M{
			"KeyType":   getTsType(t.KeyType),
			"ValueType": getTsType(t.ValueType),
		})
	}
	ma := map[string]string{
		"i32":    "number",
		"string": "string",
		"bool":   "boolean",
		"i64":    "number",
		"double": "number",
	}
	if ans, ok := ma[t.Name]; ok {
		return ans
	}
	return t.Name
}

func getModel(m *parser.Struct) string {
	//获取一个model的结构体定义
	template := `
export interface {{Name}} {
{{{Lines}}}
}
`
	var lines []string
	for _, f := range m.Fields {
		tsType := getTsType(f.Type)
		line := fu.Format("    {{FieldName}}: {{FieldType}}", bson.M{
			"FieldName": f.Name,
			"FieldType": tsType,
		})
		lines = append(lines, line)
	}
	return fu.Format(template, bson.M{
		"Name":  m.Name,
		"Lines": strings.Join(lines, "\n"),
	})
}
