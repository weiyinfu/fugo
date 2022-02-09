package util

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/aymerick/raymond"
)

func UniqueStringArray(a []string) []string {
	had := map[string]bool{}
	var b []string
	for _, i := range a {
		if had[i] {
			continue
		}
		had[i] = true
		b = append(b, i)
	}
	return b
}
func Format(template string, ctx bson.M) string {
	//正则表达式检测所有字符串，所有的字符串都应该在ctx中出现
	s, err := raymond.Render(template, ctx)
	if err != nil {
		panic(fmt.Sprintf("render template error :template=%v,error=%v", template, err))
	}
	return s
}
func ReverseStringArray(a []string) {
	for i := 0; i < len(a)/2; i++ {
		a[i] = a[len(a)-1-i]
	}
}
